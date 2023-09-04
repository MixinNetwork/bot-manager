<template>
  <nav>
    <header class="bot-container">
      <BotItem @toggleList="toggleList" :is-head="true" :is-add="!active_bot.client_id" :bot-item="active_bot"/>
    </header>
    <ul v-if="showBotList" class="select-bot-list">
      <BotItem @toggleList="toggleList" v-for="(item, idx) in bot_list" :key="idx" :bot-item="item"/>
      <BotItem @toggleList="toggleList" :is-add="true"/>
    </ul>
    <ul class="menus">
      <li v-for="(item, idx) in navMenus"
          :key="idx"
          :class="activeMenus === idx ? 'active' : ''"
          @click="clickNav(idx)"
      >
        <i class="iconfont">{{navIcon[[idx]]}}</i>
        <span class="menu-item">{{item}}</span>
      </li>
    </ul>

    <footer class="avatar" @click.stop="clickAvatar">
      <img :src="user_info.avatar_url"/>
      <div class="user-info">
        <div class="user-name">{{user_info.full_name}}</div>
        <div class="user-number">{{user_info.identity_number}}</div>
      </div>
      <div v-if="showMenus" class="bottom-info middle">
        <div :class="['bottom-more', slider ?'bottom-more-active':'']">
          <div class="bottom-button-list">
            <div @click.stop="clickLogout" class="bottom-button-item">
              <span>登出</span>
            </div>
          </div>
        </div>
      </div>
    </footer>
    <AddBot v-show="show_add_bot"/>
  </nav>
</template>

<script>

  import { mapState, mapActions } from 'vuex'
  import AddBot from '../AddBot'
  import BotItem from './BotItem'

  export default {
    name: "Nav",
    components: { AddBot, BotItem },
    data() {
      return {
        showMenus: false,
        slider: false,
        navIcon: ['\ue62b', '\ue62c', '\ue62d', '\ue62e'],
        navMenus: ["数据统计", "用户管理", "消息管理", "设置"],
        routerList: ['/data', '/user', '/messages', '/setting'],
        isBlurAndNoContent: true,
        showBotList: false
      }
    },
    computed: {
      ...mapState('system', ['activeMenus']),
      ...mapState('user', ['user_info', 'bot_list', 'active_bot', 'show_add_bot']),
    },
    methods: {
      ...mapActions('user', ['getBotList']),
      clickNav(idx) {
        let { $route, routerList } = this
        this.$DC({ activeMenus: idx })
        if (!($route.path === routerList[idx])) this.$router.push(routerList[idx])
      },
      clickLogout() {
        window.localStorage.clear()
        setTimeout(() => {
          window.location.reload()
        }, 100)
      },
      clickAvatar() {
        if (this.showMenus) {
          this.slider = false
          setTimeout(() => {
            this.showMenus = false
          }, 100)
          document.documentElement.onclick = null
        } else {
          this.showMenus = true
          setTimeout(() => {
            this.slider = true
          }, 20)
          document.documentElement.onclick = () => {
            this.clickAvatar()
          }
        }
      },
      toggleList() {
        if (this.showBotList) {
          this.showBotList = false
          window.onclick = null
        } else {
          this.showBotList = true
          window.onclick = () => {
            this.showBotList = false
          }
        }
      }
    },
    mounted() {
      this.$DC({ activeMenus: this.routerList.findIndex(item => item === this.$route.path) })
      this.getBotList()
    }
  }
</script>

