import Vue from 'vue'
import Vuex from 'vuex'
import apiClient from 'api-client'
import { EventBus } from '../event-bus'

Vue.use(Vuex)

const state = {
  initialized: false,
  samples: [ [], [], [], [], [], [], [], [] ],
  interval: 0,
  divisors: [0, 0, 0, 0, 0, 0, 0, 0],
  shadowDivisors: [0, 0, 0, 0, 0, 0, 0, 0],
  // visibilities: [false, false, false, false, false, false, false, false],
  activities: [false, false, false, false, false, false, false, false],

  // Websocket
  ws: {
    isConnected: false,
    message: '',
    reconnectError: false,
  },

  channels: {
    0: { id: 0, samples: [], active: false, divisor: 0, shadowDivisor: 0, enabled: false },
    1: { id: 1, samples: [], active: false, divisor: 0, shadowDivisor: 0, enabled: false },
    2: { id: 2, samples: [], active: false, divisor: 0, shadowDivisor: 0, enabled: false },
    3: { id: 3, samples: [], active: false, divisor: 0, shadowDivisor: 0, enabled: false },
    4: { id: 4, samples: [], active: false, divisor: 0, shadowDivisor: 0, enabled: false },
    5: { id: 5, samples: [], active: false, divisor: 0, shadowDivisor: 0, enabled: false },
    6: { id: 6, samples: [], active: false, divisor: 0, shadowDivisor: 0, enabled: false },
    7: { id: 7, samples: [], active: false, divisor: 0, shadowDivisor: 0, enabled: false },
  },
  channelsList: [0, 1, 2, 3, 4, 5, 6, 7],
  data: {},
}

const getters = {
  initialized:      state => state.initialized,
  interval:         state => state.interval,

  allChannels: state => state.channelsList.map( id => state.channels[id] ),

  samples:        (state, getters) => getters.allChannels.map( c => c.samples ).map( id => state.data[id] ),
  divisors:       (state, getters) => getters.allChannels.map( c => c.divisor ),
  shadowDivisors: (state, getters) => getters.allChannels.map( c => c.shadowDivisor ),
  // visbilities:    (state, getters) => getters.allChannels.map( c => c.enabled ),

  // divisors:         state => state.divisors,
  // shadowDivisors:   state => state.shadowDivisors,
  // visibilities:     state => state.visibilities,

  isWsConnected:    state => state.ws.isConnected,
  wsMessage:        state => state.ws.message,
  wsReconnectError: state => state.ws.reconnectError,

  channels: state => state.channels,
  data: state => state.data,
  channelData: state => ch => {
    return state.channels[ch].samples.map( id => state.data[id] )
  }
}

let sampleKey = 0

const mutations = {
  appendSamples(state, samples) {
    samples.forEach( sample => this.appendData(sample) )
  },
  appendData(state, sample) {
    state.samples[sample.Input].push( sample )

    const id = `s-${sampleKey++}`
    state.data[id] = { id: id, ...sample }
    state.channels[sample.Input].samples.push(id)
  },
  clearSamples(state) {
    for( var i = 0; i <= 7; i++ ) {
      Vue.set( state.samples, i, [] )
    }

    Object.keys(state.channels).forEach( ch => {
      state.channels[ch].samples = []
    })
    state.data = {}
  },
  setDivisor(state, payload) {
    Vue.set( state.divisors, payload.input, payload.divisor )
    state.channels[payload.input].divisor = payload.divisor
  },
  // RESET_DIVISOR(state, payload) {
  //   let ch = channels[payload.input]
  //   if (!ch) {
  //     ch
  //   }
  //   channels[payload.input].divisor = payload.divisor
  //   channels[payload.input].shadowDivisor = payload.divisor
  // },
  setShadowDivisor(state, payload) {
    Vue.set( state.shadowDivisors, payload.input, payload.divisor )
    state.channels[payload.input].shadowDivisor = payload.divisor
  },
  setInterval(state, interval) {
    state.interval = interval
  },
  SET_INTERVAL(state, interval) {
    state.interval = interval
  },
  initialize(state) {
    state.initialized = true
  },
  // setVisibility(state, payload) {
  //   Vue.set( state.visibilities, payload.channel, payload.value)
  //   state.channels[payload.input].visibility = payload.value
  // },
  SOCKET_ONOPEN(state, event)  {
    Vue.prototype.$socket = event.currentTarget
    state.ws.isConnected = true
  },
  SOCKET_ONCLOSE(state, event)  {
    state.ws.isConnected = false
  },
  SOCKET_ONERROR(state, event)  {
    console.error(state, event)
  },
  // default handler called for all methods
  SOCKET_ONMESSAGE(state, message)  {
    state.ws.message = message
  },
  // mutations for reconnect methods
  SOCKET_RECONNECT(state, count) {
    console.info(state, count)
  },
  SOCKET_RECONNECT_ERROR(state) {
    state.ws.reconnectError = true;
  },
}

const actions = {
  initialize({ commit }) {
    return apiClient
      .fetchConfig()
      .then( config => {
        commit( 'setInterval', config.Interval )
        config.Divisors.forEach( (newValue, index) => {
          const payload = {input: index, divisor: newValue}
          commit( 'setDivisor', payload )
          commit( 'setShadowDivisor', payload )
        })
        commit( 'clearSamples' )
        commit( 'initialize' )
        return true
      })
  },

  fetchData({ commit, getters }, options={}) {
    var since
    if (options.scope == 'new') {
      const s0 = getters.lastSample
      if (s0) {
        since = s0.Timestamp
      }
    }
    const mutation = `${(options && options.commit) ? options.commit : 'replace'}Samples`
    return apiClient
      .fetchSamples( {since: since} )
      .then( samples => {
        commit(mutation, samples)
        return samples
      })
  },

  updateInputDivisor({ commit }, payload) {
    return apiClient
      .updateInputDivisor( payload )
      .then( () => {
        commit( 'setDivisor', payload )
        commit( 'setShadowDivisor', payload )
        return true
      })
  },

  APPEND_DATA({ commit }, payload) {
    commit( 'appendData', payload )
    EventBus.$emit('data-appended', payload.Input)
    return Promise.resolve( payload )
  },

}

export default new Vuex.Store({
  state,
  getters,
  actions,
  mutations,
  strict: process.env.NODE_ENV !== 'production'
})

  // lastSamples: (state) => {
  //   const a = []
  //   state.samples.forEach( (b, i) => {
  //     a[i] = b.length > 0 ? b[b.length-1] : undefined
  //   })
  //   return a
  // },

  // lastSample: (state, getters) => {
  //   var last = false
  //   var ts0 = false
  //   getters.lastSamples.forEach( (s) => {
  //     if (s) {
  //       if (last) {
  //         const ts = Date.parse(s.Timestamp)
  //         if (ts < ts0) {
  //           last = s
  //           ts0 = ts
  //         }
  //       } else {
  //         last = s
  //         ts0 = Date.parse(s.Timestamp)
  //       }
  //     }
  //   })
  //   return last
  // },
  // sampleCounts:   state => state.samples.map( a => a.length ),
  // inputDivisor:   state => input => state.divisors[input],
  // data0:          state => state.samples[0],
  // data1:          state => state.samples[1],
  // data2:          state => state.samples[2],
  // data3:          state => state.samples[3],
  // data4:          state => state.samples[4],
  // data5:          state => state.samples[5],
  // data6:          state => state.samples[6],
  // data7:          state => state.samples[7],
