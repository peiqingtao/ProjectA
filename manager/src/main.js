import Vue from 'vue'
import App from './App.vue'
//elementUI
import Element from 'element-ui'
import 'element-ui/lib/locale/lang/en'
import  'element-ui/lib/theme-chalk/index.css'

Vue.use(Element)
//vue-router
import router from './router/router.js'
//css
import './assets/manger.css'

//axios vue-axios
import axios from 'axios'
import VueAxios from 'vue-axios'
Vue.use(VueAxios, axios)
Vue.prototype.$axios = axios
axios.defaults.baseURL = '/api'
// axios.defaults.headers.options['Content-Type'] = 'application/json;chaeset=utf-8';
// axios.defaults.headers.post['Content-Type'] = 'application/json;chaeset=utf-8';

axios.defaults.baseURL = '/api'   //关键代码 
Vue.config.productionTip = false

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
