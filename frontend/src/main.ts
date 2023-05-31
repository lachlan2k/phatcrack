import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { fas } from '@fortawesome/free-solid-svg-icons'
import { library } from '@fortawesome/fontawesome-svg-core'
import Toast from 'vue-toastification'

import 'vue-toastification/dist/index.css'

import App from './App.vue'
import router from './router'

import './styles.css'

const app = createApp(App)

library.add(fas)

app.component('font-awesome-icon', FontAwesomeIcon)
app.use(createPinia())
app.use(Toast, {})
app.use(router)

app.mount('#app')