<style lang="scss" scoped>
  nav {
    margin-right: 10px;
    background-color: #fff;
    border-radius: 30px;
    color: #505887;
    display: flex;
    flex-direction: column;
    align-items: center;

    position: relative;
    user-select: none;
    height: calc(100% - 20px);
    box-sizing: border-box;
    box-shadow: 0 4px 20px rgba(0, 18, 179, 0.08);
  }

  .bot-container {
    width: 240px;
    margin: 40px 32px;
    padding: 18px 0;
    height: 72px;
    background-color: #fff;
    border: 1px solid #E7F0F7;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    cursor: pointer;
    position: relative;
    z-index: 3;

  }

  .user-info {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  .user-number {
    font-size: 12px;
    color: #b1b4cf;
  }

  .user-name {
    font-size: 14px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 130px;
    margin-bottom: 1px;
  }

  .select-bot-list {
    position: absolute;
    top: 112px;
    left: 32px;
    z-index: 2;
    width: 240px;
    box-shadow: 0px 8px 40px rgba(0, 18, 179, 0.08);
    border-radius: 0 0 8px 8px;
    background-color: #fff;
    max-height: 284px;
    overflow: auto;
    padding: 12px 0;

    &::-webkit-scrollbar {
      width: 5px;
      height: 1px;
    }

    &::-webkit-scrollbar-thumb {
      background: #E0E0E0;
      border-radius: 20px;
    }

    &::-webkit-scrollbar-track {
      background: #fff;
    }

  }

  .menus {
    margin-top: 8px;
    display: flex;
    flex-direction: column;
    width: 100%;
    overflow: hidden;

    .menu-item {
      display: inline;
    }

    .iconfont {
      margin-right: 18px;
    }

    li {
      line-height: 54px;
      margin: 0 32px 15px 32px;
      padding-left: 28px;
      cursor: pointer;
      background: #fff;
    }

    .active {
      color: #396AFF;
      background: #EEF5FF;
      border-radius: 12px;
      position: relative;
    }
  }

  .avatar {
    position: absolute;
    bottom: 40px;
    cursor: pointer;
    display: flex;
    left: 32px;
    align-items: center;

    img {
      width: 46px;
      height: 46px;
      border-radius: 50%;
      border: 3px solid #fff;
      margin-right: 12px;
    }
  }


  .bottom-info {
    position: absolute;
    width: 17.5rem;
    cursor: pointer;
    bottom: 0;
    border: 0;
    background: #fff;
  }

  .user-name {
    max-width: 180px;
  }

  .bottom-more {
    position: absolute;
    top: -.5rem;
    left: -0.5rem;
    z-index: -1;

    transition: all 0.1s ease-in;
    opacity: 0;
    color: #4c3371;

    box-shadow: 0 0.125rem 1rem 0 rgba(47, 48, 50, 0.4);

    padding: 0.5rem 1rem;

    background-color: #fff;
    border-radius: 0.375rem;
    overflow: hidden;

    &.bottom-more-active {
      z-index: 2;
      opacity: 1;
      top: 0.2rem;
    }
  }

  @media screen and (max-width: $adaptWidth) {
    nav {
      background-color: transparent;
      min-width: 0;
      width: 0;
      margin: 0;
    }

    .logo {
      position: fixed;
      top: 10px;
      left: 20px;
      margin: 0;
      width: 20px;
      height: 20px;
      color: #4C4471;
      font-size: 17px;
    }

    .title {
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      color: #4C4471;
      margin: 0;
      padding-left: 50px;
      line-height: 40px;

      &::after {
        content: '';
        position: absolute;
        background-color: #F8F9FE;
        left: 0;
        right: 0;
        height: 40px;
        top: 0;
        z-index: -1;
      }
    }
    .avatar {
      position: fixed;
      top: 4px;
      bottom: initial;
      right: 20px;
      z-index: 1;

      img {
        width: 32px;
        height: 32px;
      }
    }
    .bottom-more {
      left: -1rem;
    }

    .menus {
      position: fixed;
      bottom: 0;
      left: 0;
      right: 0;
      height: 60px;
      display: flex;
      /*z-index: 10;*/
      background: #F8F9FE;
      flex-direction: row;

      .menu-item {
        display: none;
      }

      .iconfont {
        display: inline;
        color: #A5A7C8;
        font-size: 24px;
      }

      li {
        margin: 0;
        padding: 0;
        flex: 1;
        text-align: center;
        line-height: 60px;
      }

      .active {
        border-radius: 0;

        .iconfont {
          color: #4C4471;
        }

        .top-border,
        .bottom-border {
          display: none;
        }
      }
    }

    .search {
      position: fixed;
      left: 20px;
      right: 20px;
      top: 0;
      z-index: 1;

      input {
        padding: 0;
        width: 0;
      }

      .icon-holder {
        display: none;
      }

      .icon-search {
        position: absolute;
        left: initial;
        top: 12px;
        right: 80px;
        color: #4C4471;
      }
    }

    .search-active {
      input {
        width: 100%;
        padding-left: 16px;
        margin-top: 2px;
      }


    }


  }

</style>