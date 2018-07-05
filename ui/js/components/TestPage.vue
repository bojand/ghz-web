<template>
    <component-test-details :project-id="projectId" :test-id="testId"></component-test-details>
</template>

<script>
import TestDetails from './TestDetails.vue'
import RunList from './RunList.vue'

export default {
  data() {
    return {
      projectId: 0,
      testId: 0
    }
  },
  created() {
    this.projectId = this.$route.params.projectId
    this.testId = this.$route.params.testId
  },
  async beforeRouteUpdate (to, from, next) {
    console.log('TestPage beforeRouteUpdate')
    this.projectId = to.params.projectId
    this.testId = to.params.testId
    next()
  },
  components: {
    'component-test-details': TestDetails,
    'component-run-list': RunList
  }
}
</script>
