import tools from '../assets/js/tools'

class WSSClient {
  constructor() {
    let url = process.env.VUE_APP_WS_SERVER + '?token='
    this.url = url + window.localStorage.getItem('token')

    this.ws = null
    this.isAlive = false
    this.pingInterval = 0
    this.initState = false
    this._initTime = null

    const self = this
    Object.defineProperty(this, 'initTime', {
      get() {
        return self._initTime
      },
      set(created_at) {
        self._initTime = created_at
        created_at = new Date(created_at).toISOString()
        if (self.ws && self.ws.readyState === WebSocket.OPEN) {
          self.initMessage()
        }
      }
    })

    this.start()
  }

  start() {
    this.ws = new WebSocket(this.url)
    this.ws.onmessage = this._on_message.bind(this)
    this.ws.onerror = this._on_error.bind(this)
    this.ws.onclose = this._on_close.bind(this)
    this.ws.onopen = this._on_open.bind(this)

    this.pingInterval = setInterval(() => {
      if (this.ws.readyState === WebSocket.CONNECTING) return
      if (!this.isAlive) return this.ws.close()
      this.isAlive = false
      this.ws.send('ping')
    }, 1000 * 30)
  }

  async initMessage() {
    this.initState = true
    let created_at = new Date(this.initTime).toISOString()
    this.ws.send(JSON.stringify({ created_at }))
  }

  async handleMessage(message) {
  }

  async _on_message({ data }) {
    console.log(data)
    if (data === 'pong') return this.isAlive = true
    data = JSON.parse(data)
    await this.handleMessage(data)
  }

  // async _on_message(data) {
  //   console.log(data)
  //   // if (data === 'pong') return this.isAlive = true
  //   // data = JSON.parse(data)
  //   // await this.handleMessage(data)
  // }

  _on_error(e) {
    console.log(e)
  }

  async _on_close(e) {
    console.log(e)
    clearInterval(this.pingInterval)
    await tools.delay()
    this.start()
  }

  _on_open() {
    console.log('ws connected...' + new Date().toISOString())
    this.isAlive = true
    if (this.initState === false && this.initTime !== null) {
      this.initMessage()
    }

  }
}

export default new WSSClient()