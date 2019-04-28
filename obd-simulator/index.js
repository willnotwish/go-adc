import {publish} from './redis/client'
import Debug from 'debug'

const debug = Debug('obd-simulator:index')

debug("Starting up...")

let speed = 0
const publishSpeed = () => {

  debug("Speed: ", speed)

  publish('vss', speed)
    .then( result => debug("Published data. Result: ", result) )
    .catch( err => debug("Failed to publish data. Error: ", err) )

  if (speed > 80) {
    speed = 30
  } else {
    speed += 10
  }
}

setInterval(publishSpeed, 2000)

debug("Started")