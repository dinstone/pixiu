/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:25:31
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { defaultLayout, defaultPrimaryColor, naiveThemeOverrides } from '@/settings'
import { generate, getRgbStr } from '@arco-design/color'
import { useDark } from '@vueuse/core'
import { get } from 'lodash'
import { NButton, NSpace } from 'naive-ui'
import { defineStore } from 'pinia'
import { GetPreferences } from 'wailsjs/go/ipc/PreferenceApi'
import { BrowserOpenURL } from 'wailsjs/runtime/runtime.js'

export const useAppStore = defineStore('app', {
  state: () => ({
    collapsed: true,
    isDark: useDark(),
    layout: defaultLayout,
    primaryColor: defaultPrimaryColor,
    naiveThemeOverrides,
  }),
  actions: {
    getDefaultTheme() {
      return {
        layout: defaultLayout,
        dark: useDark(),
        color: defaultPrimaryColor,
      }
    },
    async loadPreferences() {
      const { code, data } = await GetPreferences()
      if (code === 0) {
        const layout = get(data, 'theme.layout')
        if (layout && this.layout !== layout) {
          this.layout = layout
        }
        const color = get(data, 'theme.color')
        if (color && this.primaryColor !== color) {
          this.primaryColor = color
        }
        const dark = get(data, 'theme.dark')
        if (this.isDark !== dark) {
          this.isDark = dark
        }
      }
    },
    async checkForUpdate(manual = false) {
      if (manual) {
        $message.loading('正在检索新版本', { key: 'checkUpdate', duration: 10000 })
      }
      try {
        const { success, data = {} } = {
          success: true,
          data: {
            version: 'v1.0.0',
            latest: 'v1.2.0',
            download_page: 'https://github.com/dinstone/pixiu',
            description: '',
          },
        } // await CheckForUpdate()

        if (success) {
          const {
            latest,
            download_page,
            description,
          } = data

          const downUrl = download_page || ''

          if (manual) {
            $message.success('检索到新版本', { key: 'checkUpdate' })
          }

          const notiRef = $notification.info({
            title: `有可用新版本 - ${latest}`,
            content: description || `新版本 ${latest}, 是否立即下载`,
            action: () =>
              h('div', { class: 'flex-box-h flex-item-expand' }, [
                h(NSpace, { wrapItem: false }, () => [
                  h(
                    NButton,
                    {
                      size: 'small',
                      secondary: true,
                      onClick: notiRef.destroy,
                    },
                    () => '稍后下载',
                  ),
                  h(
                    NButton,
                    {
                      type: 'primary',
                      size: 'small',
                      secondary: true,
                      onClick: () => BrowserOpenURL(downUrl),
                    },
                    () => '立即下载',
                  ),
                ]),
              ]),
            onPositiveClick: () => BrowserOpenURL(downUrl),
          })
          return
        }

        if (manual) {
          $message.info('当前已是最新版', { key: 'checkUpdate' })
        }
      }
      finally {
        nextTick().then(() => {
          if (manual) {
            $message.destroy('checkUpdate')
          }
        })
      }
    },
    switchCollapsed() {
      this.collapsed = !this.collapsed
    },
    setCollapsed(b) {
      this.collapsed = b
    },
    toggleDark() {
      this.isDark = !this.isDark
    },
    setLayout(v) {
      this.layout = v
    },
    setPrimaryColor(color) {
      this.primaryColor = color
    },
    setThemeColor(color = this.primaryColor, isDark = this.isDark) {
      const colors = generate(color, {
        list: true,
        dark: isDark,
      })
      document.body.style.setProperty('--primary-color', getRgbStr(colors[5]))
      this.naiveThemeOverrides.common = Object.assign(this.naiveThemeOverrides.common || {}, {
        primaryColor: colors[5],
        primaryColorHover: colors[4],
        primaryColorSuppl: colors[4],
        primaryColorPressed: colors[6],
      })
    },
  },
  persist: {
    pick: ['collapsed', 'layout', 'primaryColor', 'naiveThemeOverrides'],
    storage: sessionStorage,
  },
})
