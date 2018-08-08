<template>
  <section>

    <div class="content">
      <div>
        <canvas id="line-chart"></canvas>
      </div>
    </div>

  </section>
</template>

<script>
import _ from 'lodash'

import common from './common.js'

import Chart from 'chart.js'

const { colors } = require('../colors.js')

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

      return {
        averate: avgs,
        fastest: fasts,
        slowest: slows,
        nine5: nine5,
        dates
      }
    },
    createLineChart() {
      if (!this.runs) {
        return
      }

      const chartData = this.createChartData()

      const dates = chartData.dates

      const avgData = []
      const fastData = []
      const slowData = []
      const n5Data = []

      dates.forEach((v, i) => {
        const d = new Date(v)
        
        avgData[i] = {
          x: d,
          y: chartData.averate[i]
        }

        fastData[i] = {
          x: d,
          y: chartData.fastest[i]
        }

        slowData[i] = {
          x: d,
          y: chartData.slowest[i]
        }

        n5Data[i] = {
          x: d,
          y: chartData.nine5[i]
        }
      })

      const datasets = [
        {
          label: 'Average',
          backgroundColor: colors.blue,
          borderColor: colors.blue,
          fill: false,
          // data: chartData.averate
          data: avgData
        },
        {
          label: 'Fastest',
          backgroundColor: colors.green,
          borderColor: colors.green,
          fill: false,
          data: fastData
        },
        {
          label: 'Slowest',
          backgroundColor: colors.red,
          borderColor: colors.red,
          fill: false,
          data: slowData
        },
        {
          label: '95th',
          backgroundColor: colors.orange,
          borderColor: colors.orange,
          fill: false,
          data: n5Data
        }
      ]

      const maxLabelLength = 10

      var config = {
        type: 'line',
        data: {
          labels: dates,
          datasets: datasets
        },
        options: {
          responsive: true,
          title: {
            display: true,
            text: 'Change Over Time'
          },
          tooltips: {
            mode: 'index',
            intersect: false,
            // callbacks: {
            //   title: function(tooltipItem, data) {
            //     console.log(tooltipItem)
            //     console.log(data)
            //     const value = tooltipItem[0].xLabel
            //     return new Date(value).toLocaleString()
            //   }
            // }
          },
          hover: {
            mode: 'nearest',
            intersect: true
          },
          scales: {
            xAxes: [
              {
                display: true,
                scaleLabel: {
                  display: true,
                  labelString: 'Date'
                },
                type: 'time',
                // ticks: {
                //   callback: function(value, index, values) {
                //     const r = new Date(value).toLocaleString()
                //     return r.length > maxLabelLength ? r.substr(0, maxLabelLength) + '...' : r
                //   }
                // }
              }
            ],
            yAxes: [
              {
                display: true,
                scaleLabel: {
                  display: true,
                  labelString: 'Latency (ms)'
                }
              }
            ]
          }
        }
      }

      const ctx = document.getElementById('line-chart').getContext('2d')
      this.lineChart = new Chart(ctx, config)
    }
  }
}
</script>
