import api from '../api'
import { getAvatarColorById } from "../assets/js/color"
import tools from '@/assets/js/tools'

export default {
  namespaced: true,
  state: () => ({
    user_info: {},

    show_add_bot: false,
    bot_list: [],
    active_bot: {},

    addFavoriteModal: false,
    favorite: []
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
      let { get } = _vm.$ls
      const { commit } = ctx
      const bot_list = await api.getBot()
      commit('changeState', { bot_list })
      if (bot_list && bot_list.length > 0) {
        const cacheClientId = get('bot')
        let activeIdx = bot_list.findIndex(item => item.client_id === cacheClientId)
        if (activeIdx === -1) activeIdx = 0
        ctx.dispatch('changeActiveBot', bot_list[activeIdx])
      }
    },
    async changeActiveBot(ctx, item) {
      let { set } = _vm.$ls
      const { state, commit } = ctx
      if (state.active_bot.client_id === item.client_id) return
      commit('changeState', { active_bot: item })
      ctx.dispatch('message/toggleBotMessage', item, { root: true })
      ctx.dispatch('data/toggleStatistics', item, { root: true })
      set('bot', item.client_id)
      if (item.is_reload) window.location.href = '/'
    },
    async getBotFavorite(ctx) {
      let currentClientID = ctx.state.active_bot.client_id
      while (!currentClientID) {
        await tools.delay(0.1)
        currentClientID = ctx.state.active_bot.client_id
      }
      const favorite = await api.getBotFavorite(currentClientID)
      ctx.commit('changeState', { favorite })
    },
    addBotFavorite(ctx, id) {
      let currentClientID = ctx.state.active_bot.client_id
      return api.addBotFavorite(currentClientID, id)
    },
    delBotFavorite(ctx, id) {
      let currentClientID = ctx.state.active_bot.client_id
      return api.deleteBotFavorite(currentClientID, id)
    }
  }
}
