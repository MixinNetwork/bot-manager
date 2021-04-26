<template>
  <div class="--modal">
    <div @click="closeModal" class="modal-mask"></div>
    <div :class="['modal_main ', show ? 'modal-show' : '']">
      <slot></slot>
    </div>
  </div>
</template>

<script>
export default {
  name: "bottomModal",
  props: {
    close: { type: Number, default: 0 }
  },
  data() {
    return {
      show: false
    };
  },
  watch: {
    close() {
      this.closeModal();
    }
  },
  methods: {
    closeModal() {
      this.show = false;
      setTimeout(() => {
        this.$emit("close_modal");
      }, 300);
    }
  },
  mounted() {
    setTimeout(() => {
      this.show = true;
    }, 50);
  }
};
</script>

<style lang="scss" scoped>
.--modal {
  ol {
    list-style: none;
  }
  .modal-mask {
    position: fixed;
    top: 0;
    right: 0;
    left: 0;
    bottom: 0;
    background-color: rgba($color: #000000, $alpha: 0.4);
    z-index: 100;
  }
  .modal_main {
    z-index: 1;
    position: fixed;
    top: 100%;
    right: 0;
    left: 0;
    bottom: 0;
    overflow: hidden;
    background-color: #fff;
    z-index: 101;
    border-radius: 1rem 1rem 0 0;
    transition: top 0.3s;
    color: #7c7c7e;
  }
}
</style>