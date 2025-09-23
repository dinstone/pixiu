<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/05 21:24:09
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <SideMenu class="cus-scroll-y mt-4 h-0 flex-1" />

  <div class="mb-10">
    <UserAvatar />
  </div>
  <SettingMenu class="mb-10" />
</template>

<script setup>
import { SettingMenu, SideMenu, UserAvatar } from '@/layouts/components'
</script>
