<template lang='pug'>
  header
    h2.title.is-2 {{ title }}
    h3.subtitle.is-4( v-if="subtitle" ) {{ subtitle }}
    hr
</template>

<script>
export default {
  name: 'SectionTitle',
  props: ['title', 'subtitle'],
}
</script>
