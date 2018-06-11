<template>
  <section>
    <div class="box">
      <article class="media">
        <div class="media-left">
          <div class="media-content">
            <div class="content" v-if="!editMode">
                <p>
                <strong>{{ name }}</strong>
                <br>
                {{ description }}
                </p>
            </div>
            <div class="content" v-if="editMode">
              <b-field>
                <b-input :placeholder="name"></b-input>
              </b-field>
              <b-field>
                <b-input :placeholder="description"></b-input>
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
      name: '',
      description: '',
      loading: false,
      editMode: false
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
        const { data } = await axios.get(`http://localhost:3000/api/projects/${this.id}`)

        this.name = data.name
        this.description = data.description

        this.loading = false
      } catch (error) {
        this.loading = false
        throw error
      }
    },
    
    editClicked() {
      this.editMode = !this.editMode
    }
  },
  mounted() {
    this.loadAsyncData()
  }
}
</script>
