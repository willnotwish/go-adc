const SerialPort = require('serialport')

const sp = new SerialPort(process.env.DEVICE)

sp.on( 'error', (err) => {
  console.log("Error: ", err)
})

sp.on( 'open', console.log )

sp.write("AAABBBCCC\n")
sp.write("from Nick\n")

sp.on( 'data', data => {
  console.log("Data received: ", data)
})
