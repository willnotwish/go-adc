import axios from 'axios'
import auth from '../../services/auth'

const api = axios.create({
  baseURL: '/api',
  headers: {'Content-Type': 'application/json'},
})

api.interceptors.request.use( config => {
  return auth.getAccessToken().then( token => {
    config.headers['Authorization'] = `Bearer ${token}`
    return config
  })
})

api.interceptors.response.use( response => {
  return response
}, error => {
  console.log("intercepted api error, which will now be rejected: ", error)
  return Promise.reject(error)
})

export default {
  fetchConfig() {
    return api.get( 'config' )
      .then( response => response.data )
  },
  fetchSamples( opts={} ) {
    return api.get( 'data', { params: {s: opts.since} } )
      .then( response => response.data )
  },
  updateInputDivisor( payload ) {
    return api.post( 'control', null, { params: {i: payload.input, d: payload.divisor} } )
  },
}