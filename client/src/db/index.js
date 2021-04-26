import DB from './config'

class Donate extends DB {
  constructor(props) {
    super(props);
    this.query(`
CREATE TABLE IF NOT EXISTS messages (
  user_id             VARCHAR(36) NOT NULL,
  identity_number     VARCHAR(36) NOT NULL,
  full_name           VARCHAR(512) NOT NULL,
  avatar_url          VARCHAR(1024) NOT NULL,
  message_id          VARCHAR(36) NOT NULL PRIMARY KEY,
  category            VARCHAR(36),
  data                VARCHAR(36),
  status              VARCHAR(36) NOT NULL,
  source              VARCHAR(36) NOT NULL,
  created_at          TIMESTAMP WITH TIME ZONE NOT NULL
);`)
    window.db_drop = async () => await this.query(`drop table messages`)
    window.query = this
  }

  async addMessage(messageList) {
    let sql = `INSERT INTO messages (user_id, identity_number, full_name, message_id, avatar_url, category, data, status, source, created_at) VALUES(?,?,?,?,?,?,?,?,?,?)`
    if (!Array.isArray(messageList)) messageList = [messageList]
    messageList = messageList.map(({ user_id, identity_number, full_name, message_id, avatar_url, category, data, status, source, created_at }) => ([user_id, identity_number, full_name, message_id, avatar_url, category, data, status, source, created_at]))
    await this.queryBatch(sql, messageList)
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

  async getMessageByKey(key) {
    return await this.query(`SELECT * FROM messages WHERE data LIKE '%${key}%'`)
  }


}


export default new Donate('Donate', 'Donate Dashboard')