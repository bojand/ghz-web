<template>
  <section>

    <div class="content">
      <div id="line-chart"></div>
    </div>

  </section>
</template>

<script>
import _ from 'lodash'

import common from './common.js'

const ApexCharts = require('apexcharts')

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

      const avgs = data.map(d => d.average / 1000000)
      const fasts = data.map(d => d.fastest / 1000000)
      const slows = data.map(d => d.slowest / 1000000)

      const nine5 = _(data)
        .map(r => {
          const elem = _.find(r.latencyDistribution, ['percentage', 95])
          if (elem) {
            return elem.latency / 1000000
          }
        })
        .compact()
        .valueOf()

      const dates = data.map(d => d.date)

      const series = [
        {
          name: 'Average',
          data: avgs
        },
        {
          name: 'Fastest',
          data: fasts
        },
        {
          name: 'Slowest',
          data: slows
        },
        {
          name: '95th',
          data: nine5
        }
      ]

      return {
        series,
        dates
      }
    },
    createLineChart() {
      if (!this.runs) {
        return
      }

      const dataset = this.createChartData()

      const maxLabelLenght = 10

      var options = {
        chart: {
          height: '500',
          width: '100%',
          type: 'line',
          animations: {
            initialAnimation: {
              enabled: false
            }
          },
          zoom: {
            enabled: false
          }
        },
        dataLabels: {
          enabled: false
        },
        title: {
          text: 'Change Over Time',
          align: 'left'
        },
        grid: {
          borderColor: '#e7e7e7',
          row: {
            colors: ['#f3f3f3', 'transparent'], // takes an array which will be repeated on columns
            opacity: 0.5
          }
        },
        stroke: {
          curve: 'smooth'
        },
        series: dataset.series,
        yaxis: {
          title: {
            text: 'Latency (ms)'
          }
        },
        xaxis: {
          categories: dataset.dates,
          title: {
            text: 'Date'
          },
          labels: {
            formatter: value => {
              const r = new Date(value).toLocaleString()
              return r.length > maxLabelLenght ? r.substr(0, maxLabelLenght) + '...' : r
            }
          }
        },
        tooltip: {
          x: {
            formatter: value => new Date(value).toLocaleString()
          }
        }
      }

      var chart = new ApexCharts(document.querySelector('#line-chart'), options)

      chart.render()
    }
  }
}
</script>
