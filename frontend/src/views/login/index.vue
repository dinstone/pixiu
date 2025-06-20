<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/05 21:28:36
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <div class="wh-full flex-col bg-cover">
    <div
      class="m-auto max-w-700 min-w-345 f-c-c rounded-8 bg-opacity-20 bg-cover p-12 card-shadow auto-bg"
    >
      <div class="hidden w-380 px-20 py-35 md:block">
        <img src="@/assets/images/banner.webp" class="w-full" alt="login_banner">
      </div>

      <div class="w-320 flex-col px-20 py-32">
        <h2 class="f-c-c text-24 text-#6a6a6a font-normal">
          <img src="@/assets/images/logo.png" class="mr-12 h-50 bg-primary">
          {{ title }}
        </h2>
        <n-input
          v-model:value="loginInfo.username"
          autofocus
          class="mt-32 h-40 items-center"
          placeholder="请输入账号"
          :maxlength="20"
        >
          <template #prefix>
            <i class="i-fe:user mr-12 opacity-20" />
          </template>
        </n-input>
        <n-input
          v-model:value="loginInfo.password"
          class="mt-20 h-40 items-center"
          type="password"
          show-password-on="mousedown"
          placeholder="请输入密码"
          :maxlength="20"
          @keydown.enter="handleLogin()"
        >
          <template #prefix>
            <i class="i-fe:lock mr-12 opacity-20" />
          </template>
        </n-input>

        <n-checkbox
          class="mt-20"
          :checked="rememberRef"
          label="记住我，可一键登录"
          :on-update:checked="(val) => (rememberRef = val)"
        />

        <div class="mt-20 flex items-center">
          <n-button
            class="h-40 flex-1 rounded-5 text-16"
            type="primary"
            ghost
            @click="quickLogin()"
          >
            一键登录
          </n-button>

          <n-button
            class="ml-32 h-40 flex-1 rounded-5 text-16"
            type="primary"
            :loading="loading"
            @click="handleLogin()"
          >
            登录
          </n-button>
        </div>
      </div>
    </div>

    <TheFooter class="py-12" />
  </div>
</template>

<script setup>
import { useAuthStore } from '@/store'
import { throttle } from '@/utils'
import { useStorage } from '@vueuse/core'
import api from './api'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const title = import.meta.env.VITE_TITLE

const rememberRef = useStorage('isRemember', true)
const usernameRef = useStorage('username')
const passwordRef = useStorage('password')

const loginInfo = ref({ username: usernameRef.value, password: '' })

const captchaUrl = ref('')
const initCaptcha = throttle(() => {
  captchaUrl.value = `${import.meta.env.VITE_AXIOS_BASE_URL}/auth/captcha?${Date.now()}`
}, 500)
initCaptcha()

function quickLogin() {
  loginInfo.value.username = usernameRef.value
  loginInfo.value.password = passwordRef.value
  handleLogin(true)
}

const loading = ref(false)
async function handleLogin(isQuick) {
  const { username, password, captcha } = loginInfo.value
  if (!username || !password) {
    return $message.warning('请输入账号和密码')
  }

  try {
    loading.value = true
    $message.loading('正在验证，请稍后...', { key: 'login' })
    const res = await api.login({ username, password, captcha, isQuick })
    if (res.code === 0) {
      if (rememberRef.value) {
        usernameRef.value = username
        passwordRef.value = res.data
      }
      else {
        usernameRef.value = ''
        passwordRef.value = ''
      }
      onLoginSuccess(res.data)
    }
    else {
      passwordRef.value = ''
      $message.error(res.mesg, { key: 'login' })
    }
  }
  catch (error) {
    // 10003为验证码错误专属业务码
    if (error?.code === 10003) {
      // 为防止爆破，验证码错误则刷新验证码
      initCaptcha()
    }
    $message.destroy('login')
    console.error(error)
  }
  loading.value = false
}

async function onLoginSuccess(token) {
  authStore.setToken(token)
  $message.loading('登录中...', { key: 'login' })
  try {
    $message.success('登录成功', { key: 'login' })
    if (route.query.redirect) {
      const path = route.query.redirect
      delete route.query.redirect
      router.push({ path, query: route.query })
    }
    else {
      router.push('/')
    }
  }
  catch (error) {
    console.error(error)
    $message.destroy('login')
  }
}
</script>
