<template>
  <div class="setting">
    <Tabs :state-list="stateList" :current-state="currentState"/>
    <ul class="main-container">
      <li class="card-item">
        <i class="iconfont close-btn" @click="clickClose">&#xe636;</i>
        <img src="https://mixin-images.zeromesh.net/64i92Yf4FuaThNDZr3ujkakWpAtOQtIYTRTgPZddbCpmG7wdQzJ4gT314DzIwHVk8_JwmPub2mYjac4ncAWW9BIQ=s256" alt=""/>
        <div class="card-text-info">
          <span>2048</span>
          <span>7000101233</span>
        </div>
      </li>
      <li class="card-item">
        <i class="iconfont close-btn" @click="clickClose">&#xe636;</i>
        <img src="https://mixin-images.zeromesh.net/64i92Yf4FuaThNDZr3ujkakWpAtOQtIYTRTgPZddbCpmG7wdQzJ4gT314DzIwHVk8_JwmPub2mYjac4ncAWW9BIQ=s256" alt=""/>
        <div class="card-text-info">
          <span>2048</span>
          <span>7000101233</span>
        </div>
      </li>
      <li class="card-item-add">
        <i class="iconfont add-btn">&#xe630;</i>
        <span>添加</span>
      </li>

    </ul>
  </div>
</template>

<script>
  import Tabs from "@/components/Tabs"

  export default {
    components: { Tabs },
    data() {
      return {
        stateList: { share: "关联分享" },
        currentState: "share"
      }
    },
    methods: {
      clickClose() {
      }
    }
  }
</script>

<style lang="scss" scoped>

  .setting {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .main-container {
    @include main-container();
    padding: 32px 44px;
    display: flex;
  }

  .card-item {
    width: 223px;
    height: 75px;
    background-color: #F8F9FE;
    border-radius: 12px;

    display: flex;
    align-items: center;
    padding-left: 20px;

    position: relative;
    margin-right: 20px;
  }

  img {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    margin-right: 16px;
  }

  .card-text-info {
    display: flex;
    flex-direction: column;

    span {
      &:first-child {
        font-size: 14px;
        color: #505887;
      }

      &:last-child {
        font-size: 12px;
        color: #B1B4CF;
      }
    }
  }

  .close-btn {
    position: absolute;
    right: -10px;
    top: -10px;
    font-size: 20px;
    color: #EAEAEA;
    cursor: pointer;
  }

  .card-item-add {
    width: 223px;
    height: 75px;
    border: 1px dashed #C4C7DB;
    border-radius: 12px;
    display: flex;
    justify-content: center;
    align-items: center;
    color: #B1B4CF;
    cursor: pointer;

    .add-btn {
      font-size: 16px;
      margin-right: 8px;
    }

    span {
      font-size: 14px;
      line-height: 20px;
      letter-spacing: 2px;
    }

  }


</style>