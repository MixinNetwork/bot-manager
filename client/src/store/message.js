import db from '../db/index1'
import api from '../api'
import wss from '../api/ws'


let currentClientID = ""

export default {
  namespaced: true,
  state: () => ({
    currentState: 'message',
    stateList: {
      message: "用户消息",
      broadcast: "公告",
      key: "关键字回复",
      add: "被加好友回复"
    },
    messageTypeMap: { PLAIN_TEXT: '文字消息', PLAIN_IMAGE: '图片', BUTTON_GROUP: "按钮" },
    activeType: 'PLAIN_TEXT',
    activeContent: '',

    // 消息相关
    replyModal: false,
    originList: [],
    contactList: [],
    activeContact: {},
    searchKey: '',

    // 公告相关
    broadcastModal: false,
    broadcastModalType: 'send',
    broadcastList: [
      {
        content: "【 Mixin 公告 】Mixin 不支持 BTC 分叉巴拉巴拉巴拉巴拉巴拉巴拉巴拉巴拉巴拉巴拉",
        full_name: "Kelly 分叉巴拉巴拉巴拉巴拉巴拉巴拉巴拉巴拉巴拉巴拉",
        created_at: "2020-04-22 13:04:12",
        category: "PLAIN_TEXT",
        status: "finished"
      },
      {
        content: "【 Mixin 公告 】Mixin 不支持 BTC 分叉巴拉巴拉巴拉巴拉巴拉巴拉巴拉巴拉巴拉巴拉",
        full_name: "Kelly",
        created_at: "2020-04-22 13:04:12",
        category: "PLAIN_IMAGE",
        status: "finished"
      },
      {
        content: "【 Mixin 公告 】Mixin 不支持 BTC 分叉巴拉巴拉巴拉巴拉巴拉巴拉巴拉巴拉巴拉巴拉",
        full_name: "Kelly",
        created_at: "2020-04-22 13:04:12",
        category: "APP_BUTTON",
        status: "finished"
      },
    ],
    activeBroadcast: {},

    // 关键字回复相关
    keyModal: false,
    keyModalType: 'edit',
    keyList: [],
    activeKey: "",
    activeReplayId: "",

    helloData: "",
    helloCategory: "",
    //


  }),
  mutations: {
    changeState(state, obj) {
      for (let key in obj) {
        state[key] = obj[key]
      }
    }
  },
  actions: {
    async init(ctx) {
      const lastMessage = await db.getLastMessage()
      wss.initTime = lastMessage.created_at
      wss.handleMessage = async data => {
        if (Array.isArray(data)) {
          let addMessage = data.filter(item => item.source !== 'ACKNOWLEDGE_MESSAGE_RECEIPT')
          handleMessage(addMessage, ctx)
          _vm.$bus.$emit('message', addMessage[addMessage.length - 1])
          for (const item of data) {
            if (item.source === 'ACKNOWLEDGE_MESSAGE_RECEIPT') {
              await findMessageAndUpdate(item, ctx.state.originList)
            } else {
              await db.addMessage(item)
            }
          }
        } else {
          if (data.source === 'ACKNOWLEDGE_MESSAGE_RECEIPT') {
            await findMessageAndUpdate(data, ctx.state.originList)
          } else {
            await db.addMessage(data)
            handleMessage([data], ctx)
            _vm.$bus.$emit('message', data)
          }
        }
      }
    },
    async toggleBotMessage(ctx, item) {
      currentClientID = ctx.rootState.user.active_bot.client_id
      const originMessage = await db.getIndexMessage(item.client_id)
      handleMessage(originMessage, ctx, true)
    },
    sendMessage(ctx, payload) {
      let user_id = ctx.rootState.user.user_info.user_id,
        recipient_id = ctx.state.activeContact.user_id,
        { data, category } = payload,
        client_id = ctx.rootState.user.active_bot.client_id
      wss.ws.send(JSON.stringify({ client_id, recipient_id, data, user_id, category }))
    },
    sendBroadcast(ctx, payload) {
      let data = payload
      let group = ctx.state.broadcastGroup === 'chinese' ? 'CNY' : ''
      return api.sendSpecificUser(data, group)
    },
    async searchKey({ state, commit }, searchKey) {
      state.searchKey = searchKey
      state.contactList = state.originList.filter(({ full_name, identity_number }) => full_name.includes(searchKey) || identity_number.includes(searchKey))
    },
    async toggleState(ctx, status) {
      status = status || ctx.state.currentState
      ctx.commit('changeState', { currentState: status })
    },


    // messageReplay
    addOrUpdateMessageReplay(ctx) {
      currentClientID = ctx.rootState.user.active_bot.client_id
      const { activeType, activeContent, activeKey, activeReplayId } = ctx.state
      return api.postMessageReplay({
        replay_id: activeReplayId || "",
        keys: activeKey.split(" "),
        category: activeType,
        data: activeContent,
        client_id: currentClientID,
      })
    },
    async getMessageReplayList(ctx) {
      currentClientID = ctx.rootState.user.active_bot.client_id
      const replayList = await api.getMessageReplay(currentClientID)
      const keyList = []
      replayList.forEach(item => {
        if (["Hi", "你好"].includes(item.key)) {
          return ctx.commit('changeState', {
            helloData: item.data,
            helloCategory: item.category,
          })
        }

        const idx = keyList.findIndex(subItem => subItem.replay_id === item.replay_id)
        if (idx === -1) {
          keyList.push({
            replay_id: item.replay_id,
            key: item.key,
            category: item.category,
            content: item.data
          })
        } else {
          keyList[idx].key += ` ${item.key}`
        }
      })
      ctx.commit('changeState', { keyList })
    },
    deleteMessageReplay(ctx, replay_id) {
      currentClientID = ctx.rootState.user.active_bot.client_id
      return api.deleteMessageReplay(replay_id, currentClientID)
    },
    async clickAddReplay(ctx) {
      let { helloData, helloCategory } = ctx.state
      if (!helloCategory) {
        await ctx.dispatch('getMessageReplayList')
        let { helloData, helloCategory } = ctx.state
        ctx.commit("changeState", {
          activeType: helloCategory || "PLAIN_TEXT",
          activeContent: helloData
        })
      } else {
        ctx.commit("changeState", {
          activeType: helloCategory || "PLAIN_TEXT",
          activeContent: helloData
        })
      }
    }
  }
}

