//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"log"
	"os"
	"text/template"

	"github.com/mark3labs/anyabi.xyz/types"
)

func main() {
	// open file
	f, err := os.Create("./apis_generated.go")
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	t := template.Must(template.New("").
		Parse(`// Code generated by go generate. DO NOT EDIT.
package main

var etherscanConfig map[string]string = map[string]string{
{{ range $k, $v := . }}{{ if $v.BlockExplorers.Default.APIURL }}"{{ $v.ID }}":"{{ $v.BlockExplorers.Default.APIURL }}",
    {{end}}{{ end }}
}`))

	file, err := os.Open("./codegen/chains.json") // the file you want to read
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var chains types.Chains

	err = json.NewDecoder(file).Decode(&chains)
	if err != nil {
		log.Fatal(err)
	}

	// assign a value to the placeholder and write to file
	err = t.Execute(f, chains)
	if err != nil {
		log.Print("execute: ", err)
		return
	}

	f.Close()
}
