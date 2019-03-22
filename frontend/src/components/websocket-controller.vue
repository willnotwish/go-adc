<template lang='pug'>
  .ws-controller
    button.button( v-on:click="toggleConnection" v-bind:class="buttonModifier" )
      font-awesome-icon( v-bind:icon="icon" )
      |&nbsp;&nbsp;{{ buttonText }}

</template>

<script>
import LEDIndicator from './led-indicator'

export default {
  name: 'WebsocketController',
  data() {
    return {
      activity: 0
    }
  },
  computed: {
    isConnected() {
      return this.$store.getters.isWsConnected
    },
    lastSample() {
      return this.$store.getters.wsMessage
    },
    reconnectError() {
      return this.$store.getters.reconnectError
    },
    buttonText() {
      return this.isConnected ? "Disconnect" : "Connect"
    },
    buttonModifier() {
      return this.isConnected ? "is-light" : "is-primary"
    },
    connectionStatus() {
      return this.isConnected ? "CONNECTED" : "DISCONNECTED"
    },
    icon() {
      return this.isConnected ? 'stop': 'download'
    }
  },
  methods: {
    toggleConnection() {
      if (this.isConnected) {
        this.$disconnect()
      } else {
        this.$connect()
      }
    }
  },
  components: {
    "led-indicator": LEDIndicator
  },
  created() {
    this.$options.sockets.onmessage = () => {
      if (this.activity == 0) {
        this.activity = 1
        setTimeout( () => this.activity = 0, 200 )
      }
    }
  }
}
</script>