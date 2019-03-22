<template lang='pug'>
  .c-samples
    section-title( title="Samples" )
    table.table
      thead
        tr
          th Input
          th Voltage
          th Time
          th Date
      tbody
        sample-tr( v-for="s in samples" v-bind:sample="s" )
      tfoot
</template>

<script>

import SectionTitle from './section-title'
import SampleAsTableRow from './sample-as-tr'

export default {
  name: 'Samples',
  computed: {
    samples() {
      return this.$store.getters.samples
    }
  },
  methods: {
    f(ts) {
      const d = new Date(ts)
      return d.toLocaleTimeString()
    }
  },
  components: {
    "sample-tr": SampleAsTableRow,
    "section-title": SectionTitle
  }
}
</script>
