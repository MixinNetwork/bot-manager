<template>
  <div class="home">
    <Nav/>

    <keep-alive>
      <router-view class="-main-container"></router-view>
    </keep-alive>
  </div>
</template>

<script>
  import Nav from './Nav'

  export default {
    name: "Home",
    components: { Nav },
    mounted() {
      if (!this.$ls.get("token")) {
        this.$DC({ isLoading: true })
        this.$bus.toAuth()
      }
    }
  }
</script>

<style lang="scss" scoped>
  .home {
    display: flex;
    padding: 30px 12px 10px 32px;
    box-sizing: border-box;
    background-color: #F8F9FE;
    overflow: auto;
  }

  .-main-container {
    min-width: 700px;
  }

  @media screen and (max-width: $adaptWidth) {
    .home {
      padding: 0;
    }
    .-main-container {
      position: fixed;
      top: 40px;
      bottom: 60px;
      overflow-y: auto;
      overflow-x: hidden;
      min-width: initial;
    }
  }


</style>
