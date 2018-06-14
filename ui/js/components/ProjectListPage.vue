<template>
  <section>
    <component-project-create v-on:project-created="projectCreated"></component-project-create>
    <hr />
    <component-project-list ref="projectList"></component-project-list>
  </section>
</template>

<script>
import axios from 'axios'
import ProjectCreate from './ProjectCreate.vue'
import ProjectList from './ProjectList.vue'

export default {
  data() {
    return {
    }
  },
  methods: {
    projectCreated(p) {
      this.$refs.projectList.loadData()
    }
  },
  components: {
    'component-project-create': ProjectCreate,
    'component-project-list': ProjectList
  }
}
</script>
