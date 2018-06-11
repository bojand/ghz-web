import Buefy from 'buefy'
import 'buefy/lib/buefy.css'
import 'bulma/css/bulma.css'
import Vue from 'vue'
import VueRouter from 'vue-router'
import App from './App.vue'

Vue.use(Buefy)
Vue.use(VueRouter)

new Vue({
  el: '#app',
  render: h => h(App)
})