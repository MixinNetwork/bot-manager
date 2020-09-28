<template>
  <div class="messages">
    <Tabs :current-state="currentState" :state-list="stateList" :toggle-state="toggleState"/>
    <div class="main-container">
      <MessageList v-if="currentState==='message'"/>
      <BroadcastList v-else-if="currentState==='broadcast'"/>
      <KeyList v-else-if="currentState==='key'"/>
      <Add v-else-if="currentState==='add'"/>
    </div>
  </div>
</template>

<script>
  import { createNamespacedHelpers } from 'vuex'
  import Tabs from '@/components/Tabs'
  import MessageList from './MessageList'
  import BroadcastList from './BroadcastList'
  import KeyList from './KeyList'
  import Add from './Add'

  const { mapState, mapActions } = createNamespacedHelpers('message')

  export default {
    name: "Message",
    components: {
      BroadcastList, MessageList, KeyList, Add,
      Tabs,
    },
    computed: {
      ...mapState(['stateList', 'currentState'])
    },
    methods: {
      ...mapActions(['toggleState']),
      clickBroadcast() {
        this.$DC('message', { broadcastModal: true })
      },
    }
  }
</script>

<style lang="scss" scoped>
  .messages {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .head {
    margin-left: 20px;

    button {
      width: 180px;
      height: 82px;
      background: #FFFFFF;
      box-shadow: 0 8px 20px rgba(0, 18, 179, 0.08);
      border-radius: 8px;
      box-sizing: border-box;
      cursor: pointer;

      display: flex;
      align-items: center;

      i {
        display: inline-block;
        font-size: 22px;
        width: 54px;
        height: 54px;
        line-height: 54px;
        background-color: #F2EBFC;
        border-radius: 50%;
        margin: 0 14px;
        color: #773ACE;
      }

      span {
        font-size: 16px;
        color: #4c3371;
      }
    }
  }


  .main-container {
    @include main-container()
  }


  @media screen and (max-width: $adaptWidth) {

    .messages {
      height: initial;
      width: 100%;
      position: absolute;
    }

    .head {
      margin-top: 8px;
    }

    .main-container {
      margin: 20px 20px 0 20px;
      width: calc(100% - 40px);
      border-radius: 8px 8px 0 0;
    }
  }
</style>