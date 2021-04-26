<template>
  <div
    class="avatar"
    :style="`width:${size}px;height:${size}px`"
  >
    <img v-if="user_info.avatar_url" :src="user_info.avatar_url" alt="" class="logo">
    <span
      v-else
      class="logo"
      :style="`background:${bgColor};border-color:transparent;line-height:${size}px`"
    >{{user_info.full_name.slice(0,1)}}
    </span>
  </div>
</template>

<script>
  import { getAvatarColorById } from "../assets/js/color"

  export default {
    name: "Avatar",
    props: {
      user_info: {
        type: Object,
        default() {
          return {}
        }
      },
      size: {
        type: String,
        default: "36"
      }
    },
    computed: {
      bgColor() {
        const { avatar_url, user_id, client_id } = this.user_info
        if (avatar_url) return ""
        return getAvatarColorById(user_id || client_id)
      }
    },
  }
</script>

<style lang="scss" scoped>


  .logo {
    display: block;
    width: 100%;
    height: 100%;
    border-radius: 20px;
    text-align: center;
    color: #ffffff;
  }
</style>
