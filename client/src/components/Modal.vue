<template>
  <div v-show="show" class="modal">
    <div @click.self="close_modal" class="mask">
      <div class="main">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    props: {
      show: {
        type: Boolean,
        default: false
      },
      click_modal_close: {
        type: Boolean,
        default: true
      }
    },
    methods: {
      close_modal() {
        // if (this.click_modal_close) {
        //   this.config.show = false
        // }
      }
    }
  }
</script>

<style lang="scss" scoped>
  .modal {
    position: fixed;
    top: 0;
    right: 0;
    left: 0;
    z-index: 999;
    bottom: 0;
    cursor: default;
  }

  .mask {
    position: absolute;
    background-color: rgba(0, 0, 0, 0.5);
    top: 0;
    right: 0;
    left: 0;
    bottom: 0;
    overflow: hidden;
  }

  .main {
    border-radius: 16px;
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background-color: #fff;
    overflow: auto;
  }
</style>