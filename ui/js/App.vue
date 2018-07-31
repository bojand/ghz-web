<template>
    <section>
      <section class="hero">
        <div class="hero-body">
          <div class="container">
            <!-- <navbar></navbar> -->
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
import TestPage from './components/TestPage.vue'
import RunPage from './components/RunPage.vue'

import Navbar from './layout/Navbar.vue'
import VFooter from './layout/Footer.vue'
import Bread from './layout/Breadcrumb.vue'

import store from './store'

const routes = [
  { path: '/', redirect: '/projects' },
  {
    name: 'projects',
    path: '/projects',
    component: ProjectListPage
  },
  {
    name: 'project',
    path: '/projects/:projectId',
    component: ProjectPage,
    beforeEnter: async (to, from, next) => {
      if (!store.project || store.project.id !== to.params.projectId) {
        try {
          await store.fetchProject(to.params.projectId)
        } catch (e) {
          console.error(e)
        }
      }
      next()
    }
  },
  {
    name: 'test',
    path: '/projects/:projectId/tests/:testId',
    component: TestPage,
    beforeEnter: async (to, from, next) => {
      if (!store.project || store.project.id !== to.params.projectId) {
        try {
          await store.fetchProject(to.params.projectId)
        } catch (e) {
          console.error(e)
        }
      }
      if (!store.test || store.test.id !== to.params.testId) {
        try {
          await store.fetchTest(to.params.projectId, to.params.testId)
        } catch (e) {
          console.error(e)
        }
      }
      next()
    }
  },
  {
    name: 'run',
    path: '/projects/:projectId/tests/:testId/runs/:runId',
    component: RunPage,
    beforeEnter: async (to, from, next) => {
      if (!store.project || store.project.id !== to.params.projectId) {
        try {
          await store.fetchProject(to.params.projectId)
        } catch (e) {
          console.error(e)
        }
      }
      if (!store.test || store.test.id !== to.params.testId) {
        try {
          await store.fetchTest(to.params.projectId, to.params.testId)
        } catch (e) {
          console.error(e)
        }
      }
      if (!store.run || store.run.id !== to.params.runId) {
        try {
          await store.fetchRun(to.params.projectId, to.params.testId, to.params.runId)
        } catch (e) {
          console.error(e)
        }
      }
      next()
    }
  }
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
