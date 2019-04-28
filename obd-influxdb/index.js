import {redis, subscribe} from './redis/client'
import influxDB from './influxdb/client'

import Debug from 'debug'
const debug = Debug('obd-influxdb:index')

redis.on('subscribe', (ch, count) => {
  debug(`Subscribed to channel ${ch}. Count: ${count}`)
})

const log = Debug('obd-influxdb:influxdb')
redis.on('message', (ch, msg) => {
  log(`Message received on channel ${ch}: ${msg}`)
  const point = {
    measurement: 'pid-values',
    tags: { name: ch },
    fields: { value: msg }
  }
  log(`About to write point: ${point}`)
  influxDB.writePoints([point])
    .then( () => log('Data successfully written to InfluxDB') )
    .catch( err => log('Error caught storing data in InfluxDB', err) )
})

// this.addPoller("vss")
// this.addPoller("rpm")
// this.addPoller("temp")
// this.addPoller("load_pct")
// this.addPoller("map")

influxDB.getDatabaseNames()
  .then( names => {
    if (!names.includes(process.env.INFLUXDB_DATABASE)) {
      influxDB.createDatabase(process.env.INFLUXDB_DATABASE)
    }
  }).then( () => {
    debug('InfluxDB database ready')
    return subscribe( 'vss', 'rpm', 'temp', 'load_pct', 'map' )
  }).then( (result) => {
    debug('Redis SUBSCRIBE command returned: ', result)
  }).catch( err => {
    debug('Failed to start. Error: ', err)
    process.exit(1)
  })