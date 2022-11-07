/*
  * Author : see AUTHORS
  * Licence: MIT, see LICENSE
*/

const lskey = 'fitcrack-pagination-store'

export default {
  install (Vue) {
    if (!localStorage.getItem(lskey)) {
      localStorage.setItem(lskey, '{}')
    }

    Vue.prototype.$paginator = function (id, set) {
      const pages = JSON.parse(localStorage.getItem(lskey))
      if (set === undefined) {
        return pages[id] || 1
      } else if (set === null) {
        delete pages[id]
      } else {
        pages[id] = set
      }
      localStorage.setItem(lskey, JSON.stringify(pages))
      return pages[id]
    }
  }
}
