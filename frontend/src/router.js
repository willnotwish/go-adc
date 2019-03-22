import Vue from 'vue'
import VueRouter from 'vue-router'

// Components we can route to
import Dashboard     from './components/dashboard'
import Samples       from './components/samples'
import Config        from './components/config'
import About         from './components/about'
import DivisorEditor from './components/divisor-editor'
import AuthCallback  from './components/auth-callback'
import Profile       from './components/profile'
import NotFound      from './components/not-found'

Vue.use(VueRouter)

export default new VueRouter({
  mode: 'history',
  routes: [
    {
      path: '/dashboard',
      name: 'dashboard',
      component: Dashboard,
    },
    {
      path: '/samples',
      name: 'samples',
      component: Samples,
    },
    {
      path: '/config',
      name: 'config',
      component: Config,
    },
    {
      path: '/divisors/:input/edit',
      name: 'edit_divisor',
      component: DivisorEditor,
    },
    {
      path: '/callback',
      name: 'auth-callback',
      component: AuthCallback,
    },
    {
      path: '/about',
      name: 'about',
      component: About,
    },
    {
      path: '/profile',
      name: 'profile',
      component: Profile
    },
    {
      path: '/',
      name: 'home',
      redirect: '/dashboard'
    },
    {
      path: '*',
      name: 'not_found',
      component: NotFound
    }
  ]
})