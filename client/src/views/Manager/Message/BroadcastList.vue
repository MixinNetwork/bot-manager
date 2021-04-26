<template>
  <ul class="item-list">
    <li class="item">
      <span>内容</span>
      <span>发布人</span>
      <span>发布时间</span>
      <span>公告类型</span>
      <span>操作</span>

    </li>
    <li v-for="(item,idx) in broadcastList" :key="idx" class="item">
      <span>{{item.data}}</span>
      <span>{{item.full_name}}</span>
      <span>{{item.created_at}}</span>
      <span>{{item.category | type}}</span>
      <div class="btn">
        <button @click="clickCheck(item)">查看</button>
        <button @click="clickDelete(item)">撤回</button>
      </div>
    </li>
    <BroadcastModal />
  </ul>
</template>


<script>
  import { createNamespacedHelpers } from 'vuex'

  const { mapState, mapActions } = createNamespacedHelpers('message')
  import BroadcastModal from "./BroadcastModal"

  const TYPES = {
    PLAIN_TEXT: '文字消息',
    PLAIN_IMAGE: "图片",
    APP_BUTTON: "按钮"
  }

  const STATUS = {
    finished: "已发布"
  }
  export default {
    components: { BroadcastModal },
    computed: {
      ...mapState(['broadcastList'])
    },
    filters: {
      type: type => TYPES[type],
      status: status => STATUS[status]
    },
    methods: {
      ...mapActions(['getBroadcast', 'deleteBroadcast']),
      clickCheck(item) {
        this.$DC('message', { broadcastModalType: "check", broadcastModal: true, activeBroadcast: item })
      },
      async clickDelete(item) {
        this.$confirm("确认撤回？", async () => {
          this.$DC('message', { activeBroadcast: item })
          let data = await this.deleteBroadcast()
          console.log(data)
          if (data === 'ok') {
            this.getBroadcast()
          }
        })
      }
    },
    async mounted() {
      let t = await this.getBroadcast()
      console.log(t)
    }
  }
</script>


<style lang="scss" scoped>

  .item {
    display: grid;
    grid-template-columns: 14fr 5fr 9fr 4fr 8fr;
    grid-column-gap: 26px;
    text-align: center;
    align-items: center;
    padding: 26px;
    height: 80px;
    position: relative;
    background-color: #F8F9FE;
    font-size: 14px;
    transition: all .1s;
    border-radius: 4px;
    overflow: hidden;
    color: #505887;

    span {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    &:nth-child(2n-1) {
      background-color: #fff;
    }

    .btn {
      margin: 0 auto;
      width: 130px;

      button {
        width: 56px;
        line-height: 30px;
        border-radius: 20px;
        color: #fff;

        &:first-child {
          background-color: #396AFF;
          margin-right: 16px;
        }

        &:last-child {
          background-color: #5B73A0;
        }
      }
    }

    &:first-child {
      color: #B1B4CF;
    }
  }

  .action-button {
    position: absolute;
    top: -20px;
    right: 0;
    width: 116px;
    line-height: 30px;
    background-color: #5B73A0;
    color: #ffffff;
    border-radius: 20px;
    font-size: 16px;
    text-align: center;
    z-index: 10;

  }

  @media screen and (max-width: $adaptWidth) {

  }
</style>
