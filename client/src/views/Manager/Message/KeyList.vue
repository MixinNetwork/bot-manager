<template>
  <ul>
    <li class="item">
      <span>关键字</span>
      <span>回复类型</span>
      <span>回复内容</span>
      <span class="btn">操作</span>
    </li>
    <li v-for="(item,idx) in keyList" :key="idx" :class="['item']">
      <span>{{item.key}}</span>
      <span>{{item.category | type}}</span>
      <span>{{item.content}}</span>
      <div class="btn">
        <button @click="clickCheck(item)">详情</button>
        <button @click="clickEdit(item)">编辑</button>
        <button @click="clickDel(item)">删除</button>
      </div>
    </li>
    <KeyModal/>
  </ul>
</template>


<script>
  import { createNamespacedHelpers } from 'vuex'
  import KeyModal from "./KeyModal"

  const { mapState } = createNamespacedHelpers('message')
  const TYPES = {
    PLAIN_TEXT: '文字消息',
    PLAIN_IMAGE: "图片",
    BUTTON_GROUP: "按钮"
  }
  const STATUS = {
    finished: "已发布"
  }
  export default {
    components: { KeyModal },
    computed: {
      ...mapState(['keyList'])
    },
    filters: {
      type: type => TYPES[type],
      status: status => STATUS[status]
    },
    methods: {
      clickCheck(item) {
        this.$DC('message', { keyModalType: "check", keyModal: true })
      },
      clickEdit(item) {
        this.$DC('message', { keyModalType: "edit", keyModal: true })
      },
      clickDel(item) {
        this.$confirm("确认删除？", () => {

        })
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
        flex: 1;
      }

      &:nth-child(2) {
        width: 100px;
      }

      &:nth-child(3) {
        flex: 2;
      }

      &:nth-child(4) {
        width: 216px;
      }
    }

    .btn {
      margin: 0 10px;
      width: 216px;
      display: flex;
      justify-content: center;

      button {
        width: 56px;
        line-height: 30px;
        border-radius: 20px;
        color: #fff;
        margin-right: 16px;

        &:nth-child(1) {
          background-color: #396AFF;
        }

        &:nth-child(2) {
          background-color: #5B73A0;
        }

        &:nth-child(3) {
          background-color: #EC5A75;
          margin-right: 0;
        }


      }
    }

  }


  @media screen and (max-width: $adaptWidth) {

  }
</style>