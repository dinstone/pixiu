export const defaultLayout = 'classic'

export const defaultPrimaryColor = '#316C72'

export const naiveThemeOverrides = {
  common: {
    primaryColor: '#316C72FF',
    primaryColorHover: '#316C72E3',
    primaryColorPressed: '#2B4C59FF',
    primaryColorSuppl: '#316C72E3',
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
  code: 'Stock',
  name: '股票持仓',
  type: 'MENU',
  path: '/stock',
  icon: 'i-fe:external-link',
  component: '/src/views/stock/index.vue',
  enable: true,
  show: true,
  order: 2,
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
}]
