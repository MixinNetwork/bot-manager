<template>
  <div id="app">
    <Loading v-if="isLoading" />
    <transition :name="transitionName">
      <router-view class="page" />
    </transition>
  </div>
</template>

<script>
  import Loading from "@/components/Loading"
  import { mapState, mapActions } from "vuex"
  import tools from "@/assets/js/tools"

  export default {
    name: "app",
    components: {
      Loading
    },
    computed: {
      ...mapState("system", ["isLoading", "themeColor", "transitionName"])
    },
    methods: {
      ...mapActions(["initPage"])
    },
    watch: {
      themeColor(val) {
        tools.changeTheme(val)
      }
    },
    mounted() {
      window._vm = this
      this.initPage()
    }
  }
</script>

<style lang="scss" scoped>
  .page {
    width: 100%;
    height: 100vh;
    overflow: hidden;
  }

  .slide-left-enter-active {
    transition: transform 0.3s ease;
  }

  .slide-left-enter {
    transform: translateX(100%);
  }

  .slide-right-enter-active {
    transition: transform 0.3s ease;
  }

  .slide-right-enter {
    transform: translateX(-100%);
  }
</style>
