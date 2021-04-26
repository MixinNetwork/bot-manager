<template>
  <Modal :show="keyModal">
    <div class="modal-container">


      <div class="edit-modal" v-if="keyModalType==='edit'">
        <h3>编辑</h3>
        <i class="iconfont close-btn" @click="clickClose">&#xe636;</i>
        <div class="row">
          <span class="title">关键字</span>
          <input v-model="key" class="key-input" type="text" placeholder="请输入关键字（用空格分开）" />
        </div>
        <div class="row">
          <span class="title">回复内容</span>
          <Typing />

        </div>
        <div class="btns">
          <button @click="clickConfirmBtn()" class="btn save">保存</button>
        </div>
      </div>


      <div class="check" v-else-if="keyModalType==='check'">
        <h3>关键字回复详情</h3>
        <i class="iconfont close-btn" @click="clickClose">&#xe636;</i>
        <div class="text">
          <label>关键字</label>
          <span>{{key}}</span>
        </div>

        <div class="text">
          <label>回复内容</label>
          <span>{{activeContent}}</span>
        </div>
        <div class="btns">
          <button @click="clickEdit" class="btn edit">编辑</button>
          <button @click="clickDel" class="btn del">删除</button>
        </div>
      </div>
    </div>
  </Modal>
</template>

<script>
  import Modal from '@/components/Modal'
  import Typing from '@/components/Typing'
  import { mapState, mapActions } from 'vuex'

  export default {
    name: "Broadcast",
    components: { Modal, Typing },
    computed: {
      ...mapState('message', ['keyModalType', 'keyModal', 'activeContent']),
      key: {
        get() {
          return this.$store.state.message.activeKey
        },
        set(activeKey) {
          this.$DC('message', { activeKey })
        }
      }
    },
    methods: {
      ...mapActions('message', ['addOrUpdateMessageReplay', 'getMessageReplayList']),
      clickEdit() {
        this.$DC('message', { keyModalType: "edit" })
      },
      clickConfirmBtn() {
        console.log(this.activeContent)
        if (!this.activeContent) return this.$message('内容不能为空')
        this.$confirm("确认保存？", async () => {
          const t = await this.addOrUpdateMessageReplay()
          if (t === 'ok') {
            await this.getMessageReplayList()
            this.$message('添加成功')
            this.clickClose()
          }
        })
      },
      clickClose() {
        return this.$DC('message', { keyModal: false })
      },
      clickDel() {
        this.$confirm("确认删除？", () => {

        })
      }
    },
    mounted() {
    }
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

  .row {
    display: flex;
    margin-bottom: 28px;

    &:nth-child(3) {
      align-items: center;
    }

    .title {
      width: 80px;
      font-size: 14px;
      color: #B1B4CF;
    }
  }

  .key-input {
    background-color: #F8F9FE;
    border-radius: 12px;
    height: 44px;
    flex: 1;
    padding-left: 12px;
  }


  label {
    color: #B8BDC7;
    font-size: 14px;
    background: #F8F9FE;
  }


  .btns {
    margin: 32px 0 38px;
    display: flex;
    justify-content: center;

    .btn {
      width: 84px;
      line-height: 30px;
      border-radius: 20px;
      font-size: 16px;
      color: #fff;
    }

    .save {
      background: #396AFF;
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

    .btn {
      margin-top: 170px;
    }

    .edit {
      background-color: #5B73A0;
      margin-right: 22px;
    }

    .del {
      background-color: #EC5A75;
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

  }
</style>
