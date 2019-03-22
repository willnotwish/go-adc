<template lang='pug'>
  .c-loading-led
    led-indicator( v-bind:value="value" )
</template>

<script>

import LEDIndicator from './led-indicator'

export default {
  name: 'LoadingLED',
  props: ['value'],
}
</script>

<!-- <style type="scss">
  .input-divisor.field.is-horizontal {
    margin-bottom: 1.2em;
  }
</style> -->