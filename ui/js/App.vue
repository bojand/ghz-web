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
      <br />
      <v-footer></v-footer>
  </section>
</template>

<script>
import VueRouter from 'vue-router'
import ProjectList from './components/ProjectList.vue'
import ProjectPage from './components/ProjectPage.vue'
import Navbar from './layout/Navbar.vue'
import VFooter from './layout/Footer.vue'

const routes = [
  { path: '/', redirect: '/projects' },
  { name: 'projects', path: '/projects', component: ProjectList },
  { name: 'project', path: '/projects/:id', component: ProjectPage }
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
