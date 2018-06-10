<template>
  <div>
    <section class="section">
      <div class="container">
        <component-builds></component-builds>
      </div>
    </section>
  </div>
</template>

<script>
import Builds from './components/Builds.vue'

export default {
  name: 'app',
  components: {
    'component-builds': Builds
  }
}
</script>
