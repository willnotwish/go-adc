<template lang='pug'>
  .c-config
    section-title( title="Settings" subtitle="Sampling setup" )
    .columns
      .column
        config-item( v-for="(d, i) in divisors" v-bind:input="i" v-bind:interval="interval" v-bind:divisor="d" )
      .column
        p.content
          |Each of the eight inputs can be sampled at a maximum of {{maxFrequency}}Hz. That's once every {{interval}}ms.
        p.content
          |You may control how often an input is sampled by setting its integer <strong><em>divisor</em></strong>.
          |The higher the divisor, the lower the sampling rate.
        p.content
          |The lowest divisor (1) yields the fastest rate ({{maxFrequency}}Hz).
          |Increasing the divisor to 10 reduces the rate to {{maxFrequency/10}}Hz (once per second).
          |For once a minute, set the divisor to {{minuteDivisor}}; for once an hour, {{hourDivisor}} will do it.
        p.content Each input has its own divisor. A value of 0 disables sampling altogether on that input.

</template>

<script>
// import InputDivisor from './input-divisor'
import ConfigItem from './config-item'
import SectionTitle from './section-title'

export default {
  name: 'Config',
  computed: {
    divisors() {
      return this.$store.getters.divisors
    },
    interval() {
      return this.$store.getters.interval
    },
    maxFrequency() {
      return 1000/this.interval
    },
    minuteDivisor() {
      return 60000/this.interval
    },
    hourDivisor() {
      return 60 * this.minuteDivisor
    }
  },
  components: {
    // "input-divisor": InputDivisor,
    "section-title": SectionTitle,
    "config-item": ConfigItem
  }
}
</script>

<style type="scss">
  .input-divisor.field.is-horizontal {
    margin-bottom: 1.2em;
  }
</style>