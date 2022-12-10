import Vue from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify'
import VCalendar from 'v-calendar'
import VueGtm from '@gtm-support/vue2-gtm'
import './index.css'

// Google Analytics
Vue.use(VueGtm, {
  id: 'GTM-M677X6V',
  vueRouter: router,
})

// Use v-calendar & v-date-picker components
Vue.use(VCalendar, {
  componentPrefix: 'vc',  // Use <vc-calendar /> instead of <v-calendar />
})

Vue.config.productionTip = false

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')
