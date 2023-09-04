<template>
  <ul class="tabs">
    <li v-for="(value,key,idx) of stateList"
        :key="idx"
        :class="key === currentState ? 'active' : ''"
        @click="toggleState(key)"
    >
      <span>{{value}}</span>
      <img :src="getTabImage(value, idx)" alt="">
    </li>
  </ul>
</template>

<script>
  export default {
    name: "Tabs",
    props: {
      stateList: {
        type: Object
      },
      currentState: {
        type: String
      },
      toggleState: {
        type: Function
      }
    },
    methods: {
      getTabImage(value, idx) {
        return require(`../assets/img/tab-${idx === 0 ? 'first' : 'other'}-${value.length}.png`)
      }
    }
  }
</script>

<style lang="scss" scoped>
  .tabs {
    display: flex;
    margin-left: 20px;

    li {
      height: 64px;
      text-align: center;
      cursor: pointer;
      color: #B1B4CF;
      padding: 30px 16px 0 16px;

      img {
        display: none;
      }

      &.active {
        background-repeat: no-repeat;
        background-size: 100% 100%;
        position: relative;
        z-index: 1;
        color: #505887;

        img {
          display: block;
          position: absolute;
          left: -20px;
          right: 0;
          top: 0;
          bottom: 0;
          z-index: -1;
          pointer-events: none;
        }
      }
    }
  }
</style>
