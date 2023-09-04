<template>
  <div class="setting">
    <Tabs :state-list="stateList" :current-state="currentState" />
    <ul class="main-container">
      <li v-for="(item,idx) in favorite" :key="idx" class="card-item">
        <i class="iconfont close-btn" @click="clickDel(item)">&#xe636;</i>
        <div class="card-content">
          <Avatar :user_info="item" />
          <div class="card-text-info">
            <span>{{item.full_name}}</span>
            <span>{{item.identity_number}}</span>
          </div>
        </div>
      </li>
      <li class="card-item-add" @click="clickAddFavorite" v-if="favorite.length <3">
        <i class="iconfont add-btn">&#xe630;</i>
        <span>添加</span>
      </li>
    </ul>
    <AddFavorite />
  </div>
</template>

<script>
  import Tabs from "@/components/Tabs"
  import Avatar from '@/components/Avatar'
  import AddFavorite from './addFavorite'
  import { createNamespacedHelpers } from 'vuex'

  const { mapActions, mapState } = createNamespacedHelpers('user')

  export default {
    components: { Tabs, AddFavorite, Avatar },
    data() {
      return {
        stateList: { share: "关联分享" },
        currentState: "share"
      }
    },
    computed: {
      ...mapState(['favorite'])
    },
    methods: {
      ...mapActions(['getBotFavorite', 'delBotFavorite']),
      clickDel(item) {
        this.$confirm("确认删除？", async () => {
          const t = await this.delBotFavorite(item.user_id)
          if (t === 'ok') {
            this.$message('删除成功')
            this.getBotFavorite()
          }
        })
      },
      clickAddFavorite() {
        this.$DC('user', { addFavoriteModal: true })
      }
    },
    mounted() {
      this.getBotFavorite()
    }
  }
</script>

<style lang="scss" scoped>

  .setting {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .main-container {
    @include main-container();
    padding: 32px 24px;
    overflow: hidden;
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
  }

  .card-item {
    margin: 0 10px;
    position: relative;
    overflow: hidden;

    .avatar {
      min-width: 36px;
    }
  }

  .card-content {
    margin-top: 10px;
    margin-right: 10px;
    display: flex;
    align-items: center;
    padding: 0 20px;
    background-color: #F8F9FE;
    height: 75px;
    border-radius: 12px;
  }

  .card-text-info {
    display: flex;
    flex-direction: column;
    margin-left: 16px;
    overflow: hidden;

    span {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;

      &:first-child {
        font-size: 14px;
        color: #505887;
      }

      &:last-child {
        font-size: 12px;
        color: #B1B4CF;
      }
    }
  }

  .close-btn {
    position: absolute;
    right: 0px;
    top: 0px;
    font-size: 20px;
    color: #EAEAEA;
    cursor: pointer;
  }

  .card-item-add {
    margin: 10px;
    height: 75px;
    border: 1px dashed #C4C7DB;
    border-radius: 12px;
    display: flex;
    justify-content: center;
    align-items: center;
    color: #B1B4CF;
    cursor: pointer;

    .add-btn {
      font-size: 16px;
      margin-right: 8px;
    }

    span {
      font-size: 14px;
      line-height: 20px;
      letter-spacing: 2px;
    }

  }


</style>
