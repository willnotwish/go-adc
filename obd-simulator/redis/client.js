import Redis from 'redis'
import {promisify} from 'util'
import Debug from 'debug'

const debug = Debug('obd-simulator:redis-client')

const redis = Redis.createClient(process.env.REDIS_PORT || 6379, process.env.REDIS_HOST) // port & host. See docker-compose.yml

redis.on( 'connect', () => debug('Connected.') )
redis.on( 'error', (err) => debug('Something went wrong: ', err) )
redis.on( 'ready', () => debug('Ready') )

const publish = promisify(redis.publish).bind(redis)

export {redis, publish}