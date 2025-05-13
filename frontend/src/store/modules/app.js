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
import { defineStore } from 'pinia'
import { GetPreferences } from 'wailsjs/go/ipc/PreferenceApi'

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
