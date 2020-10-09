import request from './config'

class APIS {
  authenticate(code) {
    return request.get('/user/login', { params: { code } })
  }

  getUser(client_id, status) {
    return request.get('/user', { params: { client_id, status } })
  }

  putUser(client_id, user_id, status) {
    return request.put('/user', { client_id, user_id, status })
  }

  getData(client_id) {
    return request.get('/data', { params: { client_id } })
  }


  uploadFile(file) {
    let formData = new FormData()
    formData.append('file', file)
    return request.post('/message/uploadFile', formData, { heaaders: { 'Content-Type': 'multipart/form-data' } })
  }

  sendSpecificUser(data, group) {
    return request.post('/sendSpecificUser', { data, group })
  }

  addBot({ client_id, session_id, private_key }) {
    return request.post('/bot', { client_id, session_id, private_key })
  }

  getBot() {
    return request.get('/bot')
  }

  postMessageReplay(data) {
    return request.post('/message/replay', data)
  }

  getMessageReplay(client_id) {
    return request.get('/message/replay', { params: { client_id } })
  }

  deleteMessageReplay(replay_id, client_id) {
    return request.delete('/message/replay', { params: { replay_id, client_id } })
  }

}

export default new APIS()




