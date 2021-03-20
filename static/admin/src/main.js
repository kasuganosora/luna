import Vue from 'vue'
import App from './App.vue'
import router from './router'
import './plugins/element.js'
import store from './store'
import moment from 'moment'
import 'moment/locale/zh-cn'
import mavonEditor from 'mavon-editor'
import 'mavon-editor/dist/css/index.css'

Vue.config.productionTip = false

const mounted = () => {
  let self = this
  window.onresize = function(){
    self.$store.state.screenWidth = document.documentElement.clientWidth;
    self.$store.state.screenHeight = document.documentElement.clientHeight;
  }
}

// 中文简体
moment.locale('zh-cn');
window.ss = moment;
Vue.filter("fromNow", function (time) {
  return moment(time).fromNow()
})

Vue.use(mavonEditor)

window.app = new Vue({
  router,
  store,
  mounted,
  render: h => h(App)
}).$mount('#app')
