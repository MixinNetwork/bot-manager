import api from '../api'
import { ASSET_COLOR, ASSETS } from '../assets/js/const'

export default {
  namespaced: true,
  state: () => ({
    statistics: [],
    user_manager: []
  }),
  mutations: {
    changeState(state, obj) {
      for (let key in obj) {
        state[key] = obj[key]
      }
    }
  },
  actions: {
    async toggleStatistics(ctx, bot) {
      const resp = await api.getData(bot.client_id)
      let tmp = getStatisticsDate(resp)
      ctx.commit('changeState', tmp)
    },
    toggleAsset({ state }, type) {
      type += 'AssetObj'
      state.statistics.splice(2, 1, state[type])
    }
  }
}


const newUserList = []
const totalUserList = []

const newMessageList = []
const totalMessageList = []

let totalUserCount = 0
let totalMessageCount = 0

function getStatisticsDate({ list, today } = {}) {
  if (list) {
    for (const dataInfo of list) {
      handleDaily(dataInfo)
    }
  }
  handleDaily(today)
  let userObj = {
    name: '用户',
    total: totalUserCount,
    today: today.users,
    data: [
      {
        name: '新用户',
        color: '#1BACC0',
        list: newUserList
      },
      {
        name: '总用户',
        color: '#5099DD',
        list: totalUserList
      }
    ]
  }
  let messageObj = {
    name: '留言次数',
    total: totalMessageCount,
    today: today.messages,
    data: [
      {
        name: '当日留言',
        color: '#1BACC0',
        list: newMessageList
      },
      {
        name: '总留言',
        color: '#5099DD',
        list: totalMessageList
      }
    ]
  }
  return { statistics: [userObj, messageObj] }
}

function handleDaily(dataInfo) {
  const { date, users = 0, messages = 0 } = dataInfo
  totalUserList.push([date, totalUserCount += users])
  newUserList.push([date, users])

  totalMessageList.push([date, totalMessageCount += messages])
  newMessageList.push([date, messages])
}
