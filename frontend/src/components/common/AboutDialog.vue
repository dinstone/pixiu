<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/16 18:49:53
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <MeModal ref="modalRef" title="" :show-footer="false" width="500px">
    <n-space :size="10" :wrap="false" :wrap-item="false" align="center" vertical>
      <div class="h-64 w-64 rounded-4 bg-primary">
        <img src="@/assets/images/logo.png" alt="Logo">
      </div>
      <div class="about-app-title">
        {{ title }}
      </div>
      <n-text>{{ version }}</n-text>
      <n-space :size="5" :wrap="false" :wrap-item="false" align="center">
        <n-text class="cursor-pointer" @click="onOpenSource">
          源码工程
        </n-text>
        <n-divider vertical />
        <n-text class="cursor-pointer" @click="onOpenWebsite">
          官方网站
        </n-text>
      </n-space>
      <div class="about-copyright">
        Copyright © 2025 dinstone All rights reserved
      </div>
    </n-space>
  </MeModal>
</template>

<script setup>
import { MeModal } from '@/components'
import { useModal } from '@/composables'
import { BrowserOpenURL } from 'wailsjs/runtime/runtime.js'

const title = import.meta.env.VITE_TITLE

const [modalRef] = useModal()
const version = ref('')

onMounted(() => {
  // GetAppVersion().then(({ data }) => {
  //     version.value = data.version
  // })
  version.value = '1.0.0'
})

function onOpenSource() {
  BrowserOpenURL('https://github.com/dinstone/pixiu')
}

function onOpenWebsite() {
  BrowserOpenURL('https://github.com/zclzone/vue-naive-admin')
}

function show() {
  modalRef.value?.open()
}

defineExpose({
  show,
})
</script>

<style scoped>
.about-app-title {
  font-weight: bold;
  font-size: 18px;
  margin: 5px;
}

.about-link {
  cursor: pointer;

  &:hover {
    text-decoration: underline;
  }
}

.about-copyright {
  font-size: 12px;
}
</style>
