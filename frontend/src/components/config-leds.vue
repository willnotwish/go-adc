<template lang='pug'>
  .c-config-indicator
    .tags.has-addons
      led-indicator( v-for="(d,i) in divisors" v-bind:name="i" v-bind:value="d" )
</template>

<script>
import LEDIndicator from './led-indicator'

export default {
  name: 'ConfigLEDs',
  computed: {
    divisors() {
      return this.$store.getters.divisors
    },
  },
  components: {
    "led-indicator": LEDIndicator
  }
}
</script>

<style type="scss">
  .input-divisor.field.is-horizontal {
    margin-bottom: 1.2em;
  }
</style>