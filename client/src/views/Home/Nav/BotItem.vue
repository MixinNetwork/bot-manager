<template>
  <div class="bot-item" @click.stop="clickBotItem">
    <template v-if="isAdd">
      <i class="iconfont icon-add">&#xe630;</i>
      <div class="bot-info">
        <span class="bot-name">导入机器人</span>
      </div>
    </template>
    <template v-else>
      <Avatar :user_info="botItem" />
      <div class="bot-info">
        <span class="bot-name">{{botItem.full_name}}</span>
        <span class="bot-number">{{botItem.identity_number}}</span>
      </div>
      <i v-if="isHead" class="iconfont">&#xe62f;</i>
    </template>
  </div>
</template>


<script>
  import { mapActions, mapState } from 'vuex'
  import Avatar from "../../../components/Avatar"

  export default {
    name: "BotItem",
    components: { Avatar },
    props: {
      botItem: {
        type: Object,
        default() {
          return {}
        }
      },
      isAdd: {
        type: Boolean,
        default: false
      },
      isHead: {
        type: Boolean,
        default: false
      }
    },
    computed: {
      ...mapState('user', ['active_bot'])
    },
    methods: {
      ...mapActions('user', ['changeActiveBot']),
      clickBotItem() {
        if (!this.isHead && !this.isAdd) {
          const { client_id } = this.botItem
          this.$emit('toggleList')
          if (client_id === this.active_bot.client_id) return
          this.$ls.set('bot', client_id)
          window.location.reload()
        }
        if (this.isHead) this.$emit('toggleList')
        if (this.isAdd) {
          this.$DC('user', { show_add_bot: true })
          this.$emit('toggleList')
        }
      }
    }
  }

</script>

<style lang="scss" scoped>


  .bot-item {
    display: flex;
    align-items: center;
    padding: 14px 18px;
    cursor: pointer;
    width: 100%;
  }

  .bot-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    margin-left: 14px;
  }

  .bot-number {
    font-size: 12px;
    color: #b1b4cf;
  }

  .bot-name {
    font-size: 14px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 130px;
    margin-bottom: 1px;
  }

  .icon-add {
    font-size: 16px;
    width: 32px;
    line-height: 32px;
    text-align: center;
    background-color: #ddebff;
    border-radius: 50%;
    color: #396AFF;
    margin-right: 16px;
  }

</style>
