const OBDReader = require('serial-obd');

const log = (msg, ...params) => console.log("main. " + msg, ...params)

const options = {}
options.baudRate = 38400

log("About to create OBDReader with: ", process.env.DEVICE)
const serialOBDReader = new OBDReader(process.env.DEVICE, options)

let dataReceivedMarker = {}

serialOBDReader.on('dataReceived', data => {
  console.log(data)
  dataReceivedMarker = data
})

serialOBDReader.on('connected', function(data) {
  log("Connected. Adding pollers...")
  this.addPoller("vss")
  this.addPoller("rpm")
  this.addPoller("temp")
  this.addPoller("load_pct")
  this.addPoller("map")
  this.addPoller("frp")

  this.startPolling(2000)
})

log("About to connect to vehicle")
serialOBDReader.connect()

// const SerialPort = require('serialport')

// const log = (msg, ...params) => {
//   console.log("main. " + msg, params)
// }

// log("About to create serial port from device ", process.env.DEVICE)
// const sp = new SerialPort(process.env.DEVICE)

// sp.on( 'error', (err) => {
//   log("Error: ", err)
// })

// sp.on( 'open', () => log("Serial port opened OK") )

// const msg = "AAABBBCCCDDDEEEFFF\n"

// sp.write(msg)
// log("Sent: ", msg)

// sp.write(msg)
// log("Sent: ", msg)

// sp.write(msg)
// log("Sent: ", msg)

// sp.write(msg)
// log("Sent: ", msg)

// sp.write(msg)
// log("Sent: ", msg)

// sp.on( 'data', data => {
//   log("Data received: ", data)
// })