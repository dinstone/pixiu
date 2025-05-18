<template>
  <MeModal ref="modalRef" title="主题设置" :show-footer="false" width="600px">
    <n-space vertical>
      <n-card title="模式">
        <n-space justify="space-between">
          <n-color-picker
            id="theme-color"
            class="h-24 w-24"
            :value="appStore.primaryColor"
            :swatches="primaryColors"
            :on-update:value="(v) => setThemeColor(v)"
            :render-label="() => ''"
          />

          <i
            id="toggleTheme"
            class="h-24 w-24 cursor-pointer"
            :class="appStore.isDark ? 'i-fe:moon' : 'i-fe:sun'"
            @click="toggleThemeDark"
          />

          <n-button @click="setDefaultTheme">
            重置默认
          </n-button>
        </n-space>
      </n-card>
      <n-card title="布局">
        <n-space justify="space-between">
          <div class="flex-col cursor-pointer justify-center" @click="setThemeLayout('simple')">
            <div class="flex">
              <n-skeleton :width="20" :height="60" />
              <div class="ml-4">
                <n-skeleton :width="80" :height="60" />
              </div>
            </div>
            <n-button
              class="mt-12"
              size="small"
              :type="appStore.layout === 'simple' ? 'primary' : ''"
              ghost
            >
              简约
            </n-button>
          </div>
          <div class="flex-col cursor-pointer justify-center" @click="setThemeLayout('classic')">
            <div class="flex flex-col">
              <div class="mb-4">
                <n-skeleton :width="100" :height="10" />
              </div>
              <div class="flex">
                <n-skeleton :width="16" :height="46" />
                <div class="ml-4">
                  <n-skeleton :width="80" :height="46" />
                </div>
              </div>
            </div>
            <n-button
              class="mt-12"
              size="small"
              :type="appStore.layout === 'classic' ? 'primary' : ''"
              ghost
            >
              经典
            </n-button>
          </div>
          <div class="flex-col cursor-pointer justify-center" @click="setThemeLayout('normal')">
            <div class="flex">
              <n-skeleton :width="20" :height="60" />
              <div class="ml-4">
                <n-skeleton :width="80" :height="10" />
                <n-skeleton class="mt-4" :width="80" :height="46" />
              </div>
            </div>
            <n-button
              class="mt-12"
              size="small"
              :type="appStore.layout === 'normal' ? 'primary' : ''"
              ghost
            >
              通用
            </n-button>
          </div>

          <div class="flex-col cursor-pointer justify-center" @click="setThemeLayout('full')">
            <div class="flex">
              <n-skeleton :width="20" :height="60" />
              <div class="ml-4">
                <n-skeleton :width="80" :height="6" />
                <n-skeleton class="mt-4" :width="80" :height="4" />
                <n-skeleton class="mt-4" :width="80" :height="42" />
              </div>
            </div>
            <n-button
              class="mt-12"
              size="small"
              :type="appStore.layout === 'full' ? 'primary' : ''"
              ghost
            >
              全面
            </n-button>
          </div>
        </n-space>
        <p class="mt-16 opacity-50">
          注: 此设置仅对未设置layout或者设置成跟随系统的页面有效，菜单设置的layout优先级最高
        </p>
      </n-card>
    </n-space>
  </MeModal>
</template>

<script setup>
import { MeModal } from '@/components'
import { useModal } from '@/composables'
import { useAppStore } from '@/store'
import { getPresetColors } from '@arco-design/color'
import { useToggle } from '@vueuse/core'
import { UpdatePreferences } from 'wailsjs/go/ipc/PreferenceApi'

const primaryColors = Object.entries(getPresetColors()).map(([, value]) => value.primary)
const appStore = useAppStore()

function setDefaultTheme() {
  const theme = appStore.getDefaultTheme()
  appStore.isDark = theme.dark
  appStore.setPrimaryColor(theme.color)

  UpdatePreferences({
    'theme.dark': appStore.isDark,
    'theme.color': appStore.primaryColor,
  })
}

function setThemeColor(v) {
  appStore.setPrimaryColor(v)
  UpdatePreferences({
    'theme.color': v,
  })
}

function setThemeLayout(v) {
  appStore.setLayout(v)
  UpdatePreferences({
    'theme.layout': v,
  })
}

function toggleThemeDark() {
  appStore.toggleDark()
  useToggle(appStore.isDark)
  UpdatePreferences({
    'theme.dark': appStore.isDark,
  })
}

const [modalRef] = useModal()
function open() {
  modalRef.value?.open()
}

defineExpose({
  open,
})
</script>
