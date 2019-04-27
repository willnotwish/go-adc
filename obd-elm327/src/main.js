const Debug = require('debug')
const OBDReader = require('serial-obd')
const Influx = require('influx')

// const Pusher = require('pusher')
// const redis = require('redis')

const debug = Debug('obd:main')

// const redisClient = redis.createClient(6379, 'redis') // port & host. See docker-compose.yml

const redis = false
const pusher = false

const influxDBClient = new Influx.InfluxDB({
  host: process.env.INFLUXDB_HOST,
  port: process.env.INFLUXDB_PORT,
  database: process.env.INFLUXDB_DATABASE,
  schema: [{
    measurement: 'pid-values',
    fields: {
      value: Influx.FieldType.FLOAT
    },
    tags: ['pid']
  }]
})

// Creates the named database if it doesn't already exist
// Returns a promise which resolves with an influx DB client
setupInfluxDBClient = database => {
  return influxDBClient.getDatabaseNames()
    .then( names => {
      if (!names.includes(database)) {
        return influxDBClient.createDatabase(database)
      }
    }).then( () => {
      debug("Database ready")
      return influxDBClient
    })
}

// redisClient.on( 'connect', () => console.log('Redis client connected') )
// redisClient.on( 'error', (err) => console.log('Redis cilent. Something went wrong ', err) )

// const log = (msg, ...params) => console.log("main. " + msg, ...params)

// const pusher = new Pusher({
//   appId:    process.env.PUSHER_APP_ID,
//   key:      process.env.PUSHER_APP_KEY,
//   secret:   process.env.PUSHER_APP_SECRET,
//   cluster:  process.env.PUSHER_APP_CLUSTER,
//   encryted: true
// })
// log("Pusher client created")

const getClient = service => {
  switch (service) {
    case 'InfluxDB':
      return setupInfluxDBClient(process.env.INFLUXDB_DATABASE)
      break
    case 'Redis':
      return Promise.resolve(redis)
      break
    case 'Pusher':
      return Promise.resolve(pusher)
      break
    default:
      throw new Error("Unsupported service", service)
      break
  }
}

const storeData = (service, writeData) => {
  return getClient(service).then( client => {
    if (client) {
      return writeData(client)
    } else {
      debug(`Storage service ${service} is not available. This might be a config issue.`)
      return false
    }
  })
}

const options = {
  baudRate: 38400
}

debug("About to create OBDReader with: ", process.env.DEVICE)
const serialOBDReader = new OBDReader(process.env.DEVICE, options)

let dataReceivedMarker = {}
serialOBDReader.on('dataReceived', data => {
  debug(data)
  dataReceivedMarker = data

  // For now, a simple sanity check. Need to investigate this further
  if (data.pid) {

    // storeData( 'Redis', client => client.xadd('obd-stream', '*', data.pid, data.value, redis.print) )

    setupInfluxDBClient(process.env.INFLUXDB_DATABASE).then( client => {
    // storeData( 'InfluxDB', client => {
      const point = {
        measurement: 'pid-values',
        tags: { pid: data.pid },
        fields: { value: data.value }
      }
      return client.writePoints([point])
    }).then( () => debug('Data successfully written to InfluxDB') )
      .catch( err => debug('Error caught storing data in InfluxDB', err) )

    // storeData( 'Pusher', client => pusher.trigger('obd', 'polled-pid-values', data) )
  }
})

serialOBDReader.on('connected', function(data) {
  debug("Connected to ELM327. Adding pollers...")
  this.addPoller("vss")
  this.addPoller("rpm")
  this.addPoller("temp")
  this.addPoller("load_pct")
  this.addPoller("map")
  // this.addPoller("frp")

  this.startPolling(2000)
})

debug("About to connect to vehicle via serial port (ELM327)")
serialOBDReader.connect()