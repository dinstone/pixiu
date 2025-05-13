<!--------------------------------
 - @Author: dinstone
 - @LastEditor: dinstone
 - @LastEditTime: 2025-4-21 16:04:04
 - @Email: dinstone@163.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <router-link class="h-60 f-c-c" to="/">
    <TheLogo />
    <h2 class="ml-10 max-w-140 flex-shrink-0 text-16 color-primary font-bold">
      {{ title }}
    </h2>
  </router-link>
</template>

<script setup>
const title = import.meta.env.VITE_TITLE
</script>
