<template lang='pug'>
.divisor-summary(v-bind:class="rootHtmlClass")
  .box
    .level
      .level-left
        .level-item
          h4.title.is-4(v-bind:title="tooltip") Channel {{ input }}
        .level-item
          .tag(v-bind:class="status.htmlClass") {{ status.text }}
      .level-right
        .level-item
          a.button(href="#" v-on:click.prevent="showModal")
            font-awesome-icon( icon='edit' )
            |&nbsp;Edit

    .line
      span Divisor <strong>{{ divisor }}</strong>
      br
      span( v-if="isEnabled") <em>currently sampling {{frequencyAsText}}</em>

  .modal(v-bind:class="modalHtmlClass")
    .modal-background(v-on:click="hideModal")
    .modal-card
      header.modal-card-head
        .modal-card-title 
          h4.title.is-4 Channel {{ input }}
        button.delete( aria-label="close" v-on:click="hideModal")
      section.modal-card-body
        divisor-form( v-bind:channel="input" v-bind:interval="interval" v-on:done="hideModal" )
      footer.modal-card-foot
        p Change the sampling rate by choosing one of the presets; for finer control, adjust the divisor directly. Be sure to save your changes.

</template>

<script>

import DivisorForm from './divisor-form'

export default {
  name: 'ConfigItem',
  props: {
    input: Number,    // 0-7
    interval: Number, // in milliseconds
    divisor: Number,  // reactive
  },
  data() {
    return {
      modalHtmlClass: '',
      rootHtmlClass: '',
    }
  },
  computed: {
    pin() {
      return this.input + 1
    },
    tooltip() {
      return `MCP3008 pin ${this.pin}`
    },
    frequency() {
      return this.divisor > 0 ? 1000/(this.interval*this.divisor) : 0
    },
    isEnabled() {
      return this.divisor > 0
    },
    frequencyAsText() {
      if (this.frequency < 0.5) {
        return `once every ${this.divisor * this.interval/1000} secs`
      } else {
        return `at ${this.frequency.toFixed(2)} Hz`
      }
    },
    status() {
      if (this.isEnabled) {
        return { htmlClass: 'is-success', text: 'enabled'}
      } else {
        return { htmlClass: 'is-light', text: 'disabled'}
      }
    },
  },
  methods: {
    showModal() {
      this.modalHtmlClass = 'is-active'
    },
    hideModal() {
      this.modalHtmlClass = ''
    },
  },
  components: {
    'divisor-form': DivisorForm,
  }
}
</script>