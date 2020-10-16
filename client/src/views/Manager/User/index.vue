<template>
  <div class="users">
    <Tabs :current-state="currentState" :state-list="stateList" :toggle-state="toggleState" />
    <ul class="main-container">
      <li :class="['item', currentState==='block' && 'item-block']">
        <span>名称</span>
        <span>Mixin ID</span>
        <span>加入时间</span>
        <span v-if="currentState==='block'">屏蔽时间</span>
        <span class="btn">操作</span>
      </li>
      <li v-for="user in userList" :class="['item', currentState==='block' && 'item-block']">
        <div class="user-name">
          <Avatar :user_info="user" size="38" />
          <span>{{user.full_name}}</span>
        </div>
        <span class="user-number">{{user.identity_number}}</span>
        <span class="user-time">{{user.created_at}}</span>
        <span class="user-time" v-if="currentState==='block'">{{user.block_time}}</span>
        <button @click="clickBlockUser(user)">{{user.block_time ? "恢复":"操作"}}</button>
      </li>
    </ul>
  </div>
</template>

<script>
  import Tabs from '@/components/Tabs'
  import { mapState, mapActions } from 'vuex'
  import Avatar from "../../../components/Avatar"

  export default {
    name: "User",
    components: { Tabs, Avatar },
    computed: {
      ...mapState('users', ['stateList', 'currentState', 'userList'])
    },
    watch: {
      currentState() {
        this.updateUserList()
      }
    },
    methods: {
      ...mapActions('users', ['toggleState', 'getUserList', 'updateUserStatus']),
      updateUserList() {
        this.getUserList(this.currentState)
      },
      clickBlockUser(user) {
        const { user_id, block_time } = user
        const status = block_time ? "normal" : "block"
        const currentStatus = block_time ? "block" : "normal"
        const title = block_time ? "确认恢复？" : "确认拉黑？"
        this.$confirm(title, async () => {
          let req = await this.updateUserStatus({ user_id, status })
          if (req === "ok") {
            this.$message("操作成功")
            this.getUserList(currentStatus)
          }
        })
      }
    },
    mounted() {
      this.updateUserList()
    }
  }
</script>

<style lang="scss" scoped>

  .users {
    display: flex;
    width: 100%;
    flex-direction: column;
    color: #4C4471;
    overflow: hidden;
    height: 100%;
  }

  .main-container {
    @include main-container();
  }

  .item {
    display: grid;
    grid-template-columns: 1fr 100px 120px 100px;
    align-items: center;
    text-align: center;
    height: 80px;
    position: relative;
    background-color: #F8F9FE;
    padding: 0 32px;

    transition: all .1s;
    border-radius: 4px;
    overflow: hidden;

    &:nth-child(2n-1) {
      background-color: #fff;
    }

    &:hover {
      background: #E9E7FC;

      .session {
        background: #4C4471;
        color: #fff;
      }
    }

    img {
      width: 38px;
      height: 38px;
      border-radius: 50%;
      margin-right: 12px;
    }

    .message {
      margin-left: 12px;
      font-size: 14px;
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;

      p,
      i {

        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }

      p {
        color: #4C4471;
      }

      i {
        color: #A5A7C8;
      }
    }

    .date {
      display: flex;
      flex-direction: column;
      align-items: center;
      margin-left: 20px;

      p {
        font-size: 12px;
        color: #A5A7C8;
      }

      span {
        margin-top: 8px;
        font-size: 14px;
        color: #F05366;

        &.ok {
          color: #2FC273;
        }
      }
    }

    .session {
      font-size: 14px;
      padding: 5px 14px;
      cursor: pointer;
      color: #4C4471;
      border: 1px solid #4C4471;
      box-sizing: border-box;
      border-radius: 20px;
      margin: 0 32px;
    }
  }

  .item-block {
    grid-template-columns: 1fr 100px 140px 140px 100px;
  }

  .user-name {
    display: flex;
    align-items: center;
    text-align: center;
    margin: 2px 0 10px;
    font-weight: 500;
    font-size: 16px;
    overflow: hidden;

    .avatar {
      min-width: 38px;
    }


    span {
      text-overflow: ellipsis;
      white-space: nowrap;
      padding: 0 20px;
      overflow: hidden;
    }
  }

  .user-number {
    font-size: 14px;
    color: #B1B4CF;
  }

  .user-time {
    font-size: 12px;
    color: #B1B4CF;
  }

  button {
    height: 30px;
    width: 76px;
    color: #fff;
    background-color: #5B73A0;
    border-radius: 20px;
    margin: 0 auto;

    &:hover {
      background: #396AFF;
    }
  }


  @media screen and (max-width: $adaptWidth) {
    .users {
      padding: 8px 20px 0 20px;
      margin: 0;
    }
    .main-container {
      padding: 0;
      border-radius: 8px 8px 0 0;
    }

  }

</style>
