<template>
  <Modal :show="addFavoriteModal">
    <div class="modal-container">
      <h3>添加关联分享</h3>
      <i class="iconfont close-btn" @click="clickToggleModal">&#xe636;</i>
      <div class="row">
        <span class="title">机器人ID</span>
        <input v-model="key" class="key-input" type="text" placeholder="请输入机器人ID，如7000123456" />
      </div>
      <div class="btns">
        <button @click="clickConfirmBtn" class="btn send">发布</button>
      </div>
    </div>
  </Modal>
</template>

<script>
  import { createNamespacedHelpers } from 'vuex'
  import Modal from '@/components/Modal'

  const { mapState, mapActions } = createNamespacedHelpers('user')
  export default {
    name: "addFavorite",
    data() {
      return {
        key: ""
      }
    },
    components: { Modal },
    computed: {
      ...mapState(['addFavoriteModal'])
    },
    watch: {
      addFavoriteModal(val) {
        if (val) this.key = ""
      }
    },
    methods: {
      ...mapActions(['addBotFavorite', 'getBotFavorite']),
      clickToggleModal() {
        this.$DC('user', { addFavoriteModal: !this.addFavoriteModal })
      },
      clickConfirmBtn() {
        this.$confirm("确认添加？", async () => {
          const t = await this.addBotFavorite(this.key)
          if (t && t.app_id) {
            this.$message("添加成功")
            this.clickToggleModal()
            this.getBotFavorite()
          }
        })
      },
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
    align-items: center;
    margin-bottom: 28px;

    .title {
      width: 80px;
      font-size: 14px;
      color: #B1B4CF;
    }

    input {

      background-color: #F8F9FE;
      border-radius: 12px;
      height: 44px;
      flex: 1;
      padding-left: 12px;
    }
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
  }
</style>
