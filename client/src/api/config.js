import axios from 'axios'

const request = axios.create({
  baseURL: process.env.VUE_APP_SERVER
})

let retry = 0

request.interceptors.request.use(config => {
  config.headers.Authorization = "Bearer " + window.localStorage.getItem("token")
  let t = Number(new Date())
  config.url.includes('?') ? (config.url = config.url + '&t=' + t) : (config.url = config.url + '?t=' + t)
  return config
})

function backOff() {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve()
    }, 500)
  })
}

request.interceptors.response.use(res => {
  retry = 0
  let {data, code, description} = res.data
  if (data) return data
  _vm.$message(description)
  if ([401, 403].includes(code)) _vm.$bus.toAuth()
  return undefined
}, async err => {
  await backOff()
  if (retry > 5) throw new Error('AjaxError')
  retry++
  return request(err.config)
})


export default request
