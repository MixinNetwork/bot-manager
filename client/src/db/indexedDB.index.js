import DB from './indexedDB.config'

class Donate extends DB {
  constructor(props) {
    super(props);
    if (!this.db.objectStoreNames.contains('messages')) {
      this.MessageStore = this.db.createObjectStore('messages', { keyPath: 'message_id' })
      this.MessageStore.createIndex("created_at", "created_at", { unique: true });
    } else {
      this.MessageStore = this.db.transaction('messages', 'readwrite')
    }
  }

  async addMessage(messageList) {
    let sql = `INSERT INTO messages (user_id, identity_number, full_name, message_id, avatar_url, category, data, status, source, created_at) VALUES(?,?,?,?,?,?,?,?,?,?)`
    if (!Array.isArray(messageList)) messageList = [messageList]
    // messageList = messageList.map(({ user_id, identity_number, full_name, message_id, avatar_url, category, data, status, source, created_at }) => ([user_id, identity_number, full_name, message_id, avatar_url, category, data, status, source, created_at]))
    messageList.forEach(item => {
      this.MessageStore.add(item)
    })
    // await this.queryBatch(sql, messageList)
  }

  async getLastMessage() {
    let sql = `SELECT created_at FROM messages order by created_at desc limit 1`
    let { created_at } = (await this.query(sql))[0] || {}
    return created_at
  }

  getAllMessage() {
    let sql = `SELECT * FROM messages`
    return this.query(sql)
  }

  async getMessage(message_id) {
    return (await this.query(`SELECT * FROM messages WHERE message_id=?`, [message_id]))[0]
  }

  async updateMessage(message_id, status) {
    await this.query(`UPDATE messages SET status=? WHERE message_id=?`, [status, message_id])
  }


}


export default new Donate('Donate', 'Donate Dashboard')