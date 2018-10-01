<template>
  <section>
    <component-run-details :run="run" v-if="run"></component-run-details>
    <br />
    <div class="content">
      <span class="title is-5">
        <strong>Export</strong>
      </span>
      <br />
      <br />
      <p>
        <a class="button" :href="`http://localhost:3000/api/projects/${projectId}/tests/${testId}/runs/${runId}/export?format=json`">
          <b-icon icon="download" size="is-small"></b-icon>
          <span>JSON</span>
        </a>
        <a class="button" :href="`http://localhost:3000/api/projects/${projectId}/tests/${testId}/runs/${runId}/export?format=csv`">
          <b-icon icon="download" size="is-small"></b-icon>
          <span>CSV</span>
        </a>
      </p>
    </div>
  </section>
</template>

<script>
import RunDetail from './RunDetail.vue'

export default {
  data() {
    return {
      loading: false
    }
  },
  store: ['run'],
  created() {
    this.projectId = this.$route.params.projectId
    this.testId = this.$route.params.testId
    this.runId = this.$route.params.runId
  },
  async beforeRouteUpdate (to, from, next) {
    this.projectId = to.params.projectId
    this.testId = to.params.testId
    this.runId = to.params.runId
    next()
  },
  methods: {
    async loadData() {
      this.loading = true
      try {
        if (!this.run) {
          this.run = await this.$store.fetchTest(this.projectId, this.testId, this.runId)
        }

        this.loading = false
      } catch (e) {
        this.loading = false

        this.$snackbar.open({
          message: e.message,
          type: 'is-danger',
          position: 'is-top'
        })
      }
    }
  },
  mounted() {
    this.loadData()
  },
  components: {
    'component-run-details': RunDetail
  }
}
</script>
