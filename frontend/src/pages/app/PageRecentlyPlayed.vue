<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <section
      class="rounded-[2rem] bg-gradient-to-br from-sky-500 via-slate-900 to-black p-8 text-white"
    >
      <p class="text-sm font-bold tracking-[0.35em] text-white/70 uppercase">History</p>

      <h1 class="mt-3 text-4xl font-black md:text-6xl">Recently Played</h1>

      <p class="mt-4 max-w-2xl text-white/80">Jump back into tracks you played recently.</p>
    </section>

    <section class="mt-10 rounded-3xl border border-white/10 bg-black/20 px-6 py-16 text-center">
      <div
        class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-white/10 text-2xl text-white"
      >
        <i class="pi pi-history" />
      </div>

      <h2 class="mt-5 text-xl font-black text-white">No recent plays yet</h2>

      <p class="mt-2 text-sm text-slate-400">
        Start playing tracks and your history will appear here.
      </p>
    </section>
  </div>
</template>
