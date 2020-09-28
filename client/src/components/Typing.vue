<template>
  <div class="user">
    <div class="select">
      <span
        v-for="(value,key,idx) in messageTypeMap"
        :key="idx"
        class="group-list"
        @click="clickChangeGroup(key)"
      >
        <i class="iconfont" v-if="key === activeBroadcastType">&#xe637;</i>
        <i class="iconfont" v-else>&#xe638;</i>
        <b>{{value}}</b>
      </span>
    </div>
    <div class="edit-main">
      <textarea class="send-text" v-if="activeBroadcastType==='PLAIN_TEXT'" placeholder="在此输入公告内容" v-model="content"></textarea>
      <div class="send-image" v-if="activeBroadcastType==='PLAIN_IMAGE'">
        <input type="file" @change="getFile"/>
        <span>点击选择图片</span>
      </div>
    </div>
  </div>
</template>

<script>
  import { mapState, mapActions } from 'vuex'

  export default {
    data() {
      return {
        content: '',
      }
    },
    computed: {
      ...mapState('message', ['broadcastModalType', 'broadcastModal', 'activeBroadcastType', 'messageTypeMap'])
    },
    methods: {
      clickChangeGroup(activeBroadcastType) {
        this.$DC('message', { activeBroadcastType })
      },
      getFile(e) {
        console.log(e.target.files[0])
      }
    }
  }

</script>

<style lang="scss" scoped>
  .user {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  .send-text {
    flex: 1;
    resize: none;
    padding: 12px;
    line-height: 22px;
    color: #4C4471;
    background: #F8F9FE;
  }


  .send-image {
    text-align: center;
    position: relative;

    input {
      position: absolute;
      left: 0;
      right: 0;
      top: 0;
      bottom: 0;
      opacity: 0;
    }

    span {
      line-height: 174px;
    }
  }


  .edit-main {
    display: flex;
    flex-direction: column;
    height: 174px;
    width: 100%;
    background: #F8F9FE;
    border-radius: 12px;
    overflow: hidden;
  }

  .select {
    cursor: pointer;
    user-select: none;
    display: flex;
    margin: 2px 0 20px 5px;
  }

  .group-list {
    margin-right: 50px;
    font-size: 14px;
    color: #4C4471;
    display: flex;
    align-items: center;

    &:last-child {
      margin-right: 0;
    }

    b {
      font-weight: normal;
      margin-left: 10px;
    }
  }

  @media screen and (max-width: $adaptWidth) {

    .user {
      padding: 0 24px 20px 24px;
      flex-direction: column;
      align-items: flex-start;
    }

    .select {
      padding: 0;
      margin-top: 20px;

      .group-list {
        padding: 0;
        margin: 0;

        b {
          margin-left: 20px;
        }

        &:last-child {
          margin-left: 40px;
        }
      }
    }
  }
</style>
