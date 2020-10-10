import ECharts from "vue-echarts"

export default function ({ name, total, today, data, options = {}, formatter = null }) {
  let legend = [], series = [], color = []
  data.forEach(({ name, list, color: _color }) => {
    color.push(_color)
    legend.push({ name })
    series.push(getSeries(name, list, _color))
  })
  let addX = String(total).length * 9 + 45
  return {
    color,
    tooltip: {
      trigger: 'axis',
      formatter,
      axisPointer: {
        type: 'line',
        lineStyle: {
          type: 'dashed',
          color: 'red'
        }
      }
    },
    legend: {
      data: legend,
      icon: 'circle',
      type: '',
      itemWidth: 8,
      itemHeight: 8,
      itemGap: 10,
      bottom: 60,
      align: 'auto',
      textStyle: {
        fontSize: 12,
        color: '#A5A7C8',
        padding: [3, 10, 0, 0]
      },
      ...options
    },
    textStyle: {
      color: '#A5A7C8',
      fontFamily: 'Nunito',
      fontSize: 10,
    },
    title: {
      text: name,
      left: 30,
      top: 30,
      textStyle: {
        color: '#4C4471',
        fontWeight: '400',
        fontSize: 18
      },
      subtext: total,
      subtextStyle: {
        color: '#4C4471',
        fontSize: 14,
        lineHeight: 20
      }
    },
    graphic: {
      elements: [
        {
          type: 'text',
          style: {
            text: '+' + today,
            x: addX,
            y: 66,
            font: '12px',
            fill: '#1BACC0'
          }
        }
      ]
    },
    grid: {
      top: '32%',
      height: 100,
      left: '12%',
      right: '8%'
    },
    xAxis: {
      type: 'time',
      axisLabel: {
        fontSize: 10,
        color: '#A5A7C8',
        formatter: value => {
          let day = new Date(value).getDate()
          return day < 10 ? '0' + day : day
        }
      },
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      splitLine: {
        show: false
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        fontSize: 10,
        color: '#A5A7C8'
      },
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      splitLine: {
        lineStyle: {
          color: '#ECEFFF'
        }
      },
      splitNumber: 3,
    },
    series,
    animationDuration: 500,
  }
}


function getSeries(name, list, color) {
  return {
    data: list,
    name: name,
    type: 'line',
    smooth: true,
    showSymbol: false,
    symbol: 'circle',
    symbolSize: 1,
    zlevel: 2,
    hoverAnimation: true,
    itemStyle: {
      borderColor: '#fff',
      borderWidth: 2,
    },
    lineStyle: {
      color
    },
    areaStyle: {
      color: new ECharts.graphic.LinearGradient(0, 0, 0, 1, [
        { offset: 0, color },
        { offset: 1, color: '#fff' }
      ])
    }
  }
}
