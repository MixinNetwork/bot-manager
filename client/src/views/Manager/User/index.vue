<template>
  <div class="growth">
    <Tabs :current-state="currentState" :state-list="stateList" :toggle-state="toggleState"/>
    <ul class="main-container">
      <template v-if="currentState==='normal'">
        <li v-for="_ in 10" class="user-item">
          <img src="https://mixin-images.zeromesh.net/OsaSpGZMBV4PmQ2Om-UnDZ-878Bk37heqprakp_Sll6MWM-ciLdUQrvEDIeSF4z3t0sgfXt8Hw4zmDkiR2irag0=s256">
          <span class="user-name">Crossle</span>
          <span class="user-number">Mixin ID: 26648</span>
          <span class="user-time">加入时间：2020/02/01 12:01:11</span>
          <button>操作</button>
        </li>
      </template>
      <template v-else-if="currentState==='block'">
        <li v-for="_ in 2" class="user-item">
          <img src="https://mixin-images.zeromesh.net/OsaSpGZMBV4PmQ2Om-UnDZ-878Bk37heqprakp_Sll6MWM-ciLdUQrvEDIeSF4z3t0sgfXt8Hw4zmDkiR2irag0=s256">
          <span class="user-name">Crossle</span>
          <span class="user-number">Mixin ID: 26648</span>
          <span class="user-time">加入时间：2020/02/01 12:01:11</span>
          <span class="user-time">屏蔽时间：2020/02/01 14:01:11</span>
          <button>恢复</button>
        </li>
      </template>
    </ul>
  </div>
</template>

<script>
  import Tabs from '@/components/Tabs'
  import { mapState, mapActions } from 'vuex'

  export default {
    name: "User",
    components: { Tabs },
    computed: {
      ...mapState('growth', ['stateList', 'currentState', 'growthList'])
    },
    methods: {
      ...mapActions('growth', ['toggleState']),
      clickAction(activeGrower) {
        this.$DC('growth', { activeGrower, actionState: true })
      },
      clickState(key) {
        this.toggleState(key)
      }
    },
    mounted() {
      this.APIS.getGrowthUserList('pending')
    }
  }
</script>

<style lang="scss" scoped>

  .growth {
    display: flex;
    width: 100%;
    flex-direction: column;
    color: #4C4471;
    overflow: hidden;
    height: 100%;
  }

  .main-container {
    @include main-container();

    display: flex;
    flex-wrap: wrap;
  }

  .user-item {
    width: 223px;
    height: 226px;
    margin: 20px;
    border-radius: 12px;
    background-color: #F8F9FE;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;


    img {
      width: 48px;
      height: 48px;
      border-radius: 50%;
    }

    .user-name {
      white-space: nowrap;
      padding: 0 20px;
      width: 100%;
      overflow: hidden;
      text-overflow: ellipsis;
      text-align: center;
      margin: 2px 0 10px;
      font-weight: 500;

      font-size: 16px;
    }

    .user-number {
      font-size: 14px;
      color: #B1B4CF;
    }

    .user-time {
      font-size: 12px;
      color: #B1B4CF;
    }

    button {
      margin-top: 26px;
      height: 30px;
      width: 76px;
      color: #fff;
      background-color: #5B73A0;
      border-radius: 20px;

      &:hover {
        background: #396AFF;
      }
    }

  }

  @media screen and (max-width: $adaptWidth) {
    .growth {
      padding: 8px 20px 0 20px;
      margin: 0;
    }
    .main-container {
      padding: 0;
      border-radius: 8px 8px 0 0;
    }

  }

</style>