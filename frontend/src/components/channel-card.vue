<template lang='pug'>
  .card
    .card-header
      .card-header-title
        h6.title.is-6 Ch {{ channel.id }} &nbsp;
          led-indicator( v-bind:value="activity" )
    .card-content
      p( v-if="hasData") {{ voltage }} V
      p( v-else ) &mdash;
      p.is-size-7( v-if="hasData") {{count}}
      p.is-size-7( v-else ) &mdash;

</template>

<script>

import LED from './led-indicator'
import { EventBus } from '../event-bus'

export default {
  name: 'ChannelCard',
  props: {
    channel: Object,
  },
  data() {
    return {
      activity: false
    }
  },
  computed: {
    data() {
      return this.$store.getters.channelData( this.channel.id )
    },
    last() {
      const len = this.data.length
      if (len > 0) {
        return this.data[len-1]
      }
    },
    count() {
      return this.data.length
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
      return this.data.length > 0
    },
    divisor() {
      return this.channel.divisor
    },
  },
  methods: {
    onActivity(ch) {
      if (ch == this.channel.id ) {
        if (this.activity == 0) {
          this.activity = 1
          setTimeout( () => this.activity = 0, 200 )
        }
      }
    }
  },
  components: {
    'led-indicator': LED
  },
  created() {
    EventBus.$on( 'data-appended', this.onActivity )
  },
  destroyed() {
    EventBus.$off( 'data-appended', this.onActivity )
  }
}
</script>
