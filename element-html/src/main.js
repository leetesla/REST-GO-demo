import Vue from 'vue'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import App from './App.vue'
import Axios from 'axios'
import router from './components/router.js'

Vue.use(ElementUI)
Vue.prototype.$http = Axios
new Vue({
  el: '#app',
  router,
  render: h => h(App)
})
