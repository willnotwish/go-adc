<template lang='pug'>
  .dashboard
    header
      .notification.is-warning( v-if="!isLoggedIn" ) Log in now to see the interesting stuff.

      h2.title.is-2 Dashboard
    hr
    .level( v-if="isLoggedIn" )
      .level-left
        .level-item
          websocket-controller
      .level-right
        .level-item
          button.button( v-on:click="initialize" v-bind:disabled="!isLoggedIn" )
            font-awesome-icon( icon="sync" )
            |&nbsp;&nbsp;Sync with server
        .level-item
          button.button( v-on:click="clearData" )
            font-awesome-icon( icon="trash" )
            |&nbsp;&nbsp;Clear local data

    .notification.is-danger( v-if="error" )
      button.delete
      |{{error}}

    template( v-if="initialized" )
      hr
      .columns
        .column( v-for="channel in allChannels" )
          channel-card( v-bind:channel="channel" )
      .columns
        .column( v-for="ch in allChannels" )
          button.button.is-fullwidth.is-small(
            v-bind:class="selectChannelButtonModifier(ch)"
            v-bind:disabled="ch.samples.length === 0"
            v-on:click="toggleChannelSelection(ch)" ) {{selectChannelButtonText(ch)}}

      .columns
        .column.is-half.is-offset-one-quarter
          chart-recorder( v-bind:chartData="chartRecorderData" v-bind:options="chartRecorderOptions" )
</template>

<script>

import ChannelSummary      from './channel-summary'
import ChannelCard         from './channel-card'
import SectionTitle        from './section-title'
import LineChartRecorder   from './line-chart-recorder'
import EditModal           from './edit-modal'
import ConfigLEDs          from './config-leds'
import LEDIndicator        from './led-indicator'
import WebsocketController from './websocket-controller'

import moment from 'moment'

export default {
  name: 'Dashboard',

  props: {
    title: String,
  },

  data() {
    return {
      error: false,
      loading: 0,
      selectedChannels: [],
    }
  },

  components: {
    'channel-summary': ChannelSummary,
    'channel-card':    ChannelCard,
    'section-title':   SectionTitle,
    'chart-recorder':  LineChartRecorder,
    'edit-modal':      EditModal,
    'config-leds':     ConfigLEDs,
    'led-indicator':        LEDIndicator,
    'websocket-controller': WebsocketController,
  },

  computed: {
    allChannels() {
      return this.$store.getters.allChannels
    },
    initialized() {
      return this.$store.getters.initialized
    },
    isLoggedIn() {
      return this.$auth.isAuthenticated()
    },
    divisors() {
      return this.$store.getters.divisors
    },
    interval() {
      return this.$store.getters.interval
    },
    total() {
      var count = 0
      this.$store.getters.samples.forEach( (a) => count += a.length )
      return count
    },
    hasData() {
      return this.total > 0
    },
    chartRecorderData() {
      const data = this.$store.getters.data
      const datasets = []
      this.selectedChannels.forEach( (ch) => {
        const samples = this.$store.getters.channels[ch].samples

        const xydata = []
        if (samples.length > 0) {
          const t0 = moment().subtract( 5, 'minutes' )
          for( var i = samples.length-1; i--; i >= 0 ) {
            const sample = data[samples[i]]
            const ts = moment( sample.Timestamp )
            if (ts.isBefore(t0)) {
              break // too old so don't bother with the rest
            } else {
              xydata.unshift( { x: ts, y: sample.Voltage } )
            }
          }
        }
        datasets.push( {label: `Ch ${ch}`, data: xydata} )
      })
      return { datasets: datasets }
    },
    loadingIndication() {
      if (this.error) {
        return -1
      }
      else {
        return this.loading
      }
    },
  },

  created() {
    this.chartRecorderOptions = {
      animation: {
        duration: 0
      },
      scales: {
        xAxes: [{
          type: 'realtime',
          time: {
            displayFormats: {
              second: 'HH:mm:ss',
            },
          },
          realtime: {
            duration: 60000,  // ms
            refresh:   1000,
            delay:     1000,
            pause:    false,
            ttl:      undefined,
          }
        }],
        yAxes: [{
          display: true,
          ticks: {
            beginAtZero: true,
            max: 4,
            stepSize: 1
          },
        }],
      }
    }
  },

  methods: {
    fetchSamples() {
      this.$store.dispatch('fetchSamples', {commit: 'append', scope: 'new'})
    },
    clearSamples() {
      this.$store.commit('clearSamples')
    },
    toggleChannelSelection(ch) {
      const index = this.selectedChannels.indexOf(ch.id)
      if (index == -1) {
        this.selectedChannels.push(ch.id)
        this.selectedChannels.sort( (a,b) => a - b )
      } else {
        this.selectedChannels.splice(index, 1)
      }
    },
    selectChannelButtonText(ch) {
      return this.isChannelSelected(ch) ? "Hide" : "Show"
    },
    isChannelSelected(ch) {
      return this.selectedChannels.indexOf(ch.id) != -1
    },
    selectChannelButtonModifier(ch) {
      if (this.isChannelSelected(ch)) {
        return "is-light"
      } else {
        return "is-dark"
      }
    },

    initialize() {
      this.$store
        .dispatch('initialize')
        .then( () => {
          this.error = false
          return true
        }).catch( err => {
          console.log("Error dispatching 'initialize' action: ", err)
          this.error = err
        })
    },
    clearData() {
      this.$store.commit('clearSamples')
    },
    fetchNewData() {
      this.loading = 1
      this.$store
        .dispatch( 'fetchData', {commit: 'append', scope: 'new'} )
        .then( data => {
          this.error = false
          console.log( "Fetched data: ", data )
        }).catch( err => {
          console.log( "Error fetching new data: ", err )
          this.error = err
        }).finally( () => this.loading = 0 )
    },
  }
}
</script>