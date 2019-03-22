<template lang='pug'>
  .column
    .card
      .card-header
        .card-header-title
          h6.title.is-6 Ch {{ channel }} &nbsp;
            led-indicator( v-bind:value="activity" )
      .card-content
        p( v-if="hasData") {{ voltage }} V
        p( v-else ) &mdash;
        p.is-size-7( v-if="hasData") {{count}}
        p.is-size-7( v-else ) &mdash;
      //- .card-footer
      //-   .card-footer-item
      //-     button.button.is-small( v-bind:class="buttonClass" v-on:click="toggleVisibility" v-bind:disabled="buttonDisabled") {{ buttonText }}

</template>

<script>

import LED from './led-indicator'
export default {
  name: 'ChannelSummary',
  props: {
    channel: Number,
  },
  computed: {
    isVisible() {
      return this.$store.getters.channels[this.channel].visibility
    },
    samples() {
      return this.$store.getters.channelData(this.channel)
      // return this.$store.getters.samples[this.channel]
    },
    last() {
      if( this.hasData ) {
        const len = this.samples.length
        return this.samples[len-1]
      }
    },
    count() {
      return this.samples.length
    },
    voltage() {
      return this.last.Voltage
    },
    ts() {
      return new Date(this.last.Timestamp)
    },
    time() {
      return this.ts.toLocaleTimeString()
    },
    date() {
      return this.ts.toLocaleDateString()
    },
    hasData() {
      return this.samples.length > 0
    },
    divisor() {
      return this.$store.getters.divisors[this.channel]
    },
    buttonText() {
      return this.isVisible ? 'OFF' : 'ON'
    },
    buttonClass() {
      return this.isVisible ? 'is-success' : 'has-background-light'
    },
    buttonDisabled() {
      return this.divisor <= 0
    },
    activity() {
      return this.divisor
    }
  },
  methods: {
    toggleVisibility() {
      this.$emit('toggle-visibility', this.channel)
    },
  },
  components: {
    'led-indicator': LED
  },
}
</script>
