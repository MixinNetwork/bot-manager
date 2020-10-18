const pgsql = require('pg')

class DB {
  constructor() {
    const pool = new pgsql.Pool({
      host: 'localhost',
      port: 5432,
      user: 'neo',
      password: 'abcd1234',
      database: 'test'
    })
    this.query = async (sql, params) => {
      let client = await pool.connect()
      try {
        let { rows } = await client.query(sql, params)
        return rows
      } finally {
        await client.release()
      }
    }
  }

  get_bot(client_id) {
    return this.query(`SELECT * FROM bots WHERE client_id=$1`, [client_id])
  }

  add_user_bot({ user_id, client_id, session_id, private_key, hash }) {
    return this.query(`INSERT INTO user_bots(user_id, client_id, session_id, private_key, hash) VALUES($1,$2,$3,$4,$5)`, [user_id, client_id, session_id, private_key, hash])
  }
}

const db = new DB();

setTimeout(async () => {
  const bot = (await db.get_bot('d76b821d-ebd7-45cb-a33f-33db978370d2'))[0]
  if (bot.client_id) {
    await db.add_user_bot({ ...bot, user_id: 'b847a455-aa41-4f7d-8038-0aefbe40dcaa' })
    await db.add_user_bot({ ...bot, user_id: 'fcb87491-4fa0-4c2f-b387-262b63cbc112' })
    console.log('ok')
  }
}, 200);