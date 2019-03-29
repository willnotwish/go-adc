import Vue from 'vue'
import Vuex from 'vuex'
import apiClient from 'api-client'
import { EventBus } from '../event-bus'

Vue.use(Vuex)

const state = {
  initialized: false,
  interval: 0,
  channels: {
    0: { id: 0, samples: [], divisor: 0, shadowDivisor: 0, scale: 1 },
    1: { id: 1, samples: [], divisor: 0, shadowDivisor: 0, scale: 1 },
    2: { id: 2, samples: [], divisor: 0, shadowDivisor: 0, scale: 1 },
    3: { id: 3, samples: [], divisor: 0, shadowDivisor: 0, scale: 1 },
    4: { id: 4, samples: [], divisor: 0, shadowDivisor: 0, scale: 1 },
    5: { id: 5, samples: [], divisor: 0, shadowDivisor: 0, scale: 1 },
    6: { id: 6, samples: [], divisor: 0, shadowDivisor: 0, scale: 1 },
    7: { id: 7, samples: [], divisor: 0, shadowDivisor: 0, scale: 1 },
  },
  channelsList: [0, 1, 2, 3, 4, 5, 6, 7],
  data: {},
  ws: {
    isConnected: false,
    message: '',
    reconnectError: false,
  },
}

const getters = {
  initialized: state => state.initialized,
  interval:    state => state.interval,
  allChannels: state => state.channelsList.map( id => state.channels[id] ),

  samples:        (state, getters) => getters.allChannels.map( c => c.samples ).map( id => state.data[id] ),
  divisors:       (state, getters) => getters.allChannels.map( c => c.divisor ),
  shadowDivisors: (state, getters) => getters.allChannels.map( c => c.shadowDivisor ),

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
    const id = `s-${sampleKey++}`
    state.data[id] = { id: id, ...sample }
    state.channels[sample.Input].samples.push(id)
  },
  clearSamples(state) {
    Object.keys(state.channels).forEach( ch => {
      state.channels[ch].samples = []
    })
    state.data = {}
  },
  setDivisor(state, payload) {
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
  SOCKET_ONMESSAGE(state, message)  {
    state.ws.message = message
  },
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