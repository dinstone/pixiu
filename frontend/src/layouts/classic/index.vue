<template>
  <div class="wh-full flex flex-col">
    <AppHeader class="h-60 flex-shrink-0" style="--wails-draggable: drag"/>

    <div class="flex flex-1 overflow-hidden">
      <aside
        class="flex-col flex-shrink-0 transition-width-300"
        :class="appStore.collapsed ? 'w-60' : 'w-140'"
        border-r="1px solid light_border dark:dark_border"
        style="--wails-draggable: drag"
      >
        <SideBar />
      </aside>

      <article class="w-full flex-col flex-1 overflow-y-auto">
        <slot />
      </article>
    </div>
  </div>
</template>

<script setup>
import { useAppStore } from '@/store'
import AppHeader from './header/index.vue'
import SideBar from './sidebar/index.vue'

const appStore = useAppStore()
</script>
