<template lang='pug'>
  nav.top-nav
    router-link( v-bind:to="{ name: 'dashboard' }" ) Dashboard
    |&nbsp;|&nbsp;
    router-link( v-bind:to="{ name: 'config' }" ) Settings
    |&nbsp;|&nbsp;
    router-link( v-bind:to="{ name: 'about' }" ) About
    |&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
    a( v-if="isAuthenticated" href="#" v-on:click.prevent="logout" ) Log out
    a( v-else href="#" v-on:click.prevent="login('CSTATE')") Log in
    |&nbsp;|&nbsp;
    router-link( v-if="isAuthenticated" v-bind:to="{ name: 'profile' }") Profile
</template>

<script>
export default {
  name: 'TopNav',
  data() {
    return {
      isAuthenticated: false
    }
  },
  methods: {
    login(cstate) {
      this.$auth.login(cstate)
    },
    logout() {
      this.$auth.logOut()
    },
    handleLoginEvent(data) {
      this.isAuthenticated = data.loggedIn
    }
  },
}
</script>
