<template>
  <div class="data">
    <ul class="header">
      <li v-for="(item,idx) in header" :key="idx" class="shadow">
        <i class="iconfont">{{item.icon}}</i>
        <div class="right">
          <h1>{{statistics[idx] && statistics[idx].today}}</h1>
          <i>{{item.name}}</i>
        </div>
      </li>
    </ul>

    <ul class="charts">
      <li class="chart-item" v-for="(data, idx) in statistics" :key="idx">
        <v-chart :options="getOptions(data)"></v-chart>
      </li>
    </ul>
  </div>
</template>

<script>
  import ECharts from 'vue-echarts'
  import 'echarts/lib/chart/line'
  import 'echarts/lib/component/legend'
  import 'echarts/lib/component/legendScroll'
  import 'echarts/lib/component/title'
  import 'echarts/lib/component/polar'
  import 'echarts/lib/component/tooltip'
  import 'echarts/lib/component/graphic'

  import { mapState, mapActions } from 'vuex'
  import getOptions from '@/assets/js/echarts'

  export default {
    name: "Data",
    components: {
      vChart: ECharts
    },
    data() {
      return {
        header: [
          { icon: '\ue632', name: '新用户' },
          { icon: '\ue631', name: '留言次数' },
        ]
      }
    },
    computed: {
      ...mapState('data', ['statistics'])
    },
    methods: {
      getOptions,
      ...mapActions('data', ["initState", "toggleAsset"])
    },
    mounted() {
      this.initState()
    }
  }


</script>

<style lang="scss" scoped>


  .data {
    flex: 1;
  }

  .header {
    display: flex;
    margin-left: 20px;

    li {
      display: flex;
      flex: 1;
      height: 116px;
      background-color: #fff;
      box-shadow: 0 6px 20px rgba(0, 18, 179, 0.08);
      border-radius: 8px;
      align-items: center;
      min-width: 200px;
      margin-right: 20px;

      .iconfont {
        display: inline-block;
        width: 68px;
        height: 68px;
        line-height: 68px;
        font-size: 40px;
        margin: 0 20px;
        text-align: center;
        border-radius: 50%;
        user-select: none;

      }

      &:nth-child(1) {
        .iconfont {
          color: #773ACE;
          background: #F2EBFC;
        }
      }

      &:nth-child(2) {
        .iconfont {
          color: #26B8F4;
          background: #EDF9FE;
        }
      }

      &:nth-child(3) {
        .iconfont {
          color: #FF895F;
          background: #FEF3EF;
        }
      }

      &:nth-child(4) {
        .iconfont {
          color: #F072B4;
          background: #FFEAF5;
        }
      }


    }


    h1 {
      font-family: 'DIN Condensed', serif;

      span {
        font-family: Nunito, serif;
        font-size: 1rem;
        color: #a5a7c8;
      }
    }

    i {
      color: #ced1ea;
    }
  }


  .charts {
    display: flex;
    flex-wrap: wrap;
    margin-left: 20px;

    .chart-item {
      background: #FFFFFF;
      box-shadow: 0 6px 20px rgba(0, 18, 179, 0.08);
      border-radius: 8px;
      position: relative;
      width: calc(50% - 20px);
      height: 292px;

      margin-top: 20px;

      &:nth-child(2n-1) {
        margin-right: 20px;
      }
    }

    .echarts {
      width: 100%;
      height: 325px;
    }
  }

  .menus {
    position: absolute;
    top: 22px;
    right: 30px;
    text-align: right;

    .iconfont {
      cursor: pointer;
    }

    .menu-list {
      text-align: left;
      width: 90px;
      margin-top: 10px;
      background-color: #fff;
      box-shadow: 0 6px 20px rgba(0, 18, 179, 0.08);
      border-radius: 8px;
      padding: 6px 14px;
      color: #4C4471;

      .menu-item {
        font-size: 12px;
        line-height: 20px;
        padding: 6px 0;
        position: relative;
        cursor: pointer;

        &.active::before {
          position: absolute;
          font-family: "iconfont";
          content: '\e622';
          right: 0;
        }


      }
    }
  }


  @media screen and (max-width: $adaptWidth) {
    .header {
      flex-wrap: wrap;
      margin: 0;

      li {
        flex: initial;
        width: calc(50% - 24px);
        margin: 0 0 8px 0;
        min-width: initial;

        .iconfont {
          width: 54px;
          height: 54px;
          line-height: 54px;
          font-size: 30px;
          margin-left: 12px;
        }


        &:nth-child(2n-1) {
          margin-right: 8px;
          margin-left: 20px;
        }
      }
    }

    .charts {
      margin: 0;

      .chart-item {
        margin: 0 0 20px 0;
        width: 100%;

        &:nth-child(2n-1) {
          margin: 0 0 20px 0;
        }
      }
    }


  }

</style>