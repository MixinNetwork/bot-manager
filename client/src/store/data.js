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


let newUserList = []
let totalUserList = []

let newMessageList = []
let totalMessageList = []

let totalUserCount = 0
let totalMessageCount = 0

function getStatisticsDate({ list, today } = {}) {

  initData()
  handleListData(list, today)
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

function initData() {
  newUserList = []
  totalUserList = []
  newMessageList = []
  totalMessageList = []
  totalUserCount = 0
  totalMessageCount = 0
}

function handleListData(list, today) {
  if (list.length === 0) return handleDaily(today)
  let todayDate = today.date
  let currentDate = list[0].date
  const mapList = transferListToMap(list)
  while (true) {
    if (currentDate === todayDate) break
    handleDaily(mapList[currentDate] || { date: currentDate })
    currentDate = getNextDate(currentDate)
  }
  handleDaily(today)
}

function getNextDate(date) {
  date = new Date(Number(new Date(date)) + 86400000)
  return date.toISOString().slice(0, 10)
}

function transferListToMap(list) {
  let obj = {}
  list.forEach(item => obj[item.date] = item)
  return obj
}

function handleDaily(dataInfo) {
  const { date, users = 0, messages = 0 } = dataInfo
  totalUserList.push([date, totalUserCount += users])
  newUserList.push([date, users])
  totalMessageList.push([date, totalMessageCount += messages])
  newMessageList.push([date, messages])
}
