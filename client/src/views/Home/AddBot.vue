<template>

  <Modal :show="true" :click_modal_close="false">
    <div class="content">
      <h3>添加机器人</h3>
      <ul>
        <li>
          <span>Client ID</span>
          <input v-model="form.client_id" type="text" placeholder="Client ID">
        </li>
        <li>
          <span>Session ID</span>
          <input v-model="form.session_id" type="text" placeholder="Session ID">
        </li>
        <li>
          <span>Private Key</span>
          <input v-model="form.private_key" type="text" placeholder="Private Key">
        </li>
      </ul>
      <div class="btn">
        <button @click="clickCancel" class="cancel">取消</button>
        <button @click="clickAdd" class="add">添加</button>
      </div>
    </div>
  </Modal>

</template>

<script>
  import Modal from "@/components/Modal"
  import { mapActions } from 'vuex'

  export default {
    components: { Modal },
    data: () => ({
      form: {
        client_id: "",
        session_id: "",
        private_key: ""
      }
    }),
    methods: {
      ...mapActions('user', ['getBotList']),
      clickCancel() {
        this.$DC('user', { show_add_bot: false })
      },
      async clickAdd() {
        const { client_id } = await this.APIS.addBot(this.form) || {}
        if (client_id === this.form.client_id) {
          this.$message("添加成功")
          this.$DC('user', { show_add_bot: false })
          setTimeout(() => {
            window.location.reload()
          }, 500)
        }
      }
    }
  }

</script>

<style lang="scss" scoped>
  .content {
    width: 520px;
    height: 370px;
    padding: 44px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    color: #505887;
  }

  h3 {
    font-weight: 600;
    font-size: 22px;
    line-height: 28px;
    text-align: center;
  }

  ul {

    li {
      display: flex;
      align-items: center;
      height: 44px;
      margin-bottom: 8px;

      span {
        width: 96px;
      }

      input {
        flex: 1;
        height: 44px;
        background-color: #f8f9fe;
        border-radius: 12px;
        padding: 0 22px;

        &
        ::placeholder {
          color: #C0CCE3;
        }

      }
    }
  }

  .btn {
    text-align: center;

    button {
      width: 84px;
      height: 44px;
      color: #ffffff;
      border-radius: 30px;
    }

  }

  .cancel {
    background: #5B73A0;
    margin-right: 22px;
  }

  .add {
    background-color: #396AFF;
  }
</style>