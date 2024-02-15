// main.go
package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	_ "github.com/mark3labs/anyabi.xyz/migrations"
	"github.com/mark3labs/anyabi.xyz/ui"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

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

var etherscanConfig map[string]string = map[string]string{
	"1":         "https://api.etherscan.io/api",
	"5":         "https://api-goerli.etherscan.io/api",
	"11155111":  "https://api-sepolia.etherscan.io/api",
	"100":       "https://api.gnosisscan.io/api",
	"137":       "https://api.polygonscan.com/api",
	"80001":     "https://api-testnet.polygonscan.com/api",
	"56":        "https://api.bscscan.com/api",
	"97":        "https://api-testnet.bscscan.com/api",
	"43114":     "https://api.snowtrace.io/api",
	"43113":     "https://api-testnet.snowtrace.io/api",
	"10":        "https://api-optimistic.etherscan.io/api",
	"420":       "https://api-goerli-optimistic.etherscan.io/api",
	"42161":     "https://api.arbiscan.io/api",
	"421613":    "https://api-goerli.arbiscan.io/api",
	"421614":    "https://api-sepolia.arbiscan.io/api",
	"42170":     "https://api-nova.arbiscan.io/api",
	"250":       "https://api.ftmscan.com/api",
	"4002":      "https://api-testnet.ftmscan.com/api",
	"1284":      "https://api-moonbeam.moonscan.io/api",
	"1287":      "https://api-moonbase.moonscan.io/api",
	"1285":      "https://api-moonriver.moonscan.io/api",
	"25":        "https://api.cronoscan.com/api",
	"338":       "https://api-testnet.cronoscan.com/api",
	"42220":     "https://api.celoscan.io/api",
	"44787":     "https://api-alfajores.celoscan.io/api",
	"288":       "https://api.bobascan.com/api",
	"2888":      "https://api-testnet.bobascan.com/api",
	"534353":    "https://blockscout.scroll.io/api",
	"8453":      "https://api.basescan.org/api",
	"84531":     "https://api-goerli.basescan.org/api",
	"84532":     "https://api-sepolia.basescan.org/api",
	"59144":     "https://api.lineascan.build/api",
	"324":       "https://block-explorer-api.mainnet.zksync.io/api",
	"1101":      "https://api-zkevm.polygonscan.com/api",
	"168587773": "https://api.routescan.io/v2/network/testnet/evm/168587773/etherscan/api",
	"1088":      "https://api.routescan.io/v2/network/mainnet/evm/1088/etherscan/api",

	// TODO: finsish adding chains
}

