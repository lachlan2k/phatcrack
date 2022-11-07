/*
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
*/

import './assets/styles/font.css'
import './assets/styles/icons.css'
import '@mdi/font/css/materialdesignicons.css'

import Vue from 'vue'
import store from './store'
import App from './App.vue'
import router from './router'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import {axios} from '@bundled-es-modules/axios'
import VueAxios from 'vue-axios'
import VueMoment from 'vue-moment'
import moment from 'moment'
import mdfSetup from 'moment-duration-format'
import VueVisible from 'vue-visible'
import UserPlugin from '@/components/user/userPlugin.js'
import paginator from '@/plugins/paginator'
import VueHighlightJS from 'vue-highlightjs'

import { Chart, registerables } from 'chart.js'
import 'chartjs-adapter-moment'
Chart.register(...registerables)

document.title = window.projectName || 'Fitcrack'

mdfSetup(moment)

Vue.use(VueHighlightJS)
Vue.use(VueVisible);
Vue.use(VueMoment, {
  moment
});
Vue.use(VueAxios, axios);
Vue.axios.defaults.withCredentials = true
Vue.use(Vuetify);

Vue.use(UserPlugin);
Vue.use(paginator)

new Vue({
  el: '#app',
  router,
  store,
  vuetify: new Vuetify({
    icons: {
      iconfont: 'mdi'
    },
    theme: {
      dark: false,
      themes: window.theme
    }
  }),
  render: (h) => h(App)
}).$mount('#app');
