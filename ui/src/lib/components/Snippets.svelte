<script>
  import { CodeBlock, Tab, TabGroup } from '@skeletonlabs/skeleton'

  let tabSet = 0
</script>

<h3>Use the API</h3>
<p class="my-5">Use the API to fetch the ABI of any contract on any chain.</p>
<TabGroup>
  <Tab bind:group={tabSet} name="tab1" value={0}>Curl</Tab>
  <Tab bind:group={tabSet} name="tab2" value={1}>Typescript</Tab>
  <Tab bind:group={tabSet} name="tab3" value={2}>Go</Tab>
  <Tab bind:group={tabSet} name="tab3" value={3}>Python</Tab>
  <Tab bind:group={tabSet} name="tab3" value={4}>Ruby</Tab>
  <Tab bind:group={tabSet} name="tab3" value={5}>Rust</Tab>
  <!-- Tab Panels --->
  <svelte:fragment slot="panel">
    {#if tabSet === 0}
      <CodeBlock language="bash" code="curl https://anyabi.xyz/api/get-abi/<chainid>/<address>" />
    {:else if tabSet === 1}
      <CodeBlock
        language="typescript"
        code={`
const chainId = "1";
const address = "0x0000000000000000000000000000000000000000"; // replace with contract address

const url = \`https://anyabi.xyz/api/get-abi/\${chainId}/\${address}\`;

fetch(url)
  .then(response => response.json())
  .then(data => console.log(data));
`}
      />
    {:else if tabSet === 2}
      <CodeBlock
        language="go"
        code={`
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	chainid := "1"
	address := "0x0000000000000000000000000000000000000000" // replace with contract address
	req, err := http.NewRequest("GET", fmt.Sprintf("https://anyabi.xyz/api/get-abi/%v/%v", chainid, address), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
`}
      />
    {:else if tabSet === 3}
      <CodeBlock
        language="python"
        code={`
import requests

chainid = "1"
address = "0x0000000000000000000000000000000000000000" # replace with contract address

url = "https://anyabi.xyz/api/get-abi/" + chainid + "/" + address

response = requests.get(url)

print(response.json())
`}
      />
    {:else if tabSet === 4}
      <CodeBlock
        language="ruby"
        code={`
require 'net/http'
require 'json'

chainid = "1"
address = "0x0000000000000000000000000000000000000000" # replace with contract address

url = "https://anyabi.xyz/api/get-abi/#{chainid}/#{address}"
uri = URI(url)
response = Net::HTTP.get(uri)
puts "#{response.body}"
`}
      />
    {:else if tabSet === 5}
      <CodeBlock
        language="rust"
        code={`
use reqwest;

fn main() {
    let chain_id = "1";
    let address = "0x0000000000000000000000000000000000000000"; // replace with contract address

    let url = format!("https://anyabi.xyz/api/get-abi/{}/{}", chain_id, address);

    let body = reqwest::get(&url).unwrap().text().unwrap();

    println!("{}", body);
}
`}
      />
    {/if}
  </svelte:fragment>
</TabGroup>
<hr class="my-5" />
