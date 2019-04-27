import {redis, xrange} from './client'
import Debug from 'debug'

const debug = Debug('obd-api:redis-store')

// A redis XRANGE query yields something like this:
// [["1555501380946-0",["0C","1163"]],["1555501381126-0",["05","85"]]]

const deserialize = (record) => {
  return { id: record[0], pid: record[1][0], value: record[1][1] }
}

const store = {

  // Returns a promise. Requires a connection
  getData(options) {
    debug("getData. options: ", options)
    if (!redis.connected) {
      throw new Error("Redis not connected. Try again later")
    }
    return xrange('obd-stream', '-', '+')
      .then( data => data.map(deserialize) )
  }
}

export default store