const SerialPort = require('serialport')

console.log("About to create serial port from device ", process.env.DEVICE)
const sp = new SerialPort(process.env.DEVICE)

sp.on( 'error', (err) => {
  console.log("Error: ", err)
})

sp.on( 'open', () => console.log("Serial port opened OK") )

const msg = "AAABBBCCCDDDEEEFFF\n"

sp.write(msg)
console.log("Sent: ", msg)

sp.write(msg)
console.log("Sent: ", msg)

sp.write(msg)
console.log("Sent: ", msg)

sp.write(msg)
console.log("Sent: ", msg)

sp.write(msg)
console.log("Sent: ", msg)

sp.on( 'data', data => {
  console.log("Data received: ", data)
})
