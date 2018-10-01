<template>
  <section>
    <h2 class="subtitle"><strong>Create Test</strong></h2>
    <b-field grouped>
    <b-field>
        <b-input v-model="testName" placeholder="Name"></b-input>
    </b-field>
    <b-field>
        <b-input v-model="testDesc" placeholder="Description"></b-input>
    </b-field>
    <b-field>
        <p class="control">
        <button class="button is-primary" @click="createProject">Create</button>
        </p>
    </b-field>
    </b-field>
  </section>
</template>

<script>
import axios from 'axios'

export default {
  data() {
    return {
      testName: '',
      testDesc: ''
    }
  },
  props: {
    projectId: [String, Number]
  },
  methods: {
    async createProject() {
      let name = this.testName
      let description = this.testDesc

      try {
        const { data } = await axios.post(
          `http://localhost:3000/api/projects/${this.projectId}/tests`,
          {
            name,
            description
          }
        )

        name = data.name
        description = data.description

        this.$emit('test-created', {
          id: data.id,
          name,
          description
        })
        this.testName = ''
        this.testDesc = ''
      } catch (e) {
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
