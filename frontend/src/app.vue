<template lang='pug'>
  .container( id="app" )
    img( alt="App logo" src="./assets/logo.png" )
    top-nav
    hr
    .notification.is-danger( v-if="error")
      button.delete
      |{{error}}

    router-view

</template>

<script>
import TopNav from './components/top-nav'

export default {
  name: 'App',
  components: {
    'top-nav': TopNav
  },
  data() {
    return {
      error: false
    }
  },
  computed: {
    counts() {
      return this.$store.getters.sampleCounts
    }
  },
  methods: {
    handleLoginEvent(data) {
      console.log("App handle login event. data: ", data)
      if (data.loggedIn) {
        this.initialize()
      }
    },
    initialize() {
      this.$store
        .dispatch('initialize')
        .then( () => {
          console.log("Fetched config. Divisors: ", this.$store.getters.divisors)
          this.error = false
          return true
        }).catch( err => {
          console.log("Error dispatching 'initialize' action: ", err)
          this.error = err
        })
    },
  }
}
</script>
