package mcp3008

import (
  "github.com/stianeikeland/go-rpio"
  "log"
)

func SpiSetup() {
  log.Print("setup()...")

  if err := rpio.Open(); err != nil {
    panic(err)
  }
  log.Print("Called rpio.Open OK")

  if err := rpio.SpiBegin(rpio.Spi0); err != nil {
    panic(err)
  }

  log.Print("Called rpio.SpiBegin OK")

  rpio.SpiChipSelect(0) // Select CE0 slave
  log.Print("Called rpio.SpiChipSelect OK")
}

func SpiTeardown() {
  log.Print("teardown()...")

  rpio.SpiEnd(rpio.Spi0)
  log.Print("Called rpio.SpiEnd")

  rpio.Close()
  log.Print("Finished.")
}

// Capture a single 10-bit sample on the given input (0-7), returned as 16-bit integer
func SpiSample( channel int ) uint16 {
  // Best explanation I have read is here
  // http://www.hertaville.com/interfacing-an-spi-adc-mcp3008-chip-to-the-raspberry-pi-using-c.html

  // The first byte starts it off. It is 7 x 0s, with a deciding final 1.
  // We don't care about the byte we get back in return.
  //
  // The next byte is SGL/DIFF D2 D1 D0 in the first four bits, and don't care for the remainder.
  // Single-ended conversion (required) has SGL=1, per this table:

  // 1 0 0 0 single-ended CH0
  // 1 0 0 1 single-ended CH1
  // 1 0 1 0 single-ended CH2
  // 1 0 1 1 single-ended CH3
  // 1 1 0 0 single-ended CH4
  // 1 1 0 1 single-ended CH5
  // 1 1 1 0 single-ended CH6
  // 1 1 1 1 single-ended CH7

  // So if we're sampling input 0 (CH0), we want D2 = D1 = D0 = 0
  // In binary then, we send 10000000 (with 0's in the con't care bits).
  // In hex, that's 0x80.

  // if channel < 0 || channel >= 8 {
  //   return 0xFFFF, errors.New("channel must be between 0 and 7")
  // }

  // We're interested in the byte we receive in return.
  // It contains - as its last two bits - the two most significant bits of the 10-bit sample.

  // The third byte received is the remaining 8 bits of the sample.

  // We need to send something (anything) to get the third byte

  buffer := []byte{ 0x01, (byte(channel) << 4) | 0x80, 0x00 }
  rpio.SpiExchange(buffer)

  // Do the maths to get the sample 10-bit value from the last two bytes
  // Mask out all the bits except the last two in the middle byte (buffer[1])
  return uint16(buffer[1] & 0x03)<<8 | uint16(buffer[2])
}
