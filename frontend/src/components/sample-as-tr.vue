<template lang='pug'>
  tr
    td {{ input }}
    td {{ voltage }} V
    td {{ time }}
    td {{ date }}
</template>

<script>
export default {
  name: 'SampleAsRow',
  props: ['sample'],
  computed: {
    input() {
      return this.sample.Input
    },
    voltage() {
      return this.sample.Voltage
    },
    ts() {
      return new Date(this.sample.Timestamp)
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
