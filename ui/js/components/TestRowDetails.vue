<template>
    <article class="media">
        <div class="media-content">
            <div class="columns">
              <div class="column">
                <strong>{{ model.name }}</strong>
                <div v-if="model.description">
                  {{ model.description }}
                </div>
              </div>
              <div class="column" v-if="model.thresholds">
                  <strong>Thresholds</strong>
                  <div v-for="(value, key) in model.thresholds" :key="key">
                    {{ key }}: {{ value }}
                  </div>
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
    testId: [String, Number]
  },
  mounted() {
    this.loadData()
  },
  methods: {
    async loadData() {
      try {
        const { data } = await axios.get(
          `http://localhost:3000/api/projects/${this.projectId}/tests/${this.testId}`
        )

        this.model = data
        data.thresholds = {
          '50th': 1000,
          '90th': 2000,
          '95th': 3000,
          '99th': 4000
        }
        console.log(JSON.stringify(data))
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
