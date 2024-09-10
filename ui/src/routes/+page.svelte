<script lang="ts">
  import chains from '$lib/chains.json'
  import Snippets from '$lib/components/Snippets.svelte'
  import hljs from 'highlight.js'
  import 'highlight.js/styles/github-dark.css'
  import { CodeBlock, storeHighlightJs, Autocomplete, popup } from '@skeletonlabs/skeleton'
  import type { AutocompleteOption, PopupSettings } from '@skeletonlabs/skeleton'
  import { isAddress } from 'ethers'
  import { Pulse } from 'svelte-loading-spinners'

  const PUBLIC_API_URL = 'https://anyabi.xyz' // TODO update UI

  storeHighlightJs.set(hljs)

  let abi: any

  let form = {
    chainId: null,
    address: '',
  }

  let chainSearch = ''

  const options: AutocompleteOption[] = chains.map((chain) => ({
    label: chain.name,
    value: chain.chainId,
  }))

  let popupSettings: PopupSettings = {
    event: 'focus',
    target: 'popupAutocomplete',
    placement: 'top',
  }

  const onChainSelect = (event: any) => {
    form.chainId = event.detail.value
    chainSearch = event.detail.label
  }

  const fetchABI = async () => {
    abi = undefined

    if (!form.chainId) {
      throw new Error('No Chain Selected')
    }

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
  <input
    class="input autocomplete mb-2"
    type="search"
    name="chain"
    bind:value={chainSearch}
    placeholder="Chain..."
    use:popup={popupSettings}
  />
  <div data-popup="popupAutocomplete" class="card w-full overflow-y-auto max-h-48 max-w-lg">
    <Autocomplete bind:input={chainSearch} {options} on:selection={onChainSelect} />
  </div>
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
