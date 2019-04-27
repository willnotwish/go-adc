import express from 'express'
import logger from 'morgan'

import api from './controllers/api'

const app = express()

app.use(logger('dev'))
app.use('/api/v1', api)

export default app
