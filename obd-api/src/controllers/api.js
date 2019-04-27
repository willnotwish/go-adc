import express from 'express'
import store from '../obd/influxdb/store'
import serializer from '../obd/serializer'

import Debug from 'debug'

const debug = Debug('obd-api:controllers-api')

const router = express.Router()

router.get('/data', (req, res, next) => {
  debug('TODO. configure options from http request' )
  const options = {}
  store.getData(options).then( data => {
    debug('data received from influxdb: ', data)
    const serialized = req.query.fmt === 'raw' ? data : serializer.serialize(data)
    res.status(200).send(serialized)
  }).catch( next )
})

export default router