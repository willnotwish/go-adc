<template lang='pug'>
  form( v-on:submit.prevent="submit" )
    .field.is-horizontal
      .field-label
        label.label Presets
      .field-body
        .field.is-grouped
          .control( v-for="d in presets")
            button.button.is-small( v-bind:class="presetButtonHtmlClass(d)" v-on:click.prevent="divisor = d" v-bind:disabled="divisor == d") {{presetText(d)}}
    .field.is-horizontal
      .field-label.is-normal
        label.label Divisor
      .field-body
        .field
          .control
            input.input( type="number" v-model="divisor" v-bind:disabled="isSaving" )
          p.help( v-bind:class="hintHtmlClass" ) {{ hint }}

    .field.is-horizontal
      .field-label
      .field-body
        .field.is-grouped
          .control
            button.button.is-primary( type="submit" v-bind:disabled="isSaving || isUnchanged" ) Save changes
          .control
            button.button.is-light( v-on:click="done" v-bind:disabled="isSaving" ) Cancel
</template>

<script>

export default {
  name: 'DivisorForm',
  props: {
    channel: Number,   // 0-7
  },
  data() {
    return {
      isSaving: false,
      presets: [0, 1, 10, 100, 600]
    }
  },
  computed: {
    divisor: {
      set(d) {
        const payload = { input: this.channel, divisor: d }
        this.$store.commit('setShadowDivisor', payload)
      },
      get() {
        return this.$store.getters.channels[this.channel].shadowDivisor
        // return this.$store.getters.shadowDivisors[this.channel]
      }
    },
    checked: {
      get() {
        return this.divisor > 0
      },
      set(v) {
        if (v) {
          this.divisor = 100
        } else {
          this.divisor = 0
        }
      }
    },
    interval() {
      return this.$store.getters.interval
    },
    pin() {
      return this.channel + 1
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
        return `Sample once every ${this.divisor * this.interval/1000} secs`
      } else {
        return `Sample at ${this.frequency.toFixed(2)} Hz`
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
      if (this.isSaving) {
        return "Saving... please wait"
      } else if (this.divisor <= 0) {
        return "Channel disabled: no samples will be taken"
      } else {
        return this.frequencyAsText
      }
    },
    hintHtmlClass() {
      return ''
    },
    icon() {
      if (this.hasChanges) {
        return "sync"
      }
      else {
        return 'coffee'
      }
    },
    hasChanges() {
      return !this.isUnchanged
    },
    isUnchanged() {
      return this.divisor == this.$store.getters.divisors[this.channel]
    }
  },
  methods: {
    submit() {
      this.isSaving = true
      this.$store
        .dispatch('updateInputDivisor', {input: this.channel, divisor: Number(this.divisor)})
        .then( () => {
          this.isSaving = false
          this.done()
        }).catch( err => {
          this.isSaving = false
          console.log("Error saving changes: ", err)
        })
    },
    done() {
      this.$emit('done')
    },
    enableToggle(e) {
      console.log("Enable toggle: ", e.target.value)
    },
    usePreset(index) {
      const preset = this.presets[index]
      if (preset) {
        console.log( "Preset divisor: ", preset.divisor )
        this.divisor = preset.divisor
      } else {
        console.log( "Invalid preset. Ignoring...")
      }
    },
    presetText(d) {
      switch(d) {
        case 0:
          return "Never (OFF)"
        case 1:
          return "10x/sec"
        case 10:
          return "Every sec"
        case 100:
          return "Every 10 secs"
        case 600:
          return "Once per min"
        default:
          return "Custom"
      }
    },
    presetButtonHtmlClass(d) {
      if (d == 0) {
        return "is-danger"
      }
      return "is-primary"
    }
  }
}
</script>