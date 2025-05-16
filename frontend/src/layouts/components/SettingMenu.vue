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
import { BrowserOpenURL } from 'wailsjs/runtime/runtime.js'

const themeLayoutRef = ref(null)
const aboutDialogRef = ref(null)

const router = useRouter()

function handleMenuSelect(key, item) {
  if (key === 'theme') {
    themeLayoutRef.value.open()
    return
  }
  if (key === 'about') {
    aboutDialogRef.value.show()
    return
  }
  if (key === 'apiDoc') {
    BrowserOpenURL(item.path)
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
    label: '接口文档',
    key: 'apiDoc',
    path: 'https://apifox.com/apidoc/shared-ff4a4d32-c0d1-4caf-b0ee-6abc130f734a',
  },
  {
    label: '关于系统',
    key: 'about',
  },
]
</script>
