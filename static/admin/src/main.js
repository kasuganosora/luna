import Vue from 'vue'
import App from './App.vue'
import router from './router'
import './plugins/element.js'
import store from './store'

Vue.config.productionTip = false

const mounted = () => {
  let self = this
  window.onresize = function(){
    self.$store.state.screenWidth = document.documentElement.clientWidth;
    self.$store.state.screenHeight = document.documentElement.clientHeight;
  }
}

window.app = new Vue({
  router,
  store,
  mounted,
  render: h => h(App)
}).$mount('#app')
