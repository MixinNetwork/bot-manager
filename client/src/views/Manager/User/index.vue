<template>
  <div class="users">
    <Tabs :current-state="currentState" :state-list="stateList" :toggle-state="toggleState" />
    <ul class="main-container">
      <template>
        <li v-for="user in userList" class="user-item">
          <img :src="user.avatar_url">
          <span class="user-name">{{user.full_name}}</span>
          <span class="user-number">Mixin ID: {{user.identity_number}}</span>
          <span class="user-time">加入时间：{{user.created_at}}</span>
          <span class="user-time" v-if="user.block_time">屏蔽时间：{{user.block_time}}</span>
          <button @click="clickBlockUser(user)">{{user.block_time ? "恢复":"操作"}}</button>
        </li>
      </template>
    </ul>
  </div>
</template>

<script>
  import Tabs from '@/components/Tabs'
  import { mapState, mapActions } from 'vuex'

  export default {
    name: "User",
    components: { Tabs },
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
        // if (user.block_time) {
        //   this.$confirm("确认恢复？", async () => {
        //     let req = await this.updateUserStatus({ user_id: user.user_id, status: "normal" })
        //     if (req === "ok") {
        //       this.$message("操作成功")
        //       this.getUserList("block")
        //     }
        //   })
        // } else {
        this.$confirm(title, async () => {
          let req = await this.updateUserStatus({ user_id, status })
          if (req === "ok") {
            this.$message("操作成功")
            this.getUserList(currentStatus)
          }
        })
        // }
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

    display: flex;
    flex-wrap: wrap;
  }

  .user-item {
    width: 223px;
    height: 226px;
    margin: 20px;
    border-radius: 12px;
    background-color: #F8F9FE;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;


    img {
      width: 48px;
      height: 48px;
      border-radius: 50%;
    }

    .user-name {
      white-space: nowrap;
      padding: 0 20px;
      width: 100%;
      overflow: hidden;
      text-overflow: ellipsis;
      text-align: center;
      margin: 2px 0 10px;
      font-weight: 500;

      font-size: 16px;
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
      margin-top: 26px;
      height: 30px;
      width: 76px;
      color: #fff;
      background-color: #5B73A0;
      border-radius: 20px;

      &:hover {
        background: #396AFF;
      }
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
