<template>
  <section>

    <div class="content">
      <div class="js-line-container card--chart"></div>
    </div>

  </section>
</template>

<script>
import _ from 'lodash'

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
  methods: {
    createChartData() {
      const data = this.runs.data

      const averages = _.map(data, r => {
        return {
          date: r.createdAt,
          value: r.average
        }
      })

      const slowest = _.map(data, r => {
        return {
          date: r.createdAt,
          value: r.slowest
        }
      })

      const fastest = _.map(data, r => {
        return {
          date: r.createdAt,
          value: r.fastest
        }
      })

      //   return {
      //     dataByTopic: [
      //       {
      //         topicName: 'Average',
      //         topic: 1,
      //         dates: averages
      //       },
      //       {
      //         topicName: 'Fastest',
      //         topic: 2,
      //         dates: fastest
      //       },
      //       {
      //         topicName: 'Slowest',
      //         topic: 3,
      //         dates: slowest
      //       }
      //     ]
      //   }

      return {
        dataByTopic: [
          {
            topicName: 'Average',
            topic: 1,
            dates: [
              {
                date: '2018-01-01T16:00:00-08:00',
                value: 6.83
              },
              {
                date: '2018-01-02T16:00:00-08:00',
                value: 5.83
              },
              {
                date: '2018-01-03T16:00:00-08:00',
                value: 4.83
              },
              {
                date: '2018-01-04T16:00:00-08:00',
                value: 5.25
              },
              {
                date: '2018-01-05T16:00:00-08:00',
                value: 7.25
              },
              {
                date: '2018-01-06T16:00:00-08:00',
                value: 6.25
              },
              {
                date: '2018-01-07T16:00:00-08:00',
                value: 5.25
              },
              {
                date: '2018-01-08T16:00:00-08:00',
                value: 6.25
              },
              {
                date: '2018-01-09T16:00:00-08:00',
                value: 7.65
              }
            ]
          },
          {
            topicName: '95th',
            topic: 2,
            dates: [
              {
                date: '2018-01-01T16:00:00-08:00',
                value: 13.26
              },
              {
                date: '2018-01-02T16:00:00-08:00',
                value: 12.86
              },
              {
                date: '2018-01-03T16:00:00-08:00',
                value: 12.26
              },
              {
                date: '2018-01-04T16:00:00-08:00',
                value: 11.95
              },
              {
                date: '2018-01-05T16:00:00-08:00',
                value: 10.25
              },
              {
                date: '2018-01-06T16:00:00-08:00',
                value: 11.25
              },
              {
                date: '2018-01-07T16:00:00-08:00',
                value: 13.45
              },
              {
                date: '2018-01-08T16:00:00-08:00',
                value: 12.34
              },
              {
                date: '2018-01-09T16:00:00-08:00',
                value: 15.67
              }
            ]
          }
        ]
      }
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

        lineChart
          .isAnimated(true)
          // .aspectRatio(0.5)
          .grid('full')
          // .tooltipThreshold(600)
          .width(containerWidth)
          // .margin({
          //   top: 20,
          //   bottom: 20,
          //   left: 20,
          //   right: 20
          // })
          // .colorSchema(britecharts.colors.colorSchemas.green)
          .dateLabel('date')
          .on('customMouseOver', tooltip.show)
          .on('customMouseMove', tooltip.update)
          .on('customMouseOut', tooltip.hide)

        container.datum(dataset).call(lineChart)

        tooltip
          // In order to change the date range on the tooltip title, uncomment this line
          //   .dateFormat(chartTooltip.axisTimeCombinations.DAY)
          .title('Data')
          .valueFormatter(value => value + ' ms')
          .topicsOrder(dataset.dataByTopic.map(t => t.topic))

        tooltipContainer = d3.select('.js-line-container .metadata-group .hover-marker')
        tooltipContainer.datum([]).call(tooltip)
      }
    }
  }
}
</script>
