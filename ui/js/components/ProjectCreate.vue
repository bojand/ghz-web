<template>
  <section>
    <h2 class="subtitle"><strong>Create Project</strong></h2>
    <b-field grouped>
    <b-field>
        <b-input v-model="projectName" placeholder="Name"></b-input>
    </b-field>
    <b-field>
        <b-input v-model="projectDesc" placeholder="Description"></b-input>
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
      projectName: '',
      projectDesc: ''
    }
  },
  methods: {
    async createProject() {
      let name = this.projectName
      let description = this.projectDesc

      try {
        const { data } = await axios.post(
          'http://localhost:3000/api/projects',
          {
            name,
            description
          }
        )

        name = data.name
        description = data.description

        this.$emit('project-created', {
          id: 123,
          name,
          description
        })
        this.projectName = ''
        this.projectDesc = ''
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