function handleMessage(originMessage, ctx, initStatus) {
  const target = initStatus ? [] : ctx.state.contactList
  for (let i = 0, len = originMessage.length; i < len; i++) {
    let messageItem = originMessage[i]
    if (currentClientID !== messageItem.client_id) continue
    let { idx, obj } = findIndexOrObj(messageItem, target)
    if (idx === -1) target.unshift(obj)
    else {
      target[idx].messages.push(obj)
      target.unshift(target.splice(idx, 1)[0])
    }
    if (ctx) {
      ctx.commit('changeState', { contactList: target, originList: target })
    }
  }
}

function findIndexOrObj(messageItem, target = []) {
  let idx = target.findIndex(({ user_id }) => user_id === messageItem.user_id)
  let { message_id, category, avatar_url, full_name, identity_number, data, created_at, user_id, status } = messageItem
  created_at = formatDate(created_at)
  if (idx === -1) {
    let obj = {
      avatar_url, user_id, identity_number, full_name, repeatStatus: false,
      messages: [{ message_id, category, data, created_at, status }]
    }
    return { idx, obj }
  } else {
    full_name && (target[idx].full_name = full_name)
    avatar_url && (target[idx].avatar_url = avatar_url)
    identity_number && (target[idx].identity_number = identity_number)
    let obj = { message_id, category, data, created_at, status }
    return { idx, obj }
  }
}

async function findMessageAndUpdate({ message_id, status }, contactList) {
  await db.updateMessage(message_id, status)
  let { user_id } = await db.getMessage(message_id) || {}
  if (!user_id) return
  let wIdx = contactList.findIndex(item => item.user_id === user_id)
  let { messages } = contactList[wIdx]
  for (let len = messages.length - 1; len >= 0; len--) {
    let { message_id: _message_id } = messages[len]
    if (message_id === _message_id) {
      messages[len].status = status
      break
    }
  }
}

function formatDate(date) {
  date = new Date(date)
  let hours = date.getHours(),
    minutes = date.getMinutes()
  hours < 10 && (hours = '0' + hours)
  minutes < 10 && (minutes = '0' + minutes)
  return `${hours}:${minutes}`
}
