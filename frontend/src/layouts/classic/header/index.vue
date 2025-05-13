<template>
  <AppCard class="flex items-center px-12" border-b="1px solid light_border dark:dark_border">
    <HeaderLogo border-b="1px solid light_border dark:dark_border" :class="`${mlpx}`" />

    <AppTab class="w-0 flex-1 px-12" />

    <span class="opacity-1">|</span>
  </AppCard>
</template>

<script setup>
import { AppTab, HeaderLogo } from '@/layouts/components'
import { isMacOS } from '@/utils'
import { EventsOn, WindowIsFullscreen } from 'wailsjs/runtime/runtime.js'

const mlpx = ref('ml-40')
function onToggleFullscreen(fullscreen) {
  if (fullscreen) {
    mlpx.value = 'ml-0'
  }
  else {
    mlpx.value = isMacOS() ? 'ml-60' : 'ml-0'
  }
}

EventsOn('window_changed', (info) => {
  console.warn('window_changed event received:', info)
  const { fullscreen } = info
  onToggleFullscreen(fullscreen === true)
})

onMounted(async () => {
  const fullscreen = await WindowIsFullscreen()
  onToggleFullscreen(fullscreen === true)
})
</script>
