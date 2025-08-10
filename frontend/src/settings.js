export const defaultLayout = 'classic'

export const defaultPrimaryColor = '#D33A31'

export const naiveThemeOverrides = {
  common: {
    primaryColor: '#D33A31FF',
    primaryColorHover: '#FF6B6BFF',
    primaryColorPressed: '#D5271CFF',
    primaryColorSuppl: '#FF6B6BFF',
  },
}

export const basePermissions = [{
  code: 'Home',
  name: '首页',
  type: 'MENU',
  path: '/',
  icon: 'i-fe:home',
  component: '/src/views/home/index.vue',
  show: true,
  enable: true,
  order: 1,
}, {
  code: 'holdingList',
  name: '股票持仓',
  type: 'MENU',
  path: '/stock/holding',
  icon: 'i-fe:external-link',
  component: '/src/views/stock/hold.vue',
  enable: true,
  show: true,
  order: 2,
}, {
  code: 'clearList',
  name: '股票清仓',
  type: 'MENU',
  path: '/stock/clear/list',
  icon: 'i-fe:check-square',
  component: '/src/views/stock/clear-list.vue',
  enable: true,
  show: true,
  order: 3,
}, {
  code: 'clearHistory',
  name: '清仓历史',
  type: 'MENU',
  path: '/stock/clear/history',
  icon: 'i-fe:check-square',
  component: '/src/views/stock/clear-history.vue',
  enable: true,
  show: false,
  order: 4,
}, {
  code: 'UserProfile',
  name: '个人资料',
  type: 'MENU',
  parentId: null,
  path: '/profile',
  redirect: null,
  icon: 'i-fe:user',
  component: '/src/views/profile/index.vue',
  layout: null,
  keepAlive: null,
  method: null,
  description: null,
  show: false,
  enable: true,
  order: 99,
}, {
  code: 'GoToSite',
  name: '外链',
  type: 'MENU',
  path: 'https://GoToSite',
  order: 98,
  enable: true,
  show: false,
}]