func main() {
	godotenv.Load()
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// GET ABI
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/get-abi/:chainId/:address",
			Handler: func(c echo.Context) error {
				address := common.HexToAddress(c.PathParam("address")).String()

				log.Println("Fetching ABI...")
				name, abi, err := getABI(app, c.PathParam("chainId"), address)
				if err != nil {
					log.Println(err)
					return c.JSON(http.StatusNotFound, map[string]interface{}{"error": err.Error()})
				}
				abi = normalizeAbi(abi)
				return c.JSON(http.StatusOK, map[string]interface{}{"name": name, "abi": abi})
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		// GET ABI .json
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/get-abi/:chainId/:address/abi.json",
			Handler: func(c echo.Context) error {
				address := common.HexToAddress(c.PathParam("address")).String()

				_, abi, err := getABI(app, c.PathParam("chainId"), address)
				if err != nil {
					return c.JSON(http.StatusNotFound, map[string]interface{}{"error": err.Error()})
				}
				abi = normalizeAbi(abi)
				return c.JSON(http.StatusOK, abi)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		// POST ABI decode calldata
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/get-abi/:chainId/:address/decode",
			Handler: func(c echo.Context) error {
				address := common.HexToAddress(c.PathParam("address")).String()

				name, abi, err := getABI(app, c.PathParam("chainId"), address)
				if err != nil {
					return c.JSON(http.StatusNotFound, map[string]interface{}{"error": err.Error()})
				}
				abi = normalizeAbi(abi)

				var request struct {
					CallData string `json:"calldata"`
				}

				c.Bind(&request)

				// decode txInput method signature
				decodedSig, err := hex.DecodeString(request.CallData[2:10])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error decoding signature: %v\n", err)
					return err
				}

				// decode txInput Payload
				callDataArgs, err := hex.DecodeString(request.CallData[10:])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error decoding data: %v\n", err)
					return err
				}

				stringAbi, _ := json.Marshal(abi)

				metadata := &bind.MetaData{ABI: string(stringAbi)}
				ABI, err := metadata.GetAbi()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing ABI: %v\n", err)
					return err
				}

				method, err := ABI.MethodById(decodedSig)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error finding method: %v\n", err)
					return err
				}

				data := make(map[string]interface{})
				err = method.Inputs.UnpackIntoMap(data, callDataArgs)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error unpacking values: %v\n", err)
					return err
				}

				return c.JSON(http.StatusOK, map[string]interface{}{"name": name, "abi": abi, "args": data})
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		e.Router.GET("/*", apis.StaticDirectoryHandler(ui.BuildDirFS, true))

		return nil
	})

	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: true, // auto creates migration files when making collection changes
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func getABI(app *pocketbase.PocketBase, chainId, address string) (string, []map[string]interface{}, error) {
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
			log.Println(err)
		}
		log.Println("ABI found on Etherscan")
		return cleanName(name), abi, nil
	}

	// Try to fectch ABI from Sourcify full match
	log.Println("Fetching ABI from Sourcify (full-match)...")
	name, abi, _ = getAbiFromSourcify("full", chainId, address)
	if abi != nil {
		err := saveABI(app, chainId, address, name, abi)
		if err != nil {
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
			log.Println(err)
		}
		log.Println("ABI found on Sourcify")
		return cleanName(name), abi, nil
	}

	return "", nil, errors.New("ABI not found")
}

func getCachedABI(app *pocketbase.PocketBase, chainId, address string) (string, []map[string]interface{}, error) {
	records, err := app.Dao().FindRecordsByExpr("abis", dbx.NewExp("chainid = {:chainid} and address = {:address}", dbx.Params{"chainid": chainId, "address": address}))
	if err != nil || len(records) == 0 {
		return "", nil, err
	}
	abiString := records[0].GetString("abi")
	var abiJson []map[string]interface{}
	err = json.Unmarshal([]byte(abiString), &abiJson)
	if err != nil {
		return "", nil, err
	}
	return records[0].GetString("name"), abiJson, nil
}

func getAbiFromEtherscan(chainId, address string) (string, []map[string]interface{}, error) {
	apiUrl := fmt.Sprintf("%s?module=contract&action=getsourcecode&address=%s&apikey=%s", etherscanConfig[chainId], address, os.Getenv("CHAIN_"+chainId+"_ETHERSCAN_KEY"))

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

	// Extract ABI from interface{} type
	var abiJson []map[string]interface{}
	err = json.Unmarshal([]byte(result.Result[0].ABI), &abiJson)
	if err != nil {
		return "", nil, err
	}

	return result.Result[0].ContractName, abiJson, nil
}

func getAbiFromSourcify(matchType, chainId, address string) (string, []map[string]interface{}, error) {
	if matchType != "full" && matchType != "partial" {
		return "", nil, fmt.Errorf("invalid type")
	}

	// Replace <API_KEY> with your Etherscan API key
	apiUrl := fmt.Sprintf("https://repo.sourcify.dev/contracts/%s_match/%s/%s/metadata.json", matchType, chainId, address)

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

func saveABI(app *pocketbase.PocketBase, chainid, address, name string, abi []map[string]interface{}) error {
	collection, err := app.Dao().FindCollectionByNameOrId("abis")
	if err != nil {
		return err
	}

	record := models.NewRecord(collection)
	form := forms.NewRecordUpsert(app, record)
	form.LoadData(map[string]interface{}{
		"chainId": chainid,
		"address": address,
		"name":    name,
		"abi":     abi,
	})

	if err := form.Submit(); err != nil {
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
