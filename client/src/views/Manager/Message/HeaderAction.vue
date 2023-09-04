<template>
  <div class="action">
    <button
      v-if="currentState==='key'"
      @click="clickAddRepay"
      class="add-replay"
    >
      添加回复
    </button>
    <button
      v-if="currentState==='broadcast'"
      @click="clickAddBroadcast"
      class="add-replay"
    >
      发布公告
    </button>
  </div>
</template>

<script>
  import { createNamespacedHelpers } from 'vuex'

  const { mapState, mapActions } = createNamespacedHelpers('message')
  export default {
    name: "HeaderAction",
    data() {
      return {}
    },
    computed: {
      ...mapState(['stateList', 'currentState'])
    },
    watch: {
      currentState(val) {
        this.$DC('message', {
          activeType: "PLAIN_TEXT",
          activeKey: "",
          activeContent: ""
        })
        switch (val) {
          case "key":
            this.getMessageReplayList()
            break
          case "broadcast":
            break
          case "add":
            this.clickAddReplay()
            break
        }
      }
    },
    methods: {
      ...mapActions(['getMessageReplayList', 'clickAddReplay']),
      clickAddRepay() {
        this.$DC('message', {
          keyModalType: "edit",
          keyModal: true,
          activeContent: "",
          activeKey: "",
          activeType: 'PLAIN_TEXT',
          activeReplayId: "",
          }
        )
      },
      clickAddBroadcast() {
        this.$DC('message', { broadcastModalType: "send", broadcastModal: true })
      }
    }
  }
</script>

<style lang="scss" scoped>
  .action {
    position: absolute;
    top: 24px;
    right: 40px;
  }

  .add-replay {
    background: #5B73A0;
    border-radius: 20px;
    height: 30px;
    padding: 0 26px;
    color: #fff;
    font-size: 16px;
  }
</style>
