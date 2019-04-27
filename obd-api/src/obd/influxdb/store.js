import influxdb from './client'
import Debug from 'debug'

const debug = Debug('obd-api:influxdb-store')

const deserialize = (record) => {
  return record // pass through for now
}

const store = {

  // Returns a promise. Not sure yet if this requires a connection, or if requests are queued (Redis style)
  getData(options) {
    debug("getData. options: ", options)

    return influxdb.query('select * from "pid-values"').then( result => {
      debug("getData result: ", result)
      return result
    })
  }
}

export default store