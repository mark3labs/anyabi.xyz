<script lang="ts">
  import chains from '$lib/chains.json'
  import Snippets from '$lib/components/Snippets.svelte'
  import hljs from 'highlight.js'
  import 'highlight.js/styles/github-dark.css'
  import { CodeBlock, storeHighlightJs } from '@skeletonlabs/skeleton'
  import { isAddress } from 'ethers'
  import { Pulse } from 'svelte-loading-spinners'
  import Typeahead from 'svelte-typeahead'
  import { PUBLIC_API_URL } from '$env/static/public'

  storeHighlightJs.set(hljs)

  let abi: any

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
  <Snippets />
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
