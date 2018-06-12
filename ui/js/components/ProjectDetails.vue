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
  created() {
    this.id = this.$route.params.id
    console.log(`created id: ${this.id}`)
  },
  async beforeRouteUpdate(to, from, next) {
    console.log('beforeRouteUpdate')
    this.id = to.params.id
    await this.loadAsyncData()
    next()
  },
  methods: {
    async loadAsyncData() {
      this.loading = true
      try {
        const { data } = await axios.get(
          `http://localhost:3000/api/projects/${this.id}`
        )

        console.log(data)
        this.model = data
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
          const { data } = await axios.put(
            `http://localhost:3000/api/projects/${this.id}`,
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

          this.loadAsyncData()
        }
      }
      this.editMode = !this.editMode
    },

    async cancelClicked() {
      await this.loadAsyncData()
      this.editMode = false
    }
  },
  mounted() {
    this.loadAsyncData()
  }
}
</script>
