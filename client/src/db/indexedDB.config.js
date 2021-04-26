class DB {
  constructor(dbName, dbDesc) {
    this.open(dbName, dbDesc)
  }

  open(dbName, dbDesc) {
    this.DB = window.indexedDB.open(dbName)
    this.DB.onerror = (e) => {
      console.log(e)
      console.log('数据打开报错')
    }
    this.DB.onsuccess = (e) => {
      this.db = this.DB.result
      console.log(e)
      console.log('数据库打开成功')
    }
  }

  async queryBatch(sql, params) {
    let tx = await this.getTx()

    return new Promise(resolve => {
      let results = []
      for (let i in params) {
        tx.executeSql(sql, params[i], (tx, result) => {
          results.push(result.rows)
        })
      }
      resolve(results)
    })
  }


  async query(sql, params) {
    let tx = await this.getTx()
    return new Promise(resolve => {
      tx.executeSql(sql, params, (tx, results) => {
        let { rows } = results
        resolve(rows)
      })
    })
  }

  getTx() {
    return new Promise(resolve => {
      this.db.transaction(tx => {
        resolve(tx)
      })
    })
  }
}

export default DB