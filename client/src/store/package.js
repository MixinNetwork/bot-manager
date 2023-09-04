
const T = {
  slider: ['system', { canTransition: true }]
}


export default function (Vue) {
  Vue.prototype.$DC = function (namespace, obj) {
    typeof namespace === 'object' && (obj = namespace, namespace = 'system')
    T[namespace] && (obj = T[namespace][1], namespace = T[namespace][0])
    this.$store.commit("changeState", [namespace, obj]);
  }
  Vue.prototype.$ls = {
    get(key) {
      let value = window.localStorage.getItem(key)
      try {
        return JSON.parse(value)
      } catch (e) {
        return value
      }
    },
    set(key, value) {
      if (typeof value === 'object') value = JSON.stringify(value)
      return window.localStorage.setItem(key, value)
    }
  }
}