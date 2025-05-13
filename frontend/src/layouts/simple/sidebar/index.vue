<template>
  <SideLogo border-b="1px solid light_border dark:dark_border" />
  <SideMenu class="cus-scroll-y mt-4 h-0 flex-1" />

  <div class="mb-10">
    <UserAvatar />
  </div>
  <SettingMenu />
  <MenuCollapse />
</template>

<script setup>
import { MenuCollapse, SettingMenu, SideLogo, SideMenu, UserAvatar } from '@/layouts/components'
import { useAppStore } from '@/store'

useAppStore()
</script>
