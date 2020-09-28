import api from '../api'
import { getAvatarColorById } from "../assets/js/color"

export default {
  namespaced: true,
  state: () => ({
    token: '',
    user_info: {},
    bot_list: [],
    active_bot: {},
    show_add_bot: false
  }),
  mutations: {
    changeState(state, obj) {
      for (let key in obj) {
        state[key] = obj[key]
      }
    }
  },
  actions: {
    async authenticate(ctx, code) {
      let { set, get } = _vm.$ls
      let user = get('user')
      if (!user) _vm.$DC({ isLoading: true })
      user = await api.authenticate(code)
      if (!user.access_token) return false
      let { access_token, ...user_info } = user
      ctx.commit('changeState', { user_info })
      set('token', access_token)
      set('user', user_info)
      _vm.$DC({ isLoading: false })
      return true
    },
    async getBotList(ctx) {
      const { commit } = ctx
      const bot_list = await api.getBot()
      commit('changeState', { bot_list })
      if (bot_list && bot_list.length > 0) {
        ctx.dispatch('changeActiveBot', bot_list[0])
      }
    },
    async changeActiveBot(ctx, item) {
      const { state, commit } = ctx
      if (state.active_bot.client_id === item.client_id) return
      commit('changeState', { active_bot: item })
      ctx.dispatch('message/toggleBotMessage', item, { root: true })
    }
  }
}