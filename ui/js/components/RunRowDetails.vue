<template>
  <article class="media">
    <div class="media-content">
      <div class="columns">
        <div class="column is-narrow">
          <!-- <div class="content"> -->
            <span class="title is-5">
                <strong>Summary</strong>
            </span>
            <table class="table" style="background-color: transparent;">
              <tbody>
                <tr>
                  <th>Count</th>
                  <td>{{ model.count }}</td>
                </tr>
                <tr>
                  <th>Total</th>
                  <td>{{ model.total }} ms</td>
                </tr>
                <tr>
                  <th>Slowest</th>
                  <td>{{ model.slowest }} ms</td>
                </tr>
                <tr>
                  <th>Fastest</th>
                  <td>{{ model.fastest }} ms</td>
                </tr>
                <tr>
                  <th>Average</th>
                  <td>{{ model.average }} ms</td>
                </tr>
                <tr>
                  <th>Requests / sec</th>
                  <td>{{ Number.parseFloat(model.rps).toFixed(2) }}</td>
                </tr>
              </tbody>
            </table>
          <!-- </div> -->
        </div>
      </div>
    </div>
  </article>
</template>


<script>
import axios from 'axios'

export default {
  data() {
    return {
      model: {}
    }
  },
  props: {
    projectId: [String, Number],
    testId: [String, Number],
    runId: [String, Number]
  },
  mounted() {
    this.loadData()
  },
  methods: {
    async loadData() {
      try {
        const { data } = await axios.get(
          `http://localhost:3000/api/projects/${this.projectId}/tests/${this.testId}/runs/${this.runId}`
        )

        this.model = data
      } catch (e) {
        this.model = {}

        this.$snackbar.open({
          message: e.message,
          type: 'is-danger',
          position: 'is-top'
        })
      }
    }
  }
}
</script>
