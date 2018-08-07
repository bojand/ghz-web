<template>
  <section>

    <div class="content">
      <div class="js-line-container card--chart"></div>
    </div>

  </section>
</template>

<script>
import _ from 'lodash'

import common from './common.js'

const britecharts = require('britecharts')
const d3 = require('d3-selection')

export default {
  props: {
    runs: Object
  },
  watch: {
    runs(newVal, oldVal) {
      this.createLineChart()
    }
  },
  mounted() {
    this.createLineChart()
  },
  mixins: [common],
  methods: {
    createChartData() {
      let data = this.runs.data

      // data = _.map(data, r => {
      //   const date = new Date(r.date).toDateString()
      //   const v = _.clone(r)
      //   v.date = date
      //   v.fullDate = new Date(r.date).toISOString()
      //   return v
      // })

      const averages = _.map(data, r => {
        return {
          date: r.date,
          value: r.average / 1000000
        }
      })

      const slowest = _.map(data, r => {
        return {
          date: r.date,
          value: r.slowest / 1000000
        }
      })

      const fastest = _.map(data, r => {
        return {
          date: r.date,
          value: r.fastest / 1000000
        }
      })

      const nine5 = _(data)
        .map(r => {
          const elem = _.find(r.latencyDistribution, ['percentage', 95])
          if (elem) {
            return {
              date: r.date,
              value: elem.latency / 1000000
            }
          }
        })
        .compact()
        .valueOf()

      const chartData = {
        dataByTopic: [
          {
            topicName: 'Average',
            topic: 1,
            dates: averages
          },
          {
            topicName: 'Fastest',
            topic: 2,
            dates: fastest
          },
          {
            topicName: 'Slowest',
            topic: 3,
            dates: slowest
          },
          {
            topicName: '95th',
            topic: 4,
            dates: nine5
          }
        ]
      }

      console.log(chartData)

      return chartData

      // return {
      //   dataByTopic: [
      //     {
      //       topicName: 'Average',
      //       topic: 1,
      //       dates: [
      //         {
      //           date: '2018-01-01',
      //           value: 6.83
      //         },
      //         {
      //           date: '2018-01-02',
      //           value: 5.83
      //         },
      //         {
      //           date: '2018-01-03',
      //           value: 4.83
      //         },
      //         {
      //           date: '2018-01-04',
      //           value: 5.25
      //         },
      //         {
      //           date: '2018-01-05',
      //           value: 7.25
      //         },
      //         {
      //           date: '2018-01-06',
      //           value: 6.25
      //         },
      //         {
      //           date: '2018-01-07',
      //           value: 5.25
      //         },
      //         {
      //           date: '2018-01-08',
      //           value: 6.25
      //         },
      //         {
      //           date: '2018-01-09',
      //           value: 7.65
      //         }
      //       ]
      //     },
      //     {
      //       topicName: '95th',
      //       topic: 2,
      //       dates: [
      //         {
      //           date: '2018-01-01',
      //           value: 13.26
      //         },
      //         {
      //           date: '2018-01-02',
      //           value: 12.86
      //         },
      //         {
      //           date: '2018-01-03',
      //           value: 12.26
      //         },
      //         {
      //           date: '2018-01-04',
      //           value: 11.95
      //         },
      //         {
      //           date: '2018-01-05',
      //           value: 10.25
      //         },
      //         {
      //           date: '2018-01-06',
      //           value: 11.25
      //         },
      //         {
      //           date: '2018-01-07',
      //           value: 13.45
      //         },
      //         {
      //           date: '2018-01-08',
      //           value: 12.34
      //         },
      //         {
      //           date: '2018-01-09',
      //           value: 15.67
      //         }
      //       ]
      //     }
      //   ]
      // }
    },
    createLineChart() {
      let lineChart = britecharts.line()
      let tooltip = britecharts.tooltip()
      let container = d3.select('.js-line-container')
      let containerWidth = container.node() ? container.node().getBoundingClientRect().width : false
      let tooltipContainer
      let dataset

      if (containerWidth) {
        container.html('')
        if (!this.runs) {
          container.html(lineChart.loadingState())
          return
        }

        dataset = this.createChartData()

        const lineMargin = {
          // top: 60,
          bottom: 50
          // left: 50,
          // right: 30
        }

        lineChart
          .isAnimated(true)
          // .aspectRatio(0.5)
          .grid('full')
          // .tooltipThreshold(600)
          .width(containerWidth)
          .margin(lineMargin)
          .dateLabel('date')
          .on('customMouseOver', tooltip.show)
          .on('customMouseMove', tooltip.update)
          .on('customMouseOut', tooltip.hide)

        container.datum(dataset).call(lineChart)

        tooltip
          // In order to change the date range on the tooltip title, uncomment this line
          // .dateFormat(tooltip.axisTimeCombinations.DAY_MONTH)
          .title('Data')
          .shouldShowDateInTitle(true)
          .valueFormatter(value => this.formatFloat(value) + ' ms')
          .topicsOrder(
            dataset.dataByTopic.map(function(topic) {
              return topic.topic
            })
          )

        tooltipContainer = d3.select('.js-line-container .metadata-group .hover-marker')
        tooltipContainer.datum([]).call(tooltip)
      }
    }
  }
}
</script>
