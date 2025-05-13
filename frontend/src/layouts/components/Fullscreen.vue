<template>
  <i
    id="fullscreen"
    class="mr-16 cursor-pointer text-20"
    :class="isFullscreen ? 'i-fe:minimize' : 'i-fe:maximize'"
    @click="toggle"
  />
</template>

<script setup>
import { useFullscreen } from '@vueuse/core'

const { isFullscreen, toggle } = useFullscreen()
</script>
