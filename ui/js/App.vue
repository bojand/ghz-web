<template>
    <section>
      <section class="hero">
        <div class="hero-body">
          <div class="container">
            <navbar></navbar>
            <br />
            <bread></bread>
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

import ProjectListPage from './components/ProjectListPage.vue'
import ProjectPage from './components/ProjectPage.vue'
import Navbar from './layout/Navbar.vue'
import VFooter from './layout/Footer.vue'
import Bread from './layout/Breadcrumb.vue'

const routes = [
  { path: '/', redirect: '/projects' },
  { name: 'projects', path: '/projects', component: ProjectListPage },
  { name: 'project', path: '/projects/:id', component: ProjectPage }
]

const router = new VueRouter({
  routes
})

export default {
  name: 'app',
  components: {
    Navbar,
    VFooter,
    Bread
  },
  router
}
</script>
