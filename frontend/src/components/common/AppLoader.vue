import { Loader2 } from 'lucide-vue-next'
<template>
  <div
    class="flex items-center justify-center py-10 text-slate-400"
    role="status"
    aria-live="polite"
    aria-busy="true"
  >
    <Loader2 aria-hidden="true" class="mr-2 animate-spin"  />
    <span>Loading...</span>
  </div>
</template>
