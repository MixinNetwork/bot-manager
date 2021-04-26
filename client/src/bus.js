import Vue from 'vue'

export default new Vue({
  methods: {
    toAuth(uri = '') {
      window.localStorage.clear()
      uri = uri || window.location.pathname
      let url = `https://mixin-www.zeromesh.net/oauth/authorize?client_id=${process.env.VUE_APP_CLIENT_ID}&scope=${process.env.VUE_APP_SCOPE}&response_type=code&return_to=`
      if (window.location.href === url) return
      window.location.href = url
    }
  }
})