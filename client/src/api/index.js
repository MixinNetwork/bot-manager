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

  getImgUrl(client_id, attachment_id) {
    return request.get('/message/getViewURL', { params: { client_id, attachment_id } })
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

  getBotFavorite(client_id) {
    return request.get('/bot/favorite', { params: { client_id } })
  }

  addBotFavorite(client_id, id) {
    return request.post('/bot/favorite', {}, { params: { client_id, id } })
  }

  deleteBotFavorite(client_id, id) {
    return request.delete('/bot/favorite', { params: { client_id, id } })
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

  getBroadcast(client_id) {
    return request.get('/message/broadcast', { params: { client_id } })
  }

  postBroadcast({ client_id, category, data }) {
    return request.post('/message/broadcast', { client_id, category, data })
  }

  deleteBroadcast({ client_id, message_id }) {
    return request.delete('/message/broadcast', { params: { client_id, message_id } })
  }

}

export default new APIS()




