// main.go
package main

//go:generate go run ./codegen/gen_apis.go
//go:generate go run ./codegen/gen_chain_ids.go
//go:generate go run ./codegen/gen_blockscout_apis.go
//go:generate go fmt ./apis_generated.go
//go:generate go fmt ./chain_ids_generated.go
//go:generate go fmt ./blockscout_apis_generated.go

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	_ "github.com/mark3labs/anyabi.xyz/migrations"
	"github.com/mark3labs/anyabi.xyz/views"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core" 
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	datastar "github.com/starfederation/datastar/sdk/go"
	"golang.org/x/time/rate"
)

type GetABISignals struct {
	Address string `json:"address"`
	ChainId string `json:"chainId"`
}

type EtherscanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		SourceCode           string `json:"SourceCode"`
		ABI                  string `json:"ABI"`
		ContractName         string `json:"ContractName"`
		CompilerVersion      string `json:"CompilerVersion"`
		OptimizationUsed     string `json:"OptimizationUsed"`
		Runs                 string `json:"Runs"`
		ConstructorArguments string `json:"ConstructorArguments"`
		EVMVersion           string `json:"EVMVersion"`
		Library              string `json:"Library"`
		LicenseType          string `json:"LicenseType"`
		Proxy                string `json:"Proxy"`
		Implementation       string `json:"Implementation"`
		SwarmSource          string `json:"SwarmSource"`
	} `json:"result"`
}

type EtherscanResponseNonStandard struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		SourceCode           string `json:"SourceCode"`
		ABI                  string `json:"ABI"`
		ContractName         string `json:"ContractName"`
		CompilerVersion      string `json:"CompilerVersion"`
		OptimizationUsed     string `json:"OptimizationUsed"`
		Runs                 string `json:"Runs"`
		ConstructorArguments string `json:"ConstructorArguments"`
		EVMVersion           string `json:"EVMVersion"`
		Library              string `json:"Library"`
		LicenseType          string `json:"LicenseType"`
		Proxy                string `json:"Proxy"`
		Implementation       string `json:"Implementation"`
		SwarmSource          string `json:"SwarmSource"`
	} `json:"result"`
}

type SourcifyResponse struct {
	Compiler struct {
		Version string `json:"version"`
	} `json:"compiler"`
	Language string `json:"language"`
	Output   struct {
		Abi     []map[string]interface{} `json:"abi"`
		Devdoc  struct{}                 `json:"devdoc"`
		Userdoc struct{}                 `json:"userdoc"`
	} `json:"output"`
	Settings struct {
		CompilationTarget map[string]string `json:"compilationTarget"`
	} `json:"settings"`
	Sources struct{} `json:"sources"`
	Version int      `json:"version"`
}

