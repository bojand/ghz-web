<template>
  <section>
    <h2 class="subtitle strong"><strong>Project Details</strong></h2>
    <div class="box">
      
      <div class="content" v-if="!editMode">
        <span class="title is-5"><strong>{{ model.name }}</strong></span>
        <p>
        {{ model.description }}
        </p>
      </div>

      <div class="content" v-if="editMode">
        <div class="media">
          <div class="media-left">
            <b-field>
              <b-input placeholder="name" v-model="model.name" required></b-input>
            </b-field>
            <b-field>
              <b-input placeholder="description" v-model="model.description"></b-input>
            </b-field>
          </div>
        </div>
      </div>
      
      <nav class="level">
        <div class="level-left">
          <a class="level-item" aria-label="reply">
            <button :class="['button', editMode ? 'is-primary' : '']" @click="editClicked">
              <b-icon :icon="editMode ? 'check' : 'pencil'" size="is-small"></b-icon>
              <span>{{ editMode ? 'Save' : 'Edit' }}</span>
            </button>
          </a>
          <a v-if="editMode" class="level-item" aria-label="reply">
            <button class="button" @click="cancelClicked">
              <b-icon icon="cancel" size="is-small"></b-icon>
              <span>Cancel</span>
            </button>
          </a>
        </div>
      </nav>

    </div>
  </section>
</template>

<script>
import axios from 'axios'

export default {
  data() {
    return {
      id: null,
      loading: false,
      editMode: false,
      model: {
        name: '',
        description: ''
      }
    }
  },
  props: {
    projectId: [String, Number]
  },
  watch: {
    projectId(newVal, oldVal) {
      this.loadData()
    }
  },
  store: [ 'project' ],
  methods: {
    async loadData() {
      this.loading = true
      try {
        if (!this.project) {
          this.project = await this.$store.fetchProject(this.projectId)
        }

        Object.assign(this.model, this.project)

        this.loading = false
      } catch (e) {
        this.loading = false

        this.$snackbar.open({
          message: e.message,
          type: 'is-danger',
          position: 'is-top'
        })
      }
    },

    async editClicked() {
      if (this.editMode) {
        this.loading = true

        try {
          this.project = await this.$store.updateProject(this.model)
          Object.assign(this.model, this.project)

          this.loading = false
        } catch (e) {
          this.loading = false

          this.$snackbar.open({
            message: e.message,
            type: 'is-danger',
            position: 'is-top'
          })

          this.loadData()
        }
      }
      this.editMode = !this.editMode
    },

    async cancelClicked() {
      await this.loadData()
      this.editMode = false
    }
  },
  mounted() {
    this.loadData()
  }
}
</script>
