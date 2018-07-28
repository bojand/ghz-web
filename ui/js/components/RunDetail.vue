<template>
  <section>
    
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
                </tr>
                <tr>
                  <th>Total</th>
                  <td>{{ run.total }} ms</td>
                </tr>
                <tr>
                  <th>Slowest</th>
                  <td>{{ run.slowest }} ms</td>
                </tr>
                <tr>
                  <th>Fastest</th>
                  <td>{{ run.fastest }} ms</td>
                </tr>
                <tr>
                  <th>Average</th>
                  <td>AVG ms</td>
                </tr>
                <tr>
                  <th>Requests / sec</th>
                  <td>{{ run.rps }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    
			<div class="content">
        <span class="title is-5">
				  <strong>Historam</strong>
        </span>
				<div class="js-bar-container"></div>
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
							<td v-for="dist in run.latencyDistribution" :key="dist.latency">{{ dist.latency }} ms</td>
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

const britecharts = require('britecharts')
const d3 = require('d3-selection')

export default {
  props: {
    run: Object
  },
  watch: {
    run(newVal, oldVal) {
      this.createHistogram()
    }
  },
  mounted() {
    this.createHistogram()
  },
  methods: {
    createHistogram() {
      let barChart = britecharts.bar(),
        tooltip = britecharts.miniTooltip(),
        barContainer = d3.select('.js-bar-container'),
        containerWidth = barContainer.node()
          ? barContainer.node().getBoundingClientRect().width
          : false,
        tooltipContainer,
        dataset,
        count = this.run.count

      tooltip.numberFormat('')
      tooltip.valueFormatter(function(v) {
        var percent = v / count * 100
        return v + ' ' + '(' + Number.parseFloat(percent).toFixed(1) + ' %)'
      })

      if (containerWidth) {
        dataset = _.map(this.run.histogram, h => {
          return {
            name: Number.parseFloat(h.mark * 1000).toFixed(2) + ' ms',
            value: h.count
          }
        })

        barChart
          .isHorizontal(true)
          .isAnimated(true)
          .margin({
            left: 100,
            right: 20,
            top: 20,
            bottom: 20
          })
          .colorSchema(britecharts.colors.colorSchemas.teal)
          .width(containerWidth)
          // .yAxisPaddingBetweenChart(10)
          .height(400)
          // .hasPercentage(true)
          .enableLabels(true)
          .labelsNumberFormat('')
          .percentageAxisToMaxRatio(1.3)
          .on('customMouseOver', tooltip.show)
          .on('customMouseMove', tooltip.update)
          .on('customMouseOut', tooltip.hide)

        barChart.orderingFunction(function(a, b) {
          var nA = a.name.replace(/ms/gi, '')
          var nB = b.name.replace(/ms/gi, '')

          var vA = Number.parseFloat(nA)
          var vB = Number.parseFloat(nB)

          return vB - vA
        })

        barContainer.datum(dataset).call(barChart)

        tooltipContainer = d3.select('.js-bar-container .bar-chart .metadata-group')
        tooltipContainer.datum([]).call(tooltip)
      }
    }
  }
}
</script>
