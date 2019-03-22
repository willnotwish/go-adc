<template lang='pug'>
  .sample.column
    .card
      h4.title.is-4 Ch {{ channel }}
      template( v-if="src" )
        .line {{ voltage }} V
        .line {{ time }}
        .line {{ date }}
      p( v-else ) No data
</template>

<script>
export default {
  name: 'Sample',
  props: ['channel', 'src'],
  computed: {
    voltage() {
      return this.src.Voltage
    },
    ts() {
      return new Date(this.src.Timestamp)
    },
    time() {
      return this.ts.toLocaleTimeString()
    },
    date() {
      return this.ts.toLocaleDateString()
    }
  },
}
</script>
