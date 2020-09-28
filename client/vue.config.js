module.exports = {
  productionSourceMap: false,
  devServer: {
    port: 8082
  },
  transpileDependencies: [
    'vue-echarts',
    'resize-detector'
  ],
  css: {
    loaderOptions: {
      sass: {
        prependData: `
          @import "@/assets/scss/variable.scss";
        `
      }
    }
  }
}
