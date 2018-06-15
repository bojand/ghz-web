<template>
  <section>
    <component-project-details :project-id="projectId"></component-project-details>
    <hr />
    <component-test-create :project-id="projectId" v-on:test-created="onTestCreated"></component-test-create>
    <br>
    <component-test-list ref="testList" :project-id="projectId"></component-test-list>
  </section>
</template>

<script>
import ProjectDetails from './ProjectDetails.vue'
import TestList from './TestList.vue'
import TestCreate from './TestCreate.vue'

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
  methods: {
    onTestCreated() {
      this.$refs.testList.loadData()
    }
  },
  components: {
    'component-project-details': ProjectDetails,
    'component-test-create': TestCreate,
    'component-test-list': TestList
  }
}
</script>
