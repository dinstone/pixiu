/**********************************
 * @FilePath: naiveTools.js
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/04 22:45:20
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { useAppStore } from '@/store'
import { isNullOrUndef } from '@/utils'
import * as NaiveUI from 'naive-ui'

function setupMessage(NMessage) {
  class Message {
    static instance
    constructor() {
      // 单例模式
      if (Message.instance)
        return Message.instance

      Message.instance = this
      this.messageMap = {}
      this.removeTimer = {}
    }

    removeMessage(key, duration = 5000) {
      this.removeTimer[key] && clearTimeout(this.removeTimer[key])
      this.removeTimer[key] = setTimeout(() => {
        this.messageMap[key]?.destroy()
      }, duration)
    }

    destroy(key, duration = 200) {
      setTimeout(() => {
        this.messageMap[key]?.destroy()
      }, duration)
    }

    showMessage(type, content, option = {}) {
      // if (Array.isArray(content)) {
      //   return content.forEach(msg => NMessage[type](msg, option))
      // }

      if (!option.key) {
        return NMessage[type](content, option)
      }

      const currentMessage = this.messageMap[option.key]
      if (currentMessage) {
        currentMessage.type = type
        currentMessage.content = content
      }
      else {
        this.messageMap[option.key] = NMessage[type](content, {
          ...option,
          duration: 0,
          onAfterLeave: () => {
            delete this.messageMap[option.key]
          },
        })
      }

      this.removeMessage(option.key, option.duration)
      return this.messageMap[option.key]
    }

    loading(content, option) {
      return this.showMessage('loading', content, option)
    }

    success(content, option) {
      return this.showMessage('success', content, option)
    }

    error(content, option) {
      return this.showMessage('error', content, option)
    }

    info(content, option) {
      return this.showMessage('info', content, option)
    }

    warning(content, option) {
      return this.showMessage('warning', content, option)
    }
  }

  return new Message()
}

function setupDialog(NDialog) {
  NDialog.confirm = function (option = {}) {
    const showIcon = !isNullOrUndef(option.title)
    return NDialog[option.type || 'warning']({
      showIcon,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: option.confirm,
      onNegativeClick: option.cancel,
      onMaskClick: option.cancel,
      ...option,
    })
  }

  return NDialog
}

export function setupNaiveDiscreteApi() {
  const appStore = useAppStore()
  const configProviderProps = computed(() => ({
    theme: appStore.isDark ? NaiveUI.darkTheme : undefined,
    themeOverrides: useAppStore().naiveThemeOverrides,
  }))
  const { message, dialog, notification, loadingBar } = NaiveUI.createDiscreteApi(
    ['message', 'dialog', 'notification', 'loadingBar'],
    {
      configProviderProps,
      notificationProviderProps: {
        max: 5,
        placement: 'bottom-right',
        keepAliveOnHover: true,
        containerStyle: {
          marginBottom: '32px',
        },
      },
    },
  )

  window.$loadingBar = loadingBar
  window.$notification = notification
  window.$message = setupMessage(message)
  window.$dialog = setupDialog(dialog)
}
