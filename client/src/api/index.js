import request from './config'

class APIS {
  authenticate(code) {
    return request.get('/user/login', {params: {code}})
  }

  getGrowthUserList(status) {
    return request.get('/growth/user', {params: {status}})
  }

  updateGrowthStatus({user_id, status, auditOpinion}) {
    return request.post('/growth/user', {user_id, status, auditOpinion})
  }

  getData() {
    return request.get('/data')
  }

  uploadFile(file) {
    let formData = new FormData()
    formData.append('file', file)
    return request.post('/uploadFile', formData, {heaaders: {'Content-Type': 'multipart/form-data'}})
  }

  sendSpecificUser(data, group) {
    return request.post('/sendSpecificUser', {data, group})
  }

  addBot({client_id, session_id, private_key}) {
    return request.post('/bot', {client_id, session_id, private_key})
  }

  getBot() {
    return request.get('/bot')
  }


}

export default new APIS()




