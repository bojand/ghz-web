<template>
    <section>
      <section class="hero">
        <div class="hero-body">
          <div class="container">
            <navbar></navbar>
          </div>
        </div>
      </section>
      <section>
        <div class="container">
          <router-view></router-view>
        </div>
      </section>
      <v-footer></v-footer>
  </section>
</template>

<script>
import VueRouter from 'vue-router'
import Projects from './components/Projects.vue'
import Navbar from './layout/Navbar.vue'
import VFooter from './layout/Footer.vue'

const routes = [
  { path: '/', redirect: '/projects' },
  { path: '/projects', component: Projects }
]

const router = new VueRouter({
  routes
})

export default {
  name: 'app',
  components: {
    Navbar,
    VFooter
  },
  router
}
</script>
