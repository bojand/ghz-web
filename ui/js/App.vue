<template>
  <div>
    <section class="section">
      <div class="container">
        <component-projects></component-projects>
      </div>
    </section>
  </div>
</template>

<script>
import Projects from './components/Projects.vue'

export default {
  name: 'app',
  components: {
    'component-projects': Projects
  }
}
</script>
