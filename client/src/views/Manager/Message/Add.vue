<template>
  <div class="added">
    <Typing />
    <button class="save">保存</button>
  </div>
</template>

<script>
  import Typing from "@/components/Typing"

  export default {
    components: { Typing },
    methods: {
      clickSave() {
        this.$DC('message', { activeKey: "Hi 你好" })
      }
    }
  }
</script>

<style lang="scss" scoped>
  .added {
    padding: 114px;
  }

  .save {
    display: block;
    text-align: center;
    margin: 38px auto;
    width: 114px;
    height: 44px;
    background-color: #396AFF;
    color: #ffffff;
    border-radius: 30px;
    font-weight: 600;
    font-size: 16px;
  }
</style>
