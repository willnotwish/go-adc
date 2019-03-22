<template lang='pug'>
  .input-divisor.field.is-horizontal
    .field-label.is-normal
      label.label Input {{ channel }}
    .field-body
      .field
        .control.has-icons-right
          input.input( type="number" v-model.number.lazy="divisor" v-on:input="onInput")
          //- span.icon.is-small.is-right
            //- font-awesome-icon( v-bind:icon="icon" )
        //- p.help( v-bind:class="hintHtmlClass" ) {{ hint }}
</template>

<script>

export default {
  name: 'DivisorEditor',
  props: {
    input: Number, // 0-7
    interval: Number, // in milliseconds
  },
  data() {
    return {
      updating: false,
      state: 'idle'
    }
  },
  computed: {
    divisor: {
      get() {
        console.log("Getting divisor from store for input: ", this.input)
        return this.$store.getters.divisors[this.input]
      },
      set(newValue) {
        console.log("Setting divisor in store: ", newValue)
        this.state = 'updating'
        this.$store
          .dispatch('updateInputDivisor', {input: this.input, divisor: newValue})
          .then( () => this.state = 'idle' )
          .catch( err => this.state = 'error' )
      }
    },
    channel() {
      return this.input
    },
    isSampling() {
      return this.state == 'idle' && this.divisor > 0
    },
    isUpdating() {
      return this.state == 'updating'
    },
    isChanging() {
      return this.state == 'changing'
    },
    isDisabled() {
      return this.state == 'idle' && this.divisor == 0
    },
    frequency() {
      return this.divisor > 0 ? 1000/(this.interval*this.divisor) : 0
    },
    formattedFrequency() {
      if (this.frequency < 0.5) {
        return `once every ${this.divisor * this.interval/1000} secs`
      } else {
        return `at ${this.frequency.toFixed(2)} Hz`
      }
    },
    helpMsg() {
      if (this.isUpdating) {
        return 'Updating... please wait'
      } else if (this.isSampling) {
        return `Sampling ${this.formattedFrequency}`
      } else if (this.isChanging) {
        return `Sample ${this.formattedFrequency}. ENTER to update`
      } else if (this.isDisabled) {
        return 'Off'
      }
    },
    hint() {
      const storedValue = this.$store.getters.divisors[this.input]
      if (storedValue != this.changedValue) {
        return "Needs saving"
      }
      return "Default hint"
    },
    hintHtmlClass() {
      return 'is-success'
    },
    helpModifier() {
      if (this.isSampling) {
        return 'is-success'
      }
    },
    icon() {
      if (this.isSampling) {
        return 'check'
      } else if (this.isUpdating) {
        return 'sync'
      } else if (this.isChanging) {
        return 'info'
      } else {
        return 'ban'
      }
    }
  },
  methods: {
    onInput(e) {
      const v = e.target.value
      // this.changedValue = parseInt(v, 10)
      this.$emit('input', v)
      console.log("onInput. e.target.value: ", v)
    },
    // onChange(e) {
    //   console.log("onChange. e.target.value, value, current", e.target.value, this.value, this.changedValue)
    //   this.state = 'updating'
    //   this.$store
    //     .dispatch('updateInputDivisor', {input: this.input, divisor: e.target.value})
    //     .then( () => this.state = 'idle' )
    // },
    // onUpdate() {
    //   this.state = 'updating'
    //   this.$store
    //     .dispatch('updateInputDivisor', {input: this.input, divisor: this.cachedDivisor})
    //     .then( () => this.state = 'idle' )
    // },
  }
}
</script>