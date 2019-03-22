<template lang='pug'>
  .led( v-bind:class="modifier" )
</template>

<script>
export default {
  name: 'LEDIndicator',
  props: ['name', 'value'],
  computed: {
    title() {
      return `${this.name} status: ${this.status}`
    },
    status() {
      if (this.value > 0) {
        return 'on'
      }
      else if (this.value == 0) {
        return 'off'
      } else {
        return 'unknown'
      }
    },
    modifier() {
      if (this.value <= 0) {
        return "led--off"
      }
    }
  },
}
</script>

<style lang="scss" scoped>
  $size: 10px;
  $colour: red;
  .led {
    cursor:pointer;
    display: inline-block;
    margin: $size/4;
    margin-bottom: 0px;
    width: $size;
    height: $size;
    border-radius: $size/2;
    background-color: $colour;
    box-shadow: 0px 0px $size $size/4 rgba($colour, .75);
  }

  .led--off{
    background-color: #DDDDDD;
    box-shadow: 0px 0px 0px 0px rgba(255, 255, 190, .75);
  }
</style>