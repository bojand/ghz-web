<template>
  <section>
    <component-project-details :project-id="projectId"></component-project-details>
    <hr />
    <component-test-list :project-id="projectId"></component-test-list>
  </section>
</template>

<script>
import ProjectDetails from './ProjectDetails.vue'
import TestList from './TestList.vue'

export default {
  data() {
    return {
      projectId: 0
    }
  },
  created() {
    this.projectId = this.$route.params.id
  },
  async beforeRouteUpdate(to, from, next) {
    this.projectId = to.params.id
    next()
  },
  components: {
    'component-project-details': ProjectDetails,
    'component-test-list': TestList
  }
}
</script>
