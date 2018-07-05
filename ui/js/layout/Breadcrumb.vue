<template>
  <nav class="breadcrumb" aria-label="breadcrumbs">
    <ul v-if="parts.length">
        <router-link v-for="(part, index) in parts" :key="index" tag="li" 
            :to="{ name: part.link.name, params: part.link.params }">
            <a>{{part.label}}</a>
        </router-link>
    </ul>
  </nav>
</template>

<script>
export default {
  data() {
    return {
      parts: [
        {
          label: 'Projects',
          link: {
            name: 'projects'
          }
        },
        {
          label: 'testproject',
          link: {
            name: 'project',
            params: { id: 14 }
          }
        }
      ]
    }
  },
  mounted() {
    console.log('mounted')
    console.log(this)
    // this.$route.matched.forEach(m => {
      
    // })
  },
  watch: {
    $route(newVal, oldVal) {
      console.log('$route changed')
      console.log(newVal)
      console.log(oldVal)
    }
  },
}
</script>
