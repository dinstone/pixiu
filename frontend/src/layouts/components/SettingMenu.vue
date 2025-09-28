<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/16 19:00:00
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <div class="custom-bg-hover f-c-c cursor-pointer rounded-4 p-6 text-20 transition-all-300">
    <n-dropdown
      :options="staticMenus"
      trigger="click"
      @select="handleMenuSelect"
    >
      <i class="i-fe:settings cursor-pointer text-20" />
    </n-dropdown>

    <ThemeLayout ref="themeLayoutRef" />
    <AboutDialog ref="aboutDialogRef" />
  </div>
</template>

<script setup>
import { AboutDialog, ThemeLayout } from '@/components'
import { useAppStore } from '@/store'
import { BrowserOpenURL } from 'wailsjs/runtime/runtime.js'

const themeLayoutRef = ref(null)
const aboutDialogRef = ref(null)

const router = useRouter()
const appStore = useAppStore()

function handleMenuSelect(key, item) {
  if (key === 'theme') {
    themeLayoutRef.value.open()
    return
  }
  if (key === 'about') {
    aboutDialogRef.value.show()
    return
  }
  if (key === 'issue') {
    BrowserOpenURL(item.path)
    return
  }
  if (key === 'update') {
    appStore.checkForUpdate(true)
    return
  }
  if (!item.path)
    return
  router.push(item.path)
}

// 静态菜单数据
const staticMenus = [

  {
    label: '主题设置',
    key: 'theme',
    icon: () => h('i', { class: 'i-fe:layout text-14' }),
  },
  {
    label: '问题反馈',
    key: 'issue',
    icon: () => h('i', { class: 'i-fe:book text-14' }),
    path: 'https://github.com/dinstone/pixiu/issues',
  },
  {
    label: '检查更新',
    key: 'update',
    icon: () => h('i', { class: 'i-fe:download text-14' }),
  },
  {
    type: 'divider',
    key: 'd1',
  },
  {
    label: '关于系统',
    key: 'about',
  },
]
</script>
