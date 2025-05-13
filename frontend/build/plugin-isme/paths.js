/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:37:43
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { getPagePaths } from '..'

const PLUGIN_PAGE_PATHS_ID = 'isme:page-paths'
export function pluginPagePaths() {
  return {
    name: 'isme:page-paths',
    resolveId(id) {
      if (id === PLUGIN_PAGE_PATHS_ID)
        return `\0${PLUGIN_PAGE_PATHS_ID}`
    },
    load(id) {
      if (id === `\0${PLUGIN_PAGE_PATHS_ID}`) {
        return `export default ${JSON.stringify(getPagePaths())}`
      }
    },
  }
}
