<template>
  <n-config-provider
    class="wh-full"
    :locale="zhCN"
    :date-locale="dateZhCN"
    :theme="appStore.isDark ? darkTheme : undefined"
    :theme-overrides="appStore.naiveThemeOverrides"
  >
    <router-view v-if="Layout" v-slot="{ Component, route: curRoute }">
      <component :is="Layout">
        <transition name="fade-slide" mode="out-in" appear>
          <KeepAlive :include="keepAliveNames">
            <component :is="Component" v-if="!tabStore.reloading" :key="curRoute.fullPath" />
          </KeepAlive>
        </transition>
      </component>
    </router-view>
  </n-config-provider>
</template>

<script setup>
import { useAppStore, useTabStore } from '@/store'
import { darkTheme, dateZhCN, zhCN } from 'naive-ui'
import { computed, markRaw, watchEffect } from 'vue'

// 利用map将加载过的layout缓存起来，防止重新加载layout导致页面闪烁
const layoutMap = new Map()
function getLayout(name) {
  if (layoutMap.get(name))
    return layoutMap.get(name)
  const layout = markRaw(defineAsyncComponent(() => import(`@/layouts/${name}/index.vue`)))
  layoutMap.set(name, layout)
  return layout
}

const appStore = useAppStore()
// 根据当前路由计算布局，如果路由没有layout，则使用默认布局
const route = useRoute()
const Layout = computed(() => {
  if (!route.matched?.length)
    return null
  return getLayout(route.meta?.layout || appStore.layout)
})

const tabStore = useTabStore()
const keepAliveNames = computed(() => {
  return tabStore.tabs.filter(item => item.keepAlive).map(item => item.name)
})

// 监听主题色并动态更新
watchEffect(() => {
  appStore.setThemeColor(appStore.primaryColor, appStore.isDark)
})
</script>
