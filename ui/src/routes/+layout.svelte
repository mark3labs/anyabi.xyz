<script>
  import '@skeletonlabs/skeleton/themes/theme-modern.css'
  import '@skeletonlabs/skeleton/styles/all.css'
  import '../app.postcss'
  import Icon from '@iconify/svelte'
  import { AppShell, AppBar, storePopup } from '@skeletonlabs/skeleton'
  import { computePosition, autoUpdate, flip, shift, offset, arrow } from '@floating-ui/dom'

  storePopup.set({ computePosition, autoUpdate, flip, shift, offset, arrow })
</script>

<AppShell>
  <svelte:fragment slot="header">
    <AppBar>
      <svelte:fragment slot="lead">
        <h2>
          <img src="/logo.png" alt="AnyABI" class="inline h-10" />
          AnyABI
        </h2>
      </svelte:fragment>
      <svelte:fragment slot="trail">
        <a href="https://github.com/mark3labs/anyabi.xyz" target="_blank" rel="noreferrer">
          <Icon icon="fe:github-alt" color="white" width="32" />
        </a>
      </svelte:fragment>
    </AppBar>
  </svelte:fragment>
  <header
    class="variant-glass-surface border-b border-black/5 dark:border-white/5 px-4 py-10 md:py-20"
  >
    <div
      class="container max-w-[1200px] mx-auto grid grid-cols-1 md:grid-cols-2 gap-10 place-items-center"
    >
      <div class="space-y-4">
        <div>
          <img src="/logo.png" alt="AnyABI" class="inline w-24 md:w-48" />
          <span class="text-2xl md:text-7xl font-bold">AnyABI</span>
        </div>
        <h2 class="text-md md:text-4xl">The easiest way to fetch an ABI.</h2>
        <p class="md:!text-2xl">
          Quickly grab the ABI for <span class="font-bold">ANY<sup>*</sup></span> EVM smart contract
          on
          <span class="font-bold">ANY</span> EVM chain.
        </p>
        <div class="text-xs mt-5">
          * ABIs can only be retrieved for contracts that are verified on a block explorer like
          Etherscan or on Sourcify
        </div>
      </div>
    </div>
  </header>
  <!-- Router Slot -->
  <slot />
  <!-- ---- / ---- -->
  <!-- <svelte:fragment slot="footer">Footer</svelte:fragment> -->
</AppShell>
