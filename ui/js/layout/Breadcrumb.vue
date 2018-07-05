<template>
  <nav class="breadcrumb" aria-label="breadcrumbs">
    <ul v-if="parts.length">
        <router-link v-for="(part, index) in parts" v-if="part" :key="index" tag="li" 
            :class="index == parts.length -1 ? 'is-active' : ''" 
            :to="{ name: part.link.name, params: part.link.params }">
            <a>{{part.label}}</a>
        </router-link>
    </ul>
  </nav>
</template>

<script>
import _ from 'lodash'

export default {
  data() {
    return {
      parts: [
        {
          label: 'Projects',
          link: {
            name: 'projects'
          }
        }
      ]
    }
  },
  store: ['project', 'test', 'run'],
  mounted() {
    this.buildParts()
  },
  methods: {
    buildParts() {
      if (this.parts.length > 1) {
        this.parts.splice(1)
      }

      if (
        this.project &&
        (_.find(this.$route.matched, ['name', 'project']) || this.$route.name === 'project')
      ) {
        this.parts.push(this.getProjectPart())
      }

      if (
        this.test &&
        (_.find(this.$route.matched, ['name', 'test']) || this.$route.name === 'test')
      ) {
        this.parts.push({
          label: 'Tests',
          link: {
            name: 'tests'
          }
        })
      }
    },
    getProjectPart() {
      if (this.project) {
        return {
          label: this.project.name,
          link: {
            name: 'project',
            params: { id: this.project.id }
          }
        }
      }
      return null
    }
  },
  watch: {
    $route(newVal, oldVal) {
      this.buildParts()
    },
    project() {
      this.buildParts()
    }
  }
}
</script>
