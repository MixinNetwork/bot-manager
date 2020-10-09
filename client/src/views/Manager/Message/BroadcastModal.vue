<template>
  <Modal :show="broadcastModal">
    <div class="modal-container">
      <template v-if="broadcastModalType==='send'">
        <h3>发布公告</h3>
        <i class="iconfont close-btn" @click="clickClose">&#xe636;</i>
        <Typing />
        <div class="btns">
          <button @click="clickConfirmBtn()" class="btn send">发布</button>
        </div>
      </template>
      <div class="check" v-else-if="broadcastModalType==='check'">
        <h3>查看</h3>
        <i class="iconfont close-btn" @click="clickClose">&#xe636;</i>
        <div class="text">
          <label>内容</label>
          <span>【Mixin 公告】Mixin 不支持 BTS 分叉空投，如需领取建议通过 @BigONE 机器人划转至 BigONE 领取空投，@BigONE 机器人现已支持免费、实时划转 BTS，注意请在 2020 年 9 月 10 日 19:00（UTC +8）前划转至 BigONE 账户。</span>
        </div>
        <div class="text">
          <label>发布人</label>
          <span>Kelly</span>
        </div>
        <div class="text">
          <label>状态</label>
          <span class="green">已发布</span>
        </div>

        <div class="text">
          <label>发布时间</label>
          <span>2020-04-22 13:04:12</span>
        </div>
        <div class="btns">
          <button class="btn back">撤回</button>
        </div>
      </div>
    </div>
  </Modal>
</template>

<script>
  import Modal from '@/components//Modal'
  import { mapState, mapActions } from 'vuex'
  import Typing from '@/components/Typing'

  export default {
    name: "Broadcast",
    components: { Modal, Typing },
    data() {
      return {}
    },
    computed: {
      ...mapState('message', ['broadcastModalType', 'broadcastModal'])
    },
    methods: {
      ...mapActions('message', ['sendBroadcast']),
      clickConfirmBtn() {
        this.$confirm("确认发布？", async () => {

        })
      },
      // async clickBtn(state) {
      //   if (state === 'cancel') {
      //   } else {
      //     let data = encodeURIComponent(this.content)
      //     let t = await this.sendBroadcast(data)
      //     if (t === 'sending') {
      //       this.$message('发送中...')
      //       this.$DC('message', { broadcastModal: false })
      //       this.content = ''
      //     } else {
      //       this.$message('发送失败...')
      //     }
      //   }
      //
      // },
      clickClose() {
        return this.$DC('message', { broadcastModal: false })
      },
    },
  }
</script>

<style lang="scss" scoped>

  .modal-container {
    width: 520px;
    background-color: #fff;
    color: #505887;
    padding: 0 52px;
  }

  h3 {
    text-align: center;
    font-size: 22px;
    margin: 44px 0 30px 0;
    letter-spacing: 2px;
    font-weight: normal;
  }

  .close-btn {
    position: absolute;
    top: 20px;
    right: 20px;
    font-size: 20px;
    color: #EAEAEA;
    cursor: pointer;
  }

  .text {
    display: flex;
    flex-direction: column;
    height: 174px;
  }

  label {
    color: #B8BDC7;
    font-size: 14px;
    background: #F8F9FE;
  }


  .btns {
    margin: 56px 0 32px;
    display: flex;
    justify-content: center;

    .btn {
      width: 84px;
      height: 44px;
      border-radius: 20px;
      font-size: 16px;
      color: #fff;
    }

    .send {
      background: #396AFF;
    }

    .back {
      background-color: #5B73A0;
    }
  }

  .check {
    .text {
      flex-direction: row;
      height: initial;
      margin: 16px 0;

      label {
        min-width: 60px;
        margin-right: 16px;
        background-color: initial;
      }

      span {
        font-size: 14px;
      }
    }

    .btns {
      margin-top: 30px;
    }
  }

  .green {
    color: #2FC273;
  }


  @media screen and (max-width: $adaptWidth) {
    .modal-container {
      width: calc(100vw - 36px);
    }

    h3 {
      margin-bottom: 20px;
    }

    .text {
      padding: 0 24px 20px 24px;
    }


    .btns {
      margin-top: 60px;

    }
  }
</style>
