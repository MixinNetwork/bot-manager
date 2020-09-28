import api from '../api'
import {ASSET_COLOR, ASSETS} from '../assets/js/const'

export default {
  namespaced: true,
  state: () => ({
    statistics: [],

  }),
  mutations: {
    changeState(state, obj) {
      for (let key in obj) {
        state[key] = obj[key]
      }
    }
  },
  actions: {
    async initState(ctx) {
      const resp = await api.getData()
      let tmp = getStatisticsDate(resp)

      ctx.commit('changeState', tmp)
    },
    toggleAsset({state}, type) {
      type += 'AssetObj'
      state.statistics.splice(2, 1, state[type])
    }
  }
}


const newUserList = []
const totalUserList = []
const activeUserList = []
const uvList = []
const siteUvList = []

const assetList = []
let assetTotalAmount = 0
let assetTodayAmount = 0
const donateList = []
let donateTotalAmount = 0
let donateTodayAmount = 0

const depositList = []
let depositTotalAmount = 0
let depositTodayAmount = 0

const peopleList = []
let total_uv = 0
let total_times = 0

function getStatisticsDate({list, price_list, today} = {}) {
  if (!list) return
  for (const dataInfo of list) {
    handleDaily(dataInfo, price_list)
  }
  handleDaily(today, price_list)
  let userObj = {
    name: '用户',
    total: today.users,
    today: today.new_user,
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
  let peopleObj = {
    name: '留言次数',
    total: total_times,
    today: today.times,
    data: [
      {
        name: '当日留言',
        color: '#1BACC0',
        list: peopleList
      },
      {
        name: '总留言',
        color: '#1BACC0',
        list: peopleList
      }
    ]
  }

  return {statistics: [userObj, peopleObj]}
}


function handleDaily(dataInfo, price_list) {
  const {date, users = 0, new_user = 0, active_user = 0, uv = 0, site_uv = 0, times = 0, donate, deposit} = dataInfo
  newUserList.push([date, new_user])
  totalUserList.push([date, users])
  activeUserList.push([date, active_user])
  uvList.push([date, uv])
  siteUvList.push([date, site_uv])
  peopleList.push([date, times])
  total_uv += uv + site_uv
  total_times += times || 0


  Object.keys(ASSETS).forEach((symbol, index) => {
    let lowerSymbol = symbol.toLowerCase()
    let asset_id = ASSETS[symbol]

    generateList(donateList, donate, 'donate')
    generateList(depositList, deposit, 'deposit')


    function generateList(list, obj, type) {
      obj[lowerSymbol] = obj[lowerSymbol] || 0
      let idx = list.findIndex(({name}) => name === symbol)
      let addAmount = Number(obj[lowerSymbol]) * Number(price_list[asset_id])
      assetTotalAmount += addAmount
      type === 'donate' ? (donateTotalAmount += addAmount) : (depositTotalAmount += addAmount)
      if (date === new Date().toISOString().slice(0, 10)) {
        type === 'donate' ? (donateTodayAmount += addAmount) : (depositTodayAmount += addAmount)
        assetTodayAmount += addAmount
      }
      let totalAmount = Number(donate[lowerSymbol]) + Number(deposit[lowerSymbol] || 0)
      if (idx === -1) {
        list.push({name: symbol, color: ASSET_COLOR[index], list: [[date, obj[lowerSymbol]]]})
        type === 'donate' && assetList.push({
          name: symbol,
          color: ASSET_COLOR[index],
          list: [[date, totalAmount]]
        })
      } else {
        list[idx].list.push([date, obj[lowerSymbol]])
        type === 'donate' && assetList[idx].list.push([date, totalAmount])
      }
    }
  })
}
