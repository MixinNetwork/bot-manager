<template>
  <div class="added">
    <Typing />
    <button @click="clickSave" class="save">保存</button>
  </div>
</template>

<script>
  import Typing from "@/components/Typing"
  import { createNamespacedHelpers } from 'vuex'

  const { mapActions } = createNamespacedHelpers('message')

  export default {
    components: { Typing },
    methods: {
      ...mapActions(['addOrUpdateMessageReplay']),
      async clickSave() {
        this.$confirm('确认保存？', async () => {
          this.$DC('message', { activeKey: "hi 你好" })
          const t = await this.addOrUpdateMessageReplay()
          if (t === 'ok') {
            this.$message('保存成功')
          }
        })
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
