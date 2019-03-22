import samples from './data/samples'
import moment from 'moment'
import authService from '../../auth/auth-service'

const _fetchSamples = (mockData, delay = 0) => {
  return new Promise((resolve) => {
    setTimeout( () => resolve(mockData), delay )
  })
}

const _updateInputDivisor = (payload, delay = 0) => {
  return new Promise((resolve) => {
    console.log(`MOCK update with payload. Resolving in ${delay} ms`, payload)
    setTimeout( () => resolve(true), delay )
  })
}

const mockConfig = {
  Interval: 100,
  Divisors: [10, 0, 0, 0, 0, 0, 0, 0],
}

const _fetchConfig = (delay = 0) => {
  return new Promise((resolve) => {
    console.log(`MOCK config fetch. Resolving ${mockConfig} in ${delay} ms`)
    setTimeout( () => resolve(mockConfig), delay )
  })
}

const withAccessToken = () => {
  return authService.getAccessToken()
    .then( token => {
      console.log("With access token: ", token)
      return token
    })
}

export default {
  fetchConfig() {
    return withAccessToken().then( () => _fetchConfig( 200 ) )
  },

  fetchSamples(options) {
    // console.log("MOCK fetchSamples. options: ", options)
    let url = '/api/data'
    if (options && options.since) {
      // url = `/api/data?s=${moment(options.since).toISOString()}`
      url += `?s=${options.since}`
    }
    // console.log("MOCK fetchSamples. URL: ", url)
    return _fetchSamples(samples, 200) // wait 200ms before returning samples
  },

  updateInputDivisor(payload) {
    return _updateInputDivisor(payload, 2000)
  }
}
