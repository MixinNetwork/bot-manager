<template>
  <ul>
    <li class="item">
      <span>内容</span>
      <span>发布人</span>
      <span>发布时间</span>
      <span>公告类型</span>
      <span>状态</span>
      <span class="btn">操作</span>
    </li>
    <li v-for="(item,idx) in broadcastList" :key="idx" :class="['item']">
      <span>{{item.content}}</span>
      <span>{{item.full_name}}</span>
      <span>{{item.created_at}}</span>
      <span>{{item.category | type}}</span>
      <span>{{item.status | status}}</span>
      <div class="btn">
        <button @click="clickCheck">查看</button>
        <button>撤回</button>
      </div>
    </li>
    <BroadcastModal />
  </ul>
</template>


<script>
  import { createNamespacedHelpers } from 'vuex'

  const { mapState } = createNamespacedHelpers('message')
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
      clickCheck() {
        this.$DC('message', { broadcastModalType: "check", broadcastModal: true })
      }
    }
  }
</script>


<style lang="scss" scoped>

  .item {
    display: flex;
    align-items: center;
    height: 80px;
    position: relative;
    background-color: #F8F9FE;
    padding: 0 32px;
    font-size: 14px;

    transition: all .1s;
    border-radius: 4px;
    overflow: hidden;

    color: #505887;

    &:nth-child(2n-1) {
      background-color: #fff;
    }

    &:first-child {
      height: 50px;

      span {
        color: #B1B4CF !important;
      }
    }

    span {
      text-align: center;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      margin: 0 10px;

      &:nth-child(1) {
        flex: 3;
      }

      &:nth-child(2) {
        flex: 1;
      }

      &:nth-child(3) {
        width: 150px;
      }

      &:nth-child(4) {
        width: 60px;
      }

      &:nth-child(5) {
        width: 50px;
        color: #2FC273;
      }
    }

    .btn {
      margin: 0 10px;
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