func main() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Println("No .env file found")
	}
	sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	defer sentry.Flush(time.Second * 5)

	app := pocketbase.New()

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {

		e.Router.GET(
			"/",
			func(c *core.RequestEvent) error {
				return Render(c, http.StatusOK, views.Index(ChainIDs))
			},
		)
		e.Router.GET(
			"/static/{path...}",
			apis.Static(os.DirFS("./static"), false),
		)

		e.Router.GET("/get-abi", func(c *core.RequestEvent) error {
			// Set a longer timeout for the context
			ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
			defer cancel()
			
			// Use the new context for the request
			c.Request = c.Request.WithContext(ctx)
			
			sse := datastar.NewSSE(c.Response, c.Request)

			var store GetABISignals
			if err := datastar.ReadSignals(c.Request, &store); err != nil {
				sse.ExecuteScript("console.error('Error reading signals:', " + err.Error() + ")")
				return nil
			}

			// Get ABI from database using the struct fields
			name, abi, err := getABI(app, store.ChainId, store.Address)
			if err != nil {
				sse.MergeFragmentTempl(views.Result("", nil))
				return nil
			}

			sse.MergeFragmentTempl(views.Result(name, abi))
			return nil
		})

		limiter := rate.NewLimiter(rate.Every(time.Second), 30)
		e.Router.GET("/*", func(c *core.RequestEvent) error {
			ip := c.Request.Header.Get("Fly-Client-IP")
			
			// Check banned IPs
			record, err := app.FindFirstRecordByData("bannedIPs", "ip", ip)
			if err == nil && record.Get("ip") == ip {
				return apis.NewForbiddenError("IP banned", nil)
			}

			// Rate limit check
			if !limiter.Allow() {
				log.Println("Rate limiting IP: ", ip)
				sentry.CaptureMessage("Rate limiting IP: " + ip)
				return apis.NewTooManyRequestsError("Too many requests", nil)
			}

			return c.Next()
		})

		// GET ABI
		e.Router.GET("/api/get-abi/{chainId}/{address}", func(e *core.RequestEvent) error {
				userIp := e.Request.Header.Get("Fly-Client-IP")
				log.Println("User IP: ", userIp)
				address := common.HexToAddress(e.Request.PathValue("address")).String()

				log.Println("Fetching ABI...")
				name, abi, err := getABI(app, e.Request.PathValue("chainId"), address)
				if err != nil {
					log.Println(err)
					return e.NotFoundError("ABI not found", nil)
				}
				abi = normalizeAbi(abi)
				return e.JSON(
					http.StatusOK,
					map[string]interface{}{"name": name, "abi": abi},
				)
			})

		// GET ABI .json
		e.Router.GET("/api/get-abi/{chainId}/{address}/abi.json", func(c *core.RequestEvent) error {
			address := common.HexToAddress(c.Request.PathValue("address")).String()

			_, abi, err := getABI(app, c.Request.PathValue("chainId"), address)
			if err != nil {
				return c.NotFoundError("ABI not found", nil)
			}
			abi = normalizeAbi(abi)
			return c.JSON(http.StatusOK, abi)
		})

		// POST ABI decode calldata
		e.Router.POST("/api/get-abi/{chainId}/{address}/decode", func(c *core.RequestEvent) error {
				address := common.HexToAddress(c.Request.PathValue("address")).String()

				name, abi, err := getABI(app, c.Request.PathValue("chainId"), address)
				if err != nil {
					return c.JSON(
						http.StatusNotFound,
						map[string]interface{}{"error": err.Error()},
					)
				}
				abi = normalizeAbi(abi)

				var request struct {
					CallData string `json:"calldata"`
				}

				if err := c.BindBody(&request); err != nil {
					return err
				}

				// decode txInput method signature
				decodedSig, err := hex.DecodeString(request.CallData[2:10])
				if err != nil {
					sentry.CaptureException(err)
					fmt.Fprintf(
						os.Stderr,
						"Error decoding signature: %v\n",
						err,
					)
					return err
				}

				// decode txInput Payload
				callDataArgs, err := hex.DecodeString(request.CallData[10:])
				if err != nil {
					sentry.CaptureException(err)
					fmt.Fprintf(os.Stderr, "Error decoding data: %v\n", err)
					return err
				}

				stringAbi, _ := json.Marshal(abi)

				metadata := &bind.MetaData{ABI: string(stringAbi)}
				ABI, err := metadata.GetAbi()
				if err != nil {
					sentry.CaptureException(err)
					fmt.Fprintf(os.Stderr, "Error parsing ABI: %v\n", err)
					return err
				}

				method, err := ABI.MethodById(decodedSig)
				if err != nil {
					sentry.CaptureException(err)
					fmt.Fprintf(os.Stderr, "Error finding method: %v\n", err)
					return err
				}

				data := make(map[string]interface{})
				err = method.Inputs.UnpackIntoMap(data, callDataArgs)
				if err != nil {

					sentry.CaptureException(err)
					fmt.Fprintf(os.Stderr, "Error unpacking values: %v\n", err)
					return err
				}

				return c.JSON(
					http.StatusOK,
					map[string]interface{}{
						"name": name,
						"abi":  abi,
						"args": data,
					},
				)
			},
		)

		return e.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// Render renders a templ.Component with the given status code
func Render(e *core.RequestEvent, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(e.Request.Context(), buf); err != nil {
		return err
	}

	e.Response.Header().Set("Content-Type", "text/html")
	e.Response.WriteHeader(statusCode)
	_, err := e.Response.Write(buf.Bytes())
	return err
}

func getABI(
	app *pocketbase.PocketBase,
	chainId, address string,
) (string, []map[string]interface{}, error) {
	// Try to fetch cached ABI first
	log.Println("Fetching cached ABI...")
	name, abi, err := getCachedABI(app, chainId, address)
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	if abi != nil {
		log.Println("Cached ABI found")
		return cleanName(name), abi, nil
	}

	log.Println("Fetching ABI from Etherscan...")
	name, abi, _ = getAbiFromEtherscan(chainId, address)
	if abi != nil {
		err := saveABI(app, chainId, address, name, abi)
		if err != nil {
			sentry.CaptureException(err)
			log.Println(err)
		}
		log.Println("ABI found on Etherscan")
		return cleanName(name), abi, nil
	}

	log.Println("Fetching ABI from Etherscan...")
	name, abi, _ = getAbiFromEtherscanNonStandard(chainId, address)
	if abi != nil {
		err := saveABI(app, chainId, address, name, abi)
		if err != nil {
			sentry.CaptureException(err)
			log.Println(err)
		}
		log.Println("ABI found on Etherscan")
		return cleanName(name), abi, nil
	}

	// Add Blockscout check here
	log.Println("Fetching ABI from Blockscout...")
	name, abi, _ = getAbiFromBlockscout(chainId, address)
	if abi != nil {
		err := saveABI(app, chainId, address, name, abi)
		if err != nil {
			sentry.CaptureException(err)
			log.Println(err)
		}
		log.Println("ABI found on Blockscout")
		return cleanName(name), abi, nil
	}

	log.Println("Fetching ABI from Routescan...")
	name, abi, _ = getAbiFromRoutescan(chainId, address)
	if abi != nil {
		err := saveABI(app, chainId, address, name, abi)
		if err != nil {
			sentry.CaptureException(err)
			log.Println(err)
		}
		log.Println("ABI found on Routescan")
		return cleanName(name), abi, nil
	}

	// Try to fectch ABI from Sourcify full match
	log.Println("Fetching ABI from Sourcify (full-match)...")
	name, abi, _ = getAbiFromSourcify("full", chainId, address)
	if abi != nil {
		err := saveABI(app, chainId, address, name, abi)
		if err != nil {
			sentry.CaptureException(err)
			log.Println(err)
		}
		log.Println("ABI found on Sourcify")
		return cleanName(name), abi, nil
	}

	// Try to fectch ABI from Sourcify partial match
	log.Println("Fetching ABI from Sourcify (partial-match)...")
	name, abi, _ = getAbiFromSourcify("partial", chainId, address)
	if abi != nil {
		err := saveABI(app, chainId, address, name, abi)
		if err != nil {
			sentry.CaptureException(err)
			log.Println(err)
		}
		log.Println("ABI found on Sourcify")
		return cleanName(name), abi, nil
	}

	return "", nil, errors.New("ABI not found")
}

func getCachedABI(
	app *pocketbase.PocketBase,
	chainId, address string,
) (string, []map[string]interface{}, error) {
	records, err := app.FindAllRecords("abis", dbx.NewExp("chainid = {:chainid} and address = {:address}", dbx.Params{"chainid": chainId, "address": address}))
	if err != nil || len(records) == 0 {
		return "", nil, err
	}
	abiString := records[0].GetString("abi")
	var abiJson []map[string]interface{}
	err = json.Unmarshal([]byte(abiString), &abiJson)
	if err != nil {
		sentry.CaptureException(err)
		return "", nil, err
	}
	return records[0].GetString("name"), abiJson, nil
}

func getAbiFromEtherscan(
	chainId, address string,
) (string, []map[string]interface{}, error) {
	client := &http.Client{}
	apiKey := os.Getenv("CHAIN_" + chainId + "_ETHERSCAN_KEY")
	apiUrl := fmt.Sprintf(
		"%s?module=contract&action=getsourcecode&address=%s&apikey=%s",
		etherscanConfig[chainId],
		address,
		apiKey,
	)
	log.Println(apiUrl)

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return "", nil, err
	}

	if strings.Contains(apiUrl, "oklink") {
		request.Header.Add("Ok-Access-Key", apiKey)
	}

	// Send GET request to Etherscan API
	response, err := client.Do(request)
	if err != nil {
		return "", nil, err
	}
	defer response.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", nil, err
	}

	// Unmarshal response body JSON into interface{} type
	var result EtherscanResponse
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return "", nil, err
	}

	if result.Result[0].ContractName == "" {
		return "", nil, err
	}

	// Extract ABI from interface{} type
	var abiJson []map[string]interface{}
	err = json.Unmarshal([]byte(result.Result[0].ABI), &abiJson)
	if err != nil {
		sentry.CaptureException(err)
		return "", nil, err
	}

	return result.Result[0].ContractName, abiJson, nil
}

