<template>
  <ul>
    <li v-for="(item,idx) in contactList" :key="idx" :class="['item']">
      <Avatar :user_info="item" size="38"/>
      <div class="message">
        <p>{{item.full_name}}</p>
        <i>{{item.messages | message}}</i>
      </div>
      <div class="date">
        <p>{{item.messages[item.messages.length-1].created_at}}</p>
        <span v-if="item.messages[item.messages.length-1].status!=='sent'" class="ok">已回复</span>
        <span v-else>未回复</span>
      </div>
      <span class="session" @click="clickSession(item)">会话</span>
    </li>
    <MessageReply v-if="replyModal"/>
  </ul>
</template>


<script>
  import { createNamespacedHelpers } from 'vuex'
  import MessageReply from "./MessageReply"
  import Avatar from "../../../components/Avatar"

  const { mapState } = createNamespacedHelpers('message')
  export default {
    components: { Avatar, MessageReply },
    computed: {
      ...mapState(['contactList', 'replyModal'])
    },
    methods: {
      clickSession(activeContact) {
        this.$DC('message', { replyModal: true, activeContact })
      }
    },
    filters: {
      message(messages) {
        let message = messages[messages.length - 1]
        if (message.category === 'PLAIN_TEXT') return message.data
        if (message.category === 'PLAIN_IMAGE') return `[图片]`
        return `[未识别的消息类型，请在手机端查看]`
      }
    },
  }
</script>


<style lang="scss" scoped>
  .item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 80px;
    position: relative;
    background-color: #F8F9FE;
    padding-left: 32px;

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


  @media screen and (max-width: $adaptWidth) {

    .item {
      padding: 0 20px;


      .session {
        opacity: 0;
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
      }
    }
  }
</style>