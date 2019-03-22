<template lang='pug'>
.divisor-summary(v-bind:class="rootHtmlClass")
  .box
    .level
      .level-left
        .level-item
          h4.title.is-4(v-bind:title="tooltip") Input {{ input }}
        .level-item
          .tag(v-bind:class="status.htmlClass") {{ status.text }}
      .level-right
        .level-item
          a(href="#" v-on:click.prevent="showModal") Edit

    .line
      span Divisor <strong>{{ divisor }}</strong>
      br
      span( v-if="isEnabled") <em>currently sampling {{frequencyAsText}}</em>

  .modal(v-bind:class="modalHtmlClass")
    .modal-background(v-on:click="hideModal")
    .modal-card
      header.modal-card-head
        h4.modal-card-title Set sampling frequency
        button.delete( aria-label="close" v-on:click="hideModal")
      form( v-on:submit.prevent="saveChanges")
        section.modal-card-body
          .field.is-horizontal
            .field-label.is-normal
              label.label Input {{ input }}
            .field-body
              .field
                .control.has-icons-right
                  input.input( type="number" v-model="newDivisor" v-bind:disabled="isSaving" )
                  span.icon.is-small.is-right
                    font-awesome-icon( v-bind:icon="icon" )
                p.help( v-bind:class="hintHtmlClass" ) {{ hint }}

        footer.modal-card-foot
          button.button.button.is-success( type="submit" ) Save changes
          button.button(v-on:click="hideModal") Cancel
</template>

<script>

export default {
  name: 'DivisorSummary',
  props: {
    input: Number,    // 0-7
    interval: Number, // in milliseconds
    divisor: Number,  // reactive
  },
  data() {
    return {
      modalHtmlClass: '',
      rootHtmlClass: '',
      hintHtmlClass: '',
      isSaving: false,
    }
  },
  computed: {
    newDivisor: {
      get() {
        return this.cachedDivisor()
      },
      set(d) {
        this.pendingDivisor = d
      }
    },
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
    hint() {
      return "TODO: make hint context sensitive"
    },
    icon() {
      return "sync"
    }
  },
  methods: {
    cachedDivisor() {
      if (this.isSaving) {
        return this.pendingDivisor
      }
      else {
        return this.divisor
      }
    },
    showModal() {
      console.log("About to show modal")
      this.modalHtmlClass = 'is-active'
    },
    hideModal() {
      console.log("About to hide modal")
      this.modalHtmlClass = ''
    },
    saveChanges() {
      console.log("About to save changes by updating divisor with new value of: ", this.pendingDivisor)
      this.isSaving = true
      this.$store
        .dispatch('updateInputDivisor', {input: this.input, divisor: Number(this.pendingDivisor)})
        .then( () => {
          this.isSaving = false
          this.hideModal()
        }).catch( err => {
          this.isSaving = false
          console.log("Error saving changes: ", err)
        })
    },
    onInput(e) {
      console.log("onInput. Target value: ", e.target.value)
      this.newDivisor = e.target.value
    }
  }
}
</script>