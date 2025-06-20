import { useAuthStore } from '@/store'
import { UpdatePassword, UpdateProfile } from 'wailsjs/go/ipc/UaacApi.js'

export default {
  changePassword: data => UpdatePassword (useAuthStore().accessToken, data.newPassword),
  updateProfile: data => UpdateProfile(useAuthStore().accessToken, data),
}
