export default {
  namespaced: true,
  state: () => ({
    isLoading: false,
    transitionName: '',
    canTransition: false,
    themeColor: '',
    activeMenus: 0
  }),
  mutations: {
    changeState(state, obj) {
      for (let key in obj) {
        state[key] = obj[key]
      }
    }
  },
  actions: {},
}