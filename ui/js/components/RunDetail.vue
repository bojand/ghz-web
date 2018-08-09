<template>
  <section v-if="run.id">
    
      <div class="columns">
        <div class="column is-narrow">
          <div class="content">
            <span class="title is-5">
              <strong>Summary</strong>
            </span>
            <table class="table">
              <tbody>
                <tr>
                  <th>Count</th>
                  <td>{{ run.count }}</td>
                  <td></td>
                </tr>
                <tr>
                  <th>Total</th>
                  <td>{{ formatNano(run.total) }} ms</td>
                  <td></td>
                </tr>
                <tr>
                  <th>Slowest</th>
                  <td>{{ formatNano(run.slowest) }} ms</td>
                  <td>
                    <b-icon
                      :icon="iconifyResult(test, run.slowest, 'slowest')"
                      :type="classifyResult(test, run.slowest, 'slowest')">
                    </b-icon>
                  </td>
                </tr>
                <tr>
                  <th>Fastest</th>
                  <td>{{ formatNano(run.fastest) }} ms</td>
                  <td>
                    <b-icon
                      :icon="iconifyResult(test, run.fastest, 'fastest')"
                      :type="classifyResult(test, run.fastest, 'fastest')">
                    </b-icon>
                  </td>
                </tr>
                <tr>
                  <th>Average</th>
                  <td>{{ formatNano(run.average) }} ms</td>
                  <td>
                    <b-icon
                      :icon="iconifyResult(test, run.average, 'mean')"
                      :type="classifyResult(test, run.average, 'mean')">
                    </b-icon>
                  </td>
                </tr>
                <tr>
                  <th>Requests / sec</th>
                  <td>{{ formatFloat(run.rps, 0) }}</td>
                  <td>
                    <b-icon
                      :icon="iconifyResult(test, run.rps, 'RPS')"
                      :type="classifyResult(test, run.rps, 'RPS')">
                    </b-icon>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <div class="column">
          <div class="content">
            <span class="title is-5">
              <strong>Options</strong>
            </span>
            <b-message>
              <pre style="background-color: transparent;">{{ run.options | pretty }}</pre>
            </b-message>
          </div>
        </div>
      </div>
    
			<div class="content">
        <span class="title is-5">
				  <strong>Historam</strong>
        </span>
				<div>
          <canvas id="histogram-chart"></canvas>
        </div>
			</div>

			<div class="content" v-if="run.latencyDistribution.length > 0">
        <span class="title is-5">
				  <strong>Latency</strong>
        </span>
        <table class="table">
					<thead>
						<tr>
							<th v-for="dist in run.latencyDistribution" :key="dist.percentage">{{ dist.percentage }} %</th>
						</tr>
					</thead>
					<tbody>
						<tr>
							<td v-for="dist in run.latencyDistribution" :key="dist.latency">{{ formatNano(dist.latency) }} ms</td>
						</tr>
					</tbody>
				</table>
	  </div>

    <div class="content" v-if="run.statusCodeDistribution">
      <div class="columns">
				<div class="column is-narrow">
          <span class="title is-5">
            <strong>Status Distribution</strong>
          </span>
          <table class="table is-hoverable">
            <thead>
              <tr>
                <th>Error</th>
                <th>Count</th>
                <th>% of Total</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(value, key) in run.statusCodeDistribution" :key="key">
                <td>{{ key }}</td>
                <td>{{ value }}</td>
                <td>{{ Number.parseFloat(value / run.count * 100).toFixed(2) }}</td>
              </tr>
            </tbody>
          </table>
				</div>
      </div>
	  </div>

    <div class="content" v-if="run.errorDistribution">
        <span class="title is-5">
				  <strong>Errors</strong>
        </span>
        <table class="table is-hoverable">
					<thead>
						<tr>
              <th>Error</th>
              <th>Count</th>
              <th>% of Total</th>
            </tr>
					</thead>
					<tbody>
						<tr v-for="(value, key) in run.errorDistribution" :key="key">
							<td>{{ key }}</td>
              <td>{{ value }}</td>
              <td>{{ Number.parseFloat(value / run.count * 100).toFixed(2) }}</td>
						</tr>
					</tbody>
				</table>
	  </div>
  </section>
</template>

<script>
import _ from 'lodash'

import Chart from 'chart.js'
import common from './common.js'

const { colors } = require('../colors.js')

export default {
  props: {
    run: Object
  },
  store: ['test'],
  watch: {
    run(newVal, oldVal) {
      this.createHistogram()
    }
  },
  mixins: [common],
  mounted() {
    this.createHistogram()
  },
  filters: {
    pretty: function(value) {
      let v = value
      if (typeof v === 'string') {
        v = JSON.parse(value)
      }
      return JSON.stringify(v, null, 2)
    }
  },
  methods: {
    createHistogram() {
      if (!this.run) {
        return
      }

      console.log(this.run)

      const categories = _.map(this.run.histogram, h => {
        return Number.parseFloat(h.mark * 1000).toFixed(2)
      })

      const series = _.map(this.run.histogram, 'count')

      const totalCount = this.run.count

      const color = Chart.helpers.color

      const barChartData = {
        labels: categories,
        datasets: [
          {
            label: 'Count',
            backgroundColor: color(colors.blue)
              .alpha(0.5)
              .rgbString(),
            borderColor: colors.blue,
            borderWidth: 1,
            data: series
          }
        ]
      }

      const barOptions = {
        elements: {
          rectangle: {
            borderWidth: 2
          }
        },
        responsive: true,
        legend: {
          display: false
        },
        tooltips: {
          callbacks: {
            title: function(tooltipItem, data) {
              const value = Number.parseInt(tooltipItem[0].xLabel)
              const percent = value / totalCount * 100
              return value + ' ' + '(' + Number.parseFloat(percent).toFixed(1) + ' %)'
            }
          }
        },
        scales: {
          xAxes: [
            {
              display: true,
              scaleLabel: {
                display: true,
                labelString: 'Count'
              }
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

      const barConfig = {
        type: 'horizontalBar',
        data: barChartData,
        options: barOptions
      }

      const ctx = document.getElementById('histogram-chart').getContext('2d')
      this.histogramChart = new Chart(ctx, barConfig)
    }
  }
}
</script>
