import {InfluxDB} from 'influx'

export default new InfluxDB({
  host: process.env.INFLUXDB_HOST,
  port: process.env.INFLUXDB_PORT,
  database: process.env.INFLUXDB_DATABASE
})