func getAbiFromEtherscanNonStandard(
	chainId, address string,
) (string, []map[string]interface{}, error) {
	apiUrl := fmt.Sprintf(
		"%s?module=contract&action=getsourcecode&address=%s",
		etherscanConfig[chainId],
		address,
	)
	log.Println(apiUrl)

	// Send GET request to Etherscan API
	response, err := http.Get(apiUrl)
	if err != nil {
		return "", nil, err
	}
	defer response.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", nil, err
	}

	// Unmarshal response body JSON into interface{} type
	var result EtherscanResponseNonStandard
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return "", nil, err
	}

	if result.Result.ContractName == "" {
		return "", nil, err
	}

	// Extract ABI from interface{} type
	var abiJson []map[string]interface{}
	err = json.Unmarshal([]byte(result.Result.ABI), &abiJson)
	if err != nil {
		return "", nil, err
	}

	return result.Result.ContractName, abiJson, nil
}

func getAbiFromRoutescan(
	chainId, address string,
) (string, []map[string]interface{}, error) {
	routeScanUrl := fmt.Sprintf(
		"https://api.routescan.io/v2/network/mainnet/evm/%s/etherscan/api",
		chainId,
	)
	apiUrl := fmt.Sprintf(
		"%s?module=contract&action=getsourcecode&address=%s&apikey=%s",
		routeScanUrl,
		address,
		os.Getenv("CHAIN_"+chainId+"_ETHERSCAN_KEY"),
	)

	// Send GET request to Etherscan API
	response, err := http.Get(apiUrl)
	if err != nil {
		return "", nil, err
	}
	defer response.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", nil, err
	}

	// Unmarshal response body JSON into interface{} type
	var result EtherscanResponse
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return "", nil, err
	}

	if len(result.Result) < 1 {
		return "", nil, errors.New("not found")
	}

	// Extract ABI from interface{} type
	var abiJson []map[string]interface{}
	err = json.Unmarshal([]byte(result.Result[0].ABI), &abiJson)
	if err != nil {
		return "", nil, err
	}

	return result.Result[0].ContractName, abiJson, nil
}

