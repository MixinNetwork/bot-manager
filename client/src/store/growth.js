import api from '../api'

export default {
  namespaced: true,
  state: () => ({
    stateList: {
      normal: '正常用户',
      block: '已屏蔽',
    },
    currentState: 'normal',
    growthList: [],
    activeGrower: {},
    auditOpinion: '',
    actionState: {
      show: false
    }
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
      // let growthList = await api.getGrowthUserList(status)
      // ctx.commit('changeState', { growthList, currentState: status })
      ctx.commit('changeState', { currentState: status })
    },
    async updateState(ctx, status) {
      let { activeGrower: { user_id }, auditOpinion } = ctx.state
      return await api.updateGrowthStatus({ user_id, status, auditOpinion })
    }
  }
}