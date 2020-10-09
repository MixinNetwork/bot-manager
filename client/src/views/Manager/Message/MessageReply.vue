<template>
  <Modal :show="replyModal">
    <div class="modal-container">
      <header>
        <Avatar class="avatar" :user_info="activeContact" size="40" />
        <div class="user-info">
          <p>{{activeContact.full_name}}</p>
          <i>Mixin ID：{{activeContact.identity_number}}</i>
        </div>
        <i class="iconfont close-btn" @click="clickClose">&#xe636;</i>
      </header>

      <template v-if="!file">
        <ul ref="messageDom" class="message-container">
          <li
            v-for="(item,idx) in activeContact.messages"
            :key="idx"
            :class="item.status!=='sent' ? 'my-message' :''"
          >
            <template v-if="item.category === 'PLAIN_TEXT'">
              <div class="message-item">
                <span> {{item.data}}</span>
                <span class="time-holder">{{item.created_at}}
                  <i v-if="item.status === 'pending'" :class="`iconfont ${item.status}`">&#xe63a;</i>
                  <i v-else-if="item.status !== 'sent'" :class="`iconfont ${item.status}`">&#xe63b;</i>
                </span>
              </div>
            </template>

            <template v-else-if="item.category === 'PLAIN_IMAGE'">
              <div class="img-item">
                <img @load="messageContainerToBottom" :src="item.data" />
                <span class="time-holder">{{item.created_at}}
                  <i v-if="item.status === 'pending'" :class="`iconfont ${item.status}`">&#xe63a;</i>
                  <i v-else-if="item.status !== 'sent'" :class="`iconfont ${item.status}`">&#xe63b;</i>
                </span>
              </div>
            </template>
            <template v-else>
              <div class="message-item">
                <span>[ 未识别的消息类型，请在手机端查看 ]</span>
              </div>
            </template>

          </li>
        </ul>

        <footer>
          <i class="iconfont upload-img">&#xe643;
            <input type="file" title="" @change="imgChanged">
          </i>
          <input type="text" ref="input" v-model="message" @keydown.enter="clickSendMessage">
          <i v-if="!onceClick" class="iconfont send-icon" @click="clickSendMessage">&#xe644;</i>
          <CLoading v-else class="text-c-loading" color="#397EE4" />
        </footer>
      </template>

      <div v-else class="send-img">
        <div class="head">
          <i class="iconfont img-close" @click="clickDeleteImg">&#xe642;</i>
          <span>Preview</span>
        </div>

        <div class="img">
          <img ref="imgFile" :src="file">
        </div>

        <i v-if="!onceClick" class="iconfont img-next" @click="sendImg">&#xe641;</i>
        <CLoading v-else class="img-c-loading" color="#fff" />
      </div>

    </div>
  </Modal>
</template>

<script>
  import Modal from '@/components/Modal'
  import CLoading from '@/components/CLoading'
  import { mapState, mapActions } from 'vuex'
  import Avatar from "../../../components/Avatar"

  export default {
    name: "Reply",
    components: { Avatar, Modal, CLoading },
    data() {
      return {
        file: null,
        message: '',
        _file: null,
        onceClick: false
      }
    },
    computed: {
      ...mapState('message', ['replyModal', 'activeContact'])
    },
    methods: {
      ...mapActions('message', ['sendMessage']),
      clickChangeGroup(broadcastGroup) {
        this.$DC('message', { broadcastGroup })
      },
      clickClose() {
        this.$DC('message', { replyModal: false })
      },
      clickSendMessage() {
        if (this.onceClick) return
        this.onceClick = true
        this.sendMessage({ data: this.message, category: 'PLAIN_TEXT' })
        this.message = ''
      },
      imgChanged(e) {
        let file = e.target.files[0]
        this.file = URL.createObjectURL(file)
        this._file = file
      },
      clickDeleteImg() {
        this.file = null
      },
      async sendImg() {
        if (this.onceClick) return
        this.onceClick = true
        let { _file: file } = this
        let { clientHeight: height, clientWidth: width } = this.$refs.imgFile
        let { size, type: mime_type } = file
        let { attachment_id } = await this.APIS.uploadFile(file)
        let data = { height, width, size, mime_type, attachment_id }
        this.sendMessage({ data, category: 'PLAIN_IMAGE' })
        this._file = this.file = null
      },
      messageContainerToBottom({ status } = {}) {
        let el = this.$refs.messageDom
        if (!el) return
        if (status) this.onceClick = false
        this.$nextTick(() => {
          el.scrollTop = el.scrollHeight
        })
      }
    },
    mounted() {
      this.$bus.$on('message', status => this.messageContainerToBottom(status))
      this.messageContainerToBottom()
      this.$refs.input.focus()
    }
  }

</script>

<style lang="scss" scoped>

  .modal-container {
    width: 700px;
    height: 600px;
    display: flex;
    flex-direction: column;
  }

  .iconfont {
    cursor: pointer;
  }

  header {
    display: flex;
    height: 70px;
    align-items: center;

    .avatar {
      margin: 0 12px 0 28px;
    }

    .user-info {
      flex: 1;

      p {
        color: #4C4471;
      }

      i {
        font-size: 12px;
        color: #A5A7C8;
      }

    }

    .close-btn {
      color: #EAEAEA;
      font-size: 24px;
      margin-right: 24px;
    }

  }

  .message-container {
    background-color: #F8F9FE;
    flex: 1;
    border-top: 1px solid #ECEFFF;
    border-bottom: 1px solid #ECEFFF;
    padding: 20px 20px 0 20px;

    display: flex;
    flex-direction: column;
    overflow: auto;

    li {
      margin-bottom: 12px;


      &.my-message {
        display: flex;
        justify-content: flex-end;

        .message-item {
          text-align: left;
          background-color: #C5EDFD;
        }

        .time-holder {
          display: flex;
          justify-content: flex-end;
          width: 70px;
        }

        .sent {
          display: none;
        }

        .pending,
        .delivered,
        .read {
          padding: 0;
          display: inline;
          font-size: 12px;
          margin-left: 5px;
        }

        .read {
          color: #357AF6;
        }

      }
    }

    .message-item {
      display: inline-block;
      background-color: #fff;
      box-shadow: 0 2px 4px rgba(61, 117, 227, 0.08);
      border-radius: 8px;
      padding: 10px 12px;
      line-height: 22px;
      max-width: 400px;

      .iconfont {
        font-size: 12px;
        margin-left: 6px;
        cursor: default;
      }
    }

    .time-holder {
      display: inline-block;
      width: 60px;
      text-align: right;
      font-size: 10px;
      color: #83919E;
      transform: translateY(3px);
      float: right;
      user-select: none;

      &::after {
        clear: both;
        content: '';
        display: block;
        width: 0;
        height: 0;
        visibility: hidden;
      }
    }

    li .img-item {
      max-width: 400px;

      position: relative;

      img {
        max-width: 100%;
        max-height: 100%;
      }

      .time-holder {
        position: absolute;
        bottom: 16px;
        right: 4px;
        background: rgba(0, 0, 0, 0.3);
        border-radius: 12px;
        text-align: center;
        color: #fff;
        line-height: 18px;
        width: initial;
        padding: 4px 6px;
      }
    }

  }

  footer {
    display: flex;
    align-items: center;
    height: 52px;


    .upload-img {
      margin-left: 15px;
      position: relative;
      overflow: hidden;
      display: inline-block;
      font-size: 20px;

      input {
        position: absolute;
        top: 0;
        left: -100px;
        right: 0;
        bottom: 0;
        opacity: 0;
        cursor: pointer;
      }
    }

    .text-c-loading,
    .send-icon {
      margin-right: 10px;
      color: #397EE4;
      font-size: 20px;
    }

    input {
      flex: 1;
      margin: 0 28px;
    }
  }


  .send-img {
    padding: 20px 30px;
    text-align: center;

    .head {
      display: flex;
      align-items: center;
      color: #4c3371;
      margin-bottom: 16px;


      .img-close {
        font-size: 20px;
        margin-right: 16px;
      }
    }

    .img {
      display: flex;
      justify-content: center;
      align-items: center;
      width: 100%;
      height: 338px;
      background-color: #F8F9FE;
      border: 1px dashed #E5E5E5;
      border-radius: 8px;

      img {
        max-height: 80%;
        max-width: 80%;

      }
    }

    .img-c-loading,
    .img-next {
      display: inline-block;
      margin-top: 32px;
      width: 50px;
      height: 50px;
      line-height: 50px;
      background-color: #4C4471;
      color: #fff;
      border-radius: 50%;
      font-size: 20px;
      padding-left: 4px;

      &:hover {
        opacity: 0.9;
      }
    }

    .img-c-loading {
      display: flex;
      align-items: center;
      margin: 32px auto 0 auto;
      padding: 0;
    }
  }

  @media screen and (max-width: $adaptWidth) {
    .modal-container {
      width: calc(100vw - 36px);
      height: calc(100vh - 100px);
    }

    header {
      .user-info {
        overflow: hidden;

        p {
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }
      }
    }
  }


</style>