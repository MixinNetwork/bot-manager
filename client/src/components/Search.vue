<template>


  <div
    v-if="['message', 'user'].includes(routerName)"
    :class="['search', isBlurAndNoContent ? '' : 'search-active']"
    @click="clickSearch"
  >
    <input
      @focus="focusInput(true)"
      @blur="focusInput(false)"
      ref="searchInput"
      type="text"
      v-model="searchContent"
    />
    <i v-if="isBlurAndNoContent" class="iconfont icon-search">&#xe62a;</i>
    <span v-if="isBlurAndNoContent" class="icon-holder">Search</span>
    <i v-if="searchContent" @click="searchContent=''" class="iconfont icon-close">&#xe636;</i>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        isBlurAndNoContent: true,
      }
    },
    computed: {
      routerName() {
        return this.$route.name.toLowerCase()
      },
      searchContent: {
        get() {
          return this.$store.state[this.routerName].searchKey
        },
        set(value) {
          this.$store.dispatch(this.routerName + '/searchKey', value)
        }
      }
    },
    methods: {
      clickSearch() {
        this.$refs.searchInput.focus()
      },
      focusInput(state) {
        this.isBlurAndNoContent = !state && !this.searchContent
      }
    }
  }
</script>

<style lang="scss" scoped>

  .search {
    position: fixed;
    top: 20px;
    right: 40px;
    display: flex;
    align-items: center;
  }

  i {
    position: absolute;
  }

  input {
    height: 36px;
    width: 164px;
    box-shadow: 0 4px 12px rgba(0, 18, 179, 0.04);
    border-radius: 20px;
    padding-left: 16px;
  }

  .icon-search {
    color: #a5a7c8;
    left: 16px;
    cursor: text;
  }

  .icon-holder {
    position: absolute;
    color: #CED1EA;
    left: 40px;
    cursor: text;
    transform: translateY(1px);
  }

  .icon-close {
    color: #EAEAEA;
    right: 16px;
    cursor: pointer;
  }


  @media screen and (max-width: $adaptWidth) {

    .search {
      position: fixed;
      left: 20px;
      right: 20px;
      top: 0;
      z-index: 1;
    }

    input {
      padding: 0;
      width: 0;
    }

    .icon-holder {
      display: none;
    }

    .icon-search {
      position: absolute;
      left: initial;
      top: 12px;
      right: 80px;
      color: #4C4471;
    }
  }
</style>
