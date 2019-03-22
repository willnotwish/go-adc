<template lang='pug'>
  .modal(v-bind:class="active")
    .modal-background(v-on:click="done")
    .modal-card
      header.modal-card-head
        .modal-card-title
          h4.title.is-4 Channel {{ channel }}
        button.delete( aria-label="close" v-on:click="done")
      section.modal-card-body
        divisor-form( v-bind:channel="channel" v-on:done="done" )
      footer.modal-card-foot
        p Change the sampling rate by choosing one of the presets; for finer control, adjust the divisor directly. Be sure to save your changes.
</template>

<script>

import DivisorForm from './divisor-form'

export default {
  name: 'EditModal',

  props: {
    channel: Number
  },

  components: {
    'divisor-form': DivisorForm
  },

  computed: {
    active() {
      return this.channel >= 0 ? 'is-active' : ''
    }
  },

  methods: {
    done() {
      this.$emit('done')
    },
  }
}
</script>
