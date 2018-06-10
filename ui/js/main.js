import Vue from 'vue'
import Buefy from 'buefy'
import 'buefy/lib/buefy.css'
import 'bulma/css/bulma.css'
import App from './App.vue'

Vue.use(Buefy)

new Vue({
  el: '#app',
  render: h => h(App)
})
