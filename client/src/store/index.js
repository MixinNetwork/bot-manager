import Vue from 'vue'
import Vuex from 'vuex'
import Package from './package'

import system from './system'
import user from './user'
import users from './users'
import message from './message'
import data from './data'

Vue.use(Vuex)
Package(Vue)

export default new Vuex.Store({
  modules: {
    user,
    system,
    users,
    message,
    data
  },
  mutations: {
    changeState(state, [space, obj]) {
      if (typeof obj !== 'object') return
      for (let key in obj) {
        state[space][key] = obj[key]
      }
    },
  },
  actions: {
    initPage(ctx) {
      let { get } = _vm.$ls
      let user_info = get('user')
      if (user_info) ctx.commit('changeState', ['user', { user_info }])
      ctx.dispatch('message/init')
    }
  }
})
