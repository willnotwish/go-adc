import svc from '../services/auth'

export default {
  install(Vue) {
    Vue.prototype.$auth = svc

    Vue.mixin({
      created() {
        if (this.handleLoginEvent) {
          svc.addListener('loginEvent', this.handleLoginEvent)
        }
      },

      destroyed() {
        if (this.handleLoginEvent) {
          svc.removeListener('loginEvent', this.handleLoginEvent)
        }
      }
    })
  }
}