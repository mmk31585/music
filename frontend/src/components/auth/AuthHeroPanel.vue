<template>
  <div
    class="relative hidden min-h-screen overflow-hidden bg-[#0a0a0a] p-10 text-white lg:flex lg:flex-col lg:justify-between"
  >
    <!-- Aurora glow background -->
    <div class="absolute inset-0 opacity-30">
      <div class="aurora-glow absolute -top-20 -left-20 h-80 w-80 rounded-full bg-[#1db954] blur-3xl" />
      <div class="aurora-glow absolute top-1/3 right-10 h-60 w-60 rounded-full bg-[#a855f7] blur-3xl" style="animation-delay: 2s" />
      <div class="aurora-glow absolute -bottom-10 left-1/4 h-72 w-72 rounded-full bg-[#60a5fa] blur-3xl" style="animation-delay: 4s" />
    </div>

    <!-- Equalizer bars -->
    <div class="absolute right-12 top-1/3 z-10 flex items-end gap-[3px] opacity-40">
      <div v-for="i in 24" :key="i" class="equalizer-bar w-[3px] rounded-full bg-[#1db954]" />
    </div>

    <!-- Content -->
    <div class="relative z-10 flex flex-col justify-between h-full">
      <!-- Brand -->
      <div class="flex items-center gap-3">
        <div
          class="flex h-11 w-11 items-center justify-center rounded-xl bg-[#1db954] text-lg font-extrabold text-black shadow-lg shadow-[#1db954]/30"
        >
          M
        </div>
        <h1 class="text-xl font-bold tracking-tight">Muse</h1>
      </div>

      <!-- Hero text -->
      <div class="pb-20">
        <h2 class="font-display max-w-md text-5xl leading-[1.1] font-extrabold tracking-tight">
          <span class="text-gradient">Feel every beat.</span>
        </h2>
        <p class="mt-4 max-w-sm text-base text-white/50 leading-relaxed">
          Discover tracks, artists, albums, and genres in a smooth streaming experience.
        </p>
        <p class="mt-2 text-sm text-white/30" dir="rtl">
          هر نت یک داستان، هر آهنگ یک خاطره
        </p>
      </div>

      <!-- Footer -->
      <div class="flex items-center gap-2 text-xs text-white/30">
        <span class="inline-block h-1 w-1 rounded-full bg-[#1db954]" />
        <span>© 2026 Muse — Music streaming reimagined</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.text-gradient {
  background: linear-gradient(135deg, #1db954 0%, #60a5fa 50%, #a855f7 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.aurora-glow {
  animation: aurora-pulse 6s ease-in-out infinite alternate;
}

@keyframes aurora-pulse {
  0% {
    transform: scale(1);
    opacity: 0.4;
  }
  50% {
    transform: scale(1.15);
    opacity: 0.6;
  }
  100% {
    transform: scale(1);
    opacity: 0.4;
  }
}

.equalizer-bar {
  animation: equalizer-bounce 0.8s ease-in-out infinite alternate;
  animation-delay: calc(var(--i, 0) * 0.05s);
  height: calc(8px + (var(--bar-height, 20) * 1px));
}

.equalizer-bar:nth-child(1) { --bar-height: 12; --i: 0; }
.equalizer-bar:nth-child(2) { --bar-height: 28; --i: 1; }
.equalizer-bar:nth-child(3) { --bar-height: 8; --i: 2; }
.equalizer-bar:nth-child(4) { --bar-height: 36; --i: 3; }
.equalizer-bar:nth-child(5) { --bar-height: 18; --i: 4; }
.equalizer-bar:nth-child(6) { --bar-height: 42; --i: 5; }
.equalizer-bar:nth-child(7) { --bar-height: 14; --i: 6; }
.equalizer-bar:nth-child(8) { --bar-height: 32; --i: 7; }
.equalizer-bar:nth-child(9) { --bar-height: 24; --i: 8; }
.equalizer-bar:nth-child(10) { --bar-height: 38; --i: 9; }
.equalizer-bar:nth-child(11) { --bar-height: 10; --i: 10; }
.equalizer-bar:nth-child(12) { --bar-height: 30; --i: 11; }
.equalizer-bar:nth-child(13) { --bar-height: 22; --i: 12; }
.equalizer-bar:nth-child(14) { --bar-height: 40; --i: 13; }
.equalizer-bar:nth-child(15) { --bar-height: 16; --i: 14; }
.equalizer-bar:nth-child(16) { --bar-height: 34; --i: 15; }
.equalizer-bar:nth-child(17) { --bar-height: 26; --i: 16; }
.equalizer-bar:nth-child(18) { --bar-height: 44; --i: 17; }
.equalizer-bar:nth-child(19) { --bar-height: 10; --i: 18; }
.equalizer-bar:nth-child(20) { --bar-height: 28; --i: 19; }
.equalizer-bar:nth-child(21) { --bar-height: 38; --i: 20; }
.equalizer-bar:nth-child(22) { --bar-height: 14; --i: 21; }
.equalizer-bar:nth-child(23) { --bar-height: 32; --i: 22; }
.equalizer-bar:nth-child(24) { --bar-height: 20; --i: 23; }

@keyframes equalizer-bounce {
  0% {
    transform: scaleY(0.4);
    opacity: 0.5;
  }
  100% {
    transform: scaleY(1);
    opacity: 1;
  }
}
</style>
