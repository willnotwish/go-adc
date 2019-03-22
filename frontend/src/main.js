import Vue from 'vue'
import App from './app'
import store from './store'
import router from './router'

import AuthPlugin from './plugins/auth'
import Buefy from 'buefy'

import './websocket'
import './icons'
import 'buefy/dist/buefy.css'
import './assets/main.scss'

Vue.config.productionTip = false

Vue.use(AuthPlugin)
Vue.use(Buefy)

new Vue({
  render: h => h(App),
  store,
  router
}).$mount('#app')
