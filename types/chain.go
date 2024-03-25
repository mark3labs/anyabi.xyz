package types

type Chain struct {
	BlockExplorers struct {
		Default struct {
			Name   string `json:"name"`
			URL    string `json:"url"`
			APIURL string `json:"apiUrl"`
		} `json:"default"`
	} `json:"blockExplorers"`
	Name    string `json:"name"`
	Network string `json:"network"`
	RPCUrls struct {
		Public struct {
			HTTP      []string `json:"http"`
			WebSocket []string `json:"webSocket"`
		} `json:"public"`
		Default struct {
			HTTP      []string `json:"http"`
			WebSocket []string `json:"webSocket"`
		} `json:"default"`
	} `json:"rpcUrls"`
	NativeCurrency struct {
		Name     string `json:"name"`
		Symbol   string `json:"symbol"`
		Decimals int    `json:"decimals"`
	} `json:"nativeCurrency"`
	ID      int  `json:"id"`
	Testnet bool `json:"testnet"`
}

type Chains map[string]Chain
