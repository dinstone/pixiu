<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/16 18:50:18
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <div
    id="menu-collapse"
    class="custom-bg-hover f-c-c cursor-pointer rounded-4 p-6 text-20 transition-all-300"
    @click="appStore.switchCollapsed"
  >
    <i :class="appStore.collapsed ? 'i-line-md-menu-unfold-right' : 'i-line-md-menu-fold-left'" />
  </div>
</template>

<script setup>
import { useAppStore } from '@/store'

const appStore = useAppStore()
</script>
