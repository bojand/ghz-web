<template>
    <article class="media">
        <div class="media-content">
          <strong>{{ model.name }}</strong>
          <component-status-tags :thresholds="model.thresholds" v-if="model.thresholds"></component-status-tags>
          <div v-if="model.description">
            <p>
              {{ model.description }}
            </p>
          </div>
        </div>
    </article>
</template>

<script>
import axios from 'axios'

import StatusTags from './StatusTags.vue'

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
      } catch (e) {
        this.model = {}

        this.$snackbar.open({
          message: e.message,
          type: 'is-danger',
          position: 'is-top'
        })
      }
    }
  },
  components: {
    'component-status-tags': StatusTags
  }
}
</script>
