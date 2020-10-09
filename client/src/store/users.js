import api from '../api'
import tools from '../assets/js/tools'

export default {
  namespaced: true,
  state: () => ({
    stateList: {
      normal: '正常用户',
      block: '已屏蔽',
    },
    currentState: 'normal',
    userList: []
  }),
  mutations: {
    changeState(state, obj) {
      for (let key in obj) {
        state[key] = obj[key]
      }
    }
  },
  actions: {
    async toggleState(ctx, status) {
      status = status || ctx.state.currentState
      ctx.commit('changeState', { currentState: status })
    },
    async getUserList(ctx, status) {
      let currentClientID = ctx.rootState.user.active_bot.client_id
      while (true) {
        await tools.delay(0.1)
        if (!currentClientID) {
          currentClientID = ctx.rootState.user.active_bot.client_id
        } else break
      }
      const userList = await api.getUser(currentClientID, status)
      ctx.commit('changeState', { userList })
    },
    async updateUserStatus(ctx, params) {
      let currentClientID = ctx.rootState.user.active_bot.client_id
      const { user_id, status } = params
      return await api.putUser(currentClientID, user_id, status)
    }
  }
}
