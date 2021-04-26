<template>
  <Modal :show="visible">
    <div class="confirm-container">
      <h3>{{title}}</h3>
      <!--        <p>确认发送给所有{{messageTypeMap[activeBroadcastType]}}用户么？</p>-->
      <div class="btns">
        <button @click="visible=false" class="cancel">{{cancel}}</button>
        <button @click="onSuccess" class="confirm">{{confirm}}</button>
      </div>
    </div>
  </Modal>
</template>
<script>
  import Modal from '../../Modal'

  export default {
    components: { Modal },
    data() {
      return {
        visible: false,
        title: "确认？",
        confirm: "确认",
        cancel: "取消",
        onSuccess() {
        }
      }
    }
  }
</script>
<style lang="scss" scoped>

  .confirm-container {
    width: 350px;
    text-align: center;
    color: #4C4471;
    padding: 74px 16px;
  }

  h3 {
    font-size: 22px;
    color: #505887;
    font-weight: 400;
    margin-bottom: 42px;
  }

  button {
    width: 84px;
    height: 30px;
    color: #ffffff;
    border-radius: 20px;
  }

  .confirm {
    background-color: #396AFF;
  }

  .cancel {
    background: #5B73A0;
    margin-right: 10px;
  }
</style>
