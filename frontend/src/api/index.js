/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/04 22:50:38
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { useAuthStore } from '@/store'
import { request } from '@/utils'
import { GetUserDetail } from 'wailsjs/go/ipc/UaacApi.js'

export default {
  // 获取用户信息
  getUserDetail: () => GetUserDetail(useAuthStore().accessToken),
  // 登出
  logout: () => '',
  // 切换当前角色
  switchCurrentRole: role => request.post(`/auth/current-role/switch/${role}`),
  // 验证菜单路径
  validateMenuPath: path => ({ path, data: true }),
}
