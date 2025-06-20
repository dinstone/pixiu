import { AuthenAccessToken, AuthenPassword } from 'wailsjs/go/ipc/UaacApi.js'

export default {

  login: (data) => {
    if (data.isQuick) {
      return AuthenAccessToken(data.username, data.password)
    }
    else {
      return AuthenPassword(data.username, data.password)
    }
  },

}
