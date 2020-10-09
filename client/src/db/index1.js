class DB {
  constructor(dbName, dbVersion) {
    this.open(dbName, dbVersion)
  }

  open(dbName, dbVersion) {
    this.request = window.indexedDB.open(dbName, dbVersion)
    this.request.onsuccess = (e) => {
      this.db = this.request.result
    }
    this.request.onupgradeneeded = e => {
      let res = e.target.result
      if (!res.objectStoreNames.contains('messages')) {
        const objectStore = res.createObjectStore('messages', { keyPath: 'message_id' })
        objectStore.createIndex("client_id", "client_id", { unique: false })
      } else {
      }
    }
  }

  async addMessage(messageList) {
    if (!Array.isArray(messageList)) messageList = [messageList]
    let store = await this.getTx()
    for (let message of messageList) {
      store.add(message)
    }
  }

  async getAllMessage() {
    let store = await this.getTx()
    let req = store.getAll()
    return resolveData(req, (a, b) => {
      let aTime = new Date(a.created_at)
      let bTime = new Date(b.created_at)
      return aTime - bTime
    })
  }

  async getLastMessage() {
    let store = await this.getTx()
    let req = store.getAll()
    let data = await resolveData(req, (a, b) => {
      let aTime = new Date(a.created_at)
      let bTime = new Date(b.created_at)
      return aTime - bTime
    })
    if (data.length > 0) {
      return { created_at: data[data.length - 1].created_at }
    }
    return { created_at: 0 }
  }


  async getIndexMessage(client_id) {
    let store = await this.getTx()
    const index = store.index("client_id")
    const req = index.getAll(client_id)
    return resolveData(req, (a, b) => {
      let aTime = new Date(a.created_at)
      let bTime = new Date(b.created_at)
      return aTime - bTime
    })

  }

  async getMessage(message_id) {
    let store = await this.getTx()
    let req = store.get(message_id)
    return resolveData(req)
  }

  async updateMessage(message_id, status) {
    let message = await this.getMessage(message_id)
    if (!message) return
    message.status = status
    let store = await this.getTx()
    let req = store.put(message)
    return resolveData(req)
  }

  async getTx() {
    while (!this.db) {
      await wait()
    }
    return this.db.transaction(['messages'], 'readwrite').objectStore('messages')
  }
}

function wait() {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve()
    }, 200)
  })
}

function resolveData(req, fn) {
  let state = false
  return new Promise((resolve, reject) => {
    req.onsuccess = () => {
      state = true
      let { result } = req
      if (!result) return resolve([])
      if (fn) result.sort(fn)
      resolve(result)
    }
    req.onerror = e => {
      state = true
      reject(e)
    }
    setTimeout(() => {
      if (state) return
      resolve([])
    }, 5000)
  })
}

export default new DB('donate', 1)