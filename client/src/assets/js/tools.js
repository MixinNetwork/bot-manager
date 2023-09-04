import { getAvatarColorById } from './color'

class Tools {
  async delay(num = 1) {
    return new Promise(resolve => {
      setTimeout(() => {
        resolve()
      }, num * 1000)
    })
  }

  getUrlParameter(name) {
    name = name.replace(/[\[]/, '\\[').replace(/[\]]/, '\\]')
    var regex = new RegExp('[\\?&]' + name + '=([^&#]*)')
    var results = regex.exec(window.location.search)
    return results === null ? '' : decodeURIComponent(results[1].replace(/\+/g, ' '))
  }

  getTheme() {
    let metas = document.getElementsByTagName('meta')
    for (let i = 0; i < metas.length; i++) {
      if (metas[i].name === 'theme-color') {
        return metas[i].content
      }
    }
  }

  changeTheme(color) {
    let head = document.getElementsByTagName('head')[0]
    let metas = document.getElementsByTagName('meta')
    for (let i = 0; i < metas.length; i++) {
      if (metas[i].name === 'theme-color') {
        head.removeChild(metas[i])
      }
    }
    let meta = document.createElement('meta')
    meta.name = 'theme-color'
    meta.content = color
    head.appendChild(meta)
    reloadTheme()
  }

  getEarnDate(date) {
    date = new Date(date)
    let _year = date.getFullYear(),
      _month = date.getMonth() + 1,
      _date = date.getDate(),
      _hours = date.getHours(),
      _minutes = date.getMinutes()
    _month = addZero(_month)
    _date = addZero(_date)
    _hours = addZero(_hours)
    _minutes = addZero(_minutes)
    return `${_year}-${_month}-${_date} ${_hours}:${_minutes}`
  }

  environment() {
    return environment()
  }

  getAvatarColor(id) {
    return getAvatarColorById(id)
  }

  getConversationId() {
    var ctx
    switch (environment()) {
      case 'iOS':
        ctx = prompt('MixinContext.getContext()')
        return JSON.parse(ctx).conversation_id
      case 'Android':
        ctx = window.MixinContext.getContext()
        return JSON.parse(ctx).conversation_id
      default:
        return undefined
    }
  }

  getUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
      const r = (Math.random() * 16) | 0
      const v = c === 'x' ? r : (r & 0x3) | 0x8
      return v.toString(16)
    })
  }

  getStorageContent(key) {
    let t = window.localStorage.getItem(key)
    return t ? JSON.parse(t) : {}
  }
}

function reloadTheme() {
  switch (environment()) {
    case 'iOS':
      return window.webkit.messageHandlers.reloadTheme && window.webkit.messageHandlers.reloadTheme.postMessage('')
    case 'Android':
      return window.MixinContext.reloadTheme()
  }
}

function environment() {
  if (window.webkit && window.webkit.messageHandlers && window.webkit.messageHandlers.MixinContext) {
    return 'iOS'
  }
  if (window.MixinContext && window.MixinContext.getContext) {
    return 'Android'
  }
  return undefined
}

function addZero(num) {
  return num >= 10 ? num : '0' + num
}

export default new Tools()
