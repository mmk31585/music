<script setup lang="ts">
import { useRouter } from 'vue-router'
import { computed } from 'vue'

const router = useRouter()

const homeRoute = computed(() => '/')
</script>

<template>
  <section class="bg-indigo-500 relative z-10 py-[120px] min-h-screen">
    <div class="max-w-7xl mx-auto w-full">
      <div class="flex">
        <div class="w-full px-4">
          <div class="mx-auto max-w-[400px] text-center">
            <h2
              class="mb-2 text-[50px] font-bold leading-none text-white sm:text-[80px] md:text-[100px]"
            >
              ۴۰۴
            </h2>
            <h4 class="mb-3 text-[22px] font-semibold leading-tight text-white">
              صفحه مورد نظر شما پیدا نشد!
            </h4>
            <p class="mb-8 text-lg text-white">لطفا دوباره تلاش نمایید</p>

            <div class="flex gap-3 justify-center items-center">
              <router-link
                v-if="homeRoute"
                :to="homeRoute"
                class="inline-block text-purple-500 rounded-lg border border-white px-8 py-3 text-center text-base font-semibold transition bg-white hover:text-black hover:bg-opacity-80 hover:shadow-xl"
              >
                خانه
              </router-link>
              <button
                class="inline-block rounded-lg border border-white px-8 py-3 text-center text-base font-semibold text-white transition hover:bg-white hover:text-purple-500 hover:shadow-xl"
                @click="router.go(-1)"
              >
                بازگشت
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div
      class="absolute top-0 left-0 -z-10 flex h-full w-full items-center justify-between space-x-5 md:space-x-8 lg:space-x-14"
    >
      <div
        class="h-full w-1/3 bg-linear-to-t from-[#FFFFFF14] to-[#C4C4C400]"
      ></div>
      <div class="flex h-full w-1/3">
        <div
          class="h-full w-1/2 bg-linear-to-b from-[#FFFFFF14] to-[#C4C4C400]"
        ></div>
        <div
          class="h-full w-1/2 bg-gradient-to-t from-[#FFFFFF14] to-[#C4C4C400]"
        ></div>
      </div>
      <div
        class="h-full w-1/3 bg-gradient-to-b from-[#FFFFFF14] to-[#C4C4C400]"
      ></div>
    </div>
  </section>
</template>
