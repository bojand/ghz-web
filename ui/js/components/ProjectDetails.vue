<template>
  <section>
    <h2 class="subtitle strong"><strong>Project Details</strong></h2>
    <div class="box">
      <article class="media">
        <div class="media-left">
          <div class="media-content">
            <div class="content" v-if="!editMode">
                <p>
                <strong>{{ model.name }}</strong>
                <br>
                {{ model.description }}
                </p>
            </div>
            <div class="content" v-if="editMode">
              <b-field>
                <b-input :placeholder="model.name" v-model="model.name"></b-input>
              </b-field>
              <b-field>
                <b-input :placeholder="model.description" v-model="model.description"></b-input>
              </b-field>
            </div>
            <nav class="level is-mobile">
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
        </div>
      </article>
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
  methods: {
    async loadData() {
      this.loading = true
      try {
        const { data } = await axios.get(
          `http://localhost:3000/api/projects/${this.projectId}`
        )

        this.model = data
        this.loading = false

        this.$store.project = this.model
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
          const { data } = await axios.put(
            `http://localhost:3000/api/projects/${this.projectId}`,
            this.model
          )

          this.model = data
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
  },
  beforeDestroy () {
    console.log('beforeDestroy')
    this.$store.project = null
  },
  beforeRouteLeave (to, from, next) {
    console.log('beforeRouteLeave')
    this.$store.project = null
  }
}
</script>