func getAbiFromSourcify(
	matchType, chainId, address string,
) (string, []map[string]interface{}, error) {
	if matchType != "full" && matchType != "partial" {
		return "", nil, fmt.Errorf("invalid type")
	}

	// Replace <API_KEY> w0xF2ee649caB7a0edEdED7a27821B0aCDF77778aeDith your Etherscan API key
	apiUrl := fmt.Sprintf(
		"https://repo.sourcify.dev/contracts/%s_match/%s/%s/metadata.json",
		matchType,
		chainId,
		address,
	)

	// Send GET request to Etherscan API
	response, err := http.Get(apiUrl)
	if err != nil {
		return "", nil, err
	}
	defer response.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", nil, err
	}

	// Unmarshal response body JSON into interface{} type
	var result SourcifyResponse
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return "", nil, err
	}

	var contractName string
	for _, contract := range result.Settings.CompilationTarget {
		contractName = contract
		break
	}

	return contractName, result.Output.Abi, nil
}

func saveABI(
	app *pocketbase.PocketBase,
	chainid, address, name string,
	abi []map[string]interface{},
) error {
	collection, err := app.FindCollectionByNameOrId("abis")
	if err != nil {
		return err
	}

	record := core.NewRecord(collection)
	record.Set("chainId", chainid)
	record.Set("address", address) 
	record.Set("name", name)
	record.Set("abi", abi)

	if err := app.Save(record); err != nil {
		return err
	}

	return nil
}

func cleanName(name string) string {
	if strings.Contains(name, ":") {
		return strings.Split(name, ":")[1]
	} else {
		return name
	}
}

func getAbiFromBlockscout(
	chainId, address string,
) (string, []map[string]interface{}, error) {
	// Check if we have a Blockscout API for this chain
	apiUrl, exists := blockscoutConfig[chainId]
	if !exists {
		return "", nil, fmt.Errorf("no blockscout API for chain %s", chainId)
	}

	url := fmt.Sprintf(
		"%s?module=contract&action=getsourcecode&address=%s",
		apiUrl,
		address,
	)
	log.Println(url)

	// Send GET request to Blockscout API
	response, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}
	defer response.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", nil, err
	}

	// Try standard format first
	var result EtherscanResponse
	err = json.Unmarshal(responseBody, &result)
	if err == nil && len(result.Result) > 0 && result.Result[0].ABI != "" {
		var abiJson []map[string]interface{}
		err = json.Unmarshal([]byte(result.Result[0].ABI), &abiJson)
		if err == nil {
			return result.Result[0].ContractName, abiJson, nil
		}
	}

	// Try non-standard format
	var nonStandardResult EtherscanResponseNonStandard
	err = json.Unmarshal(responseBody, &nonStandardResult)
	if err == nil && nonStandardResult.Result.ABI != "" {
		var abiJson []map[string]interface{}
		err = json.Unmarshal([]byte(nonStandardResult.Result.ABI), &abiJson)
		if err == nil {
			return nonStandardResult.Result.ContractName, abiJson, nil
		}
	}

	return "", nil, fmt.Errorf("unable to get ABI from Blockscout")
}

func normalizeAbi(abi []map[string]interface{}) []map[string]interface{} {
	newAbi := []map[string]interface{}{}
	// loop through each item in the array and if "type" == "function" make sure any "outputs" parameter exists. If not set "outputs" to []
	for _, item := range abi {
		if item["type"] == "function" {
			if _, ok := item["outputs"]; !ok {
				item["outputs"] = []map[string]interface{}{}
			}
		}
		newAbi = append(newAbi, item)
	}

	return newAbi
}
