<script lang="ts">
  import chains from '$lib/chains.json'
  import hljs from 'highlight.js'
  import 'highlight.js/styles/github-dark.css'
  import { CodeBlock, storeHighlightJs, Tab, TabGroup } from '@skeletonlabs/skeleton'
  import { isAddress } from 'ethers'
  import { Pulse } from 'svelte-loading-spinners'
  import Typeahead from 'svelte-typeahead'
  import { PUBLIC_API_URL } from '$env/static/public'

  storeHighlightJs.set(hljs)

  let abi: any
  let tabSet = 0

  let form = {
    chainId: 1,
    address: '',
  }

  const fetchABI = async () => {
    abi = undefined
    if (!isAddress(form.address)) {
      throw new Error('Invalid Address')
    }
    const res = await fetch(`${PUBLIC_API_URL}/api/get-abi/${form.chainId}/${form.address}`)
    if (res.ok) {
      return await res.json()
    } else {
      throw new Error('ABI Not Found')
    }
  }
</script>

<div class="container mx-auto my-20 p-2">
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
          code={"const response = await fetch(`https://anyabi.xyz/api/get-abi/${chainid}/${address}`, {method: 'GET'});"}
        />
      {:else if tabSet === 2}
        <CodeBlock
          language="go"
          code={`req, err := http.NewRequest("GET", fmt.Sprintf("https://anyabi.xyz/api/get-abi/%v/%v", chainid, address), nil)`}
        />
      {:else if tabSet === 3}
        <CodeBlock
          language="python"
          code={`response = requests.get(f"https://anyabi.xyz/api/get-abi/{chainid}/{address}")`}
        />
      {:else if tabSet === 4}
        <CodeBlock
          language="ruby"
          code={`response = Net::HTTP.get_response(URI(f"https://anyabi.xyz/api/get-abi/{chainid}/{address}"))`}
        />
      {:else if tabSet === 5}
        <CodeBlock
          language="rust"
          code={`let response = reqwest::get(&format!("https://anyabi.xyz/api/get-abi/{}/{}", chainid, address))?;`}
        />
      {/if}
    </svelte:fragment>
  </TabGroup>
  <hr class="my-5" />
  <h3 class="mb-3">Search for an ABI</h3>
  <Typeahead
    hideLabel
    showDropdownOnFocus
    showAllResultsOnFocus
    data={chains}
    extract={(item) => item.name}
    on:select={({ detail }) => (form.chainId = detail.original.chainId)}
    placeholder="Chain..."
  />
  <div class="input-group input-group-divider grid-cols-[1fr_auto]">
    <input type="search" placeholder="Contract Address..." bind:value={form.address} />
    <button
      class="variant-glass-primary hover:variant-filled-primary"
      on:click={() => (abi = fetchABI())}>Let's Go!</button
    >
  </div>
  {#if abi}
    <div>
      {#await abi}
        <div class="flex justify-center mt-20 p-5">
          <Pulse color="#EC4899" />
        </div>
      {:then result}
        <div class="mt-10">
          <h3 class="mb-3">{result.name}</h3>
          <CodeBlock language="json" code={JSON.stringify(result.abi, null, 2)} />
        </div>
      {:catch e}
        <div class="mt-10">
          <aside class="alert variant-ghost-error">
            <!-- Message -->
            <div class="alert-message text-center">
              <h3>{e.message}</h3>
            </div>
          </aside>
        </div>
      {/await}
    </div>
  {/if}
</div>

<style lang="postcss">
  :global([data-svelte-typeahead]) {
    @apply !rounded-full !bg-surface-500 !text-primary-500 !border-0 !shadow-none !ring-0 !outline-none !m-1;
  }

  :global([data-svelte-search] input) {
    @apply !rounded-full !border-0 !ring !ring-surface-500 !bg-surface-700 focus:!ring-primary-500 focus:!outline-none;
  }

  :global([data-svelte-typeahead] ul) {
    @apply !bg-surface-500 rounded !text-white;
  }

  :global([data-svelte-typeahead] mark) {
    @apply !bg-surface-200;
  }

  :global([data-svelte-typeahead] li.selected) {
    @apply !bg-primary-500;
  }

  :global([data-svelte-search]) {
    @apply rounded-full mb-3 text-white;
  }
</style>
