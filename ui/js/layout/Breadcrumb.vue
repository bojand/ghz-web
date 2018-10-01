<template>
  <nav class="breadcrumb" aria-label="breadcrumbs">
    <ul v-if="parts.length">
        <router-link v-for="(part, index) in parts" v-if="part" :key="index" tag="li" 
            :class="index == parts.length -1 ? 'is-active' : ''" 
            :to="{ name: part.link.name, params: part.link.params }">
            <a>
              <b-icon v-if="part.icon" :icon="part.icon" size="is-small"></b-icon>
              {{part.label}}
            </a>
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
          icon: 'view-dashboard',
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

      if (this.project && this.$route.params.projectId) {
        this.parts.push(this.getProjectPart())
      }

      if (this.project && this.test && this.$route.params.testId) {
        this.parts.push(
          {
            icon: 'gauge',
            label: 'Tests',
            link: {
              name: 'project',
              params: { projectId: this.project.id }
            }
          },
          {
            label: this.test.name,
            link: {
              name: 'test',
              params: { projectId: this.project.id, testId: this.test.id }
            }
          }
        )
      }

      if (this.project && this.test && this.run && this.$route.params.runId) {
        this.parts.push(
          {
            icon: 'poll-box',
            label: 'Runs',
            link: {
              name: 'test',
              params: { projectId: this.project.id, testId: this.test.id }
            }
          },
          {
            label: this.run.id,
            link: {
              name: 'run',
              params: { projectId: this.project.id, testId: this.test.id, runId: this.run.id }
            }
          }
        )
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
    },
    test() {
      this.buildParts()
    },
    run() {
      this.buildParts()
    }
  }
}
</script>
