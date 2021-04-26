class DB {
  constructor(dbName, dbDesc) {
    this.open(dbName, dbDesc)
  }

  open(dbName, dbDesc) {
    this.db = openDatabase(dbName, '1.0', dbDesc, 10 * 1024 * 1024)
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