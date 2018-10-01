import Buefy from 'buefy'
import Vue from 'vue'
import VueRouter from 'vue-router'
import VueStash from 'vue-stash'

import App from './App.vue'
import store from './store'

Vue.use(Buefy)
Vue.use(VueRouter)
Vue.use(VueStash)

new Vue({
  el: '#app',
  render: h => h(App),
  data: { store }
})
