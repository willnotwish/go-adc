package main

import (
	"github.com/stianeikeland/go-rpio"
	"github.com/schollz/progressbar"
	"flag"
	"fmt"
	"time"
	"log"
	"os"
)

func main() {

	log.SetPrefix("go-adc: ")

	// CLI setup
  var sampleCount int
  var outputFilename string

  interval, err := time.ParseDuration("100ms")
  if err != nil {
  	panic(err)
  }

  flag.IntVar(&sampleCount, "count", 1000, "number of samples to take")
  flag.DurationVar(&interval, "interval", interval, "interval between samples")
  flag.StringVar(&outputFilename, "output", "/results/output.txt", "output file name")

  flag.Parse()

	setupSpi()
	defer teardownSpi()

  outputFile, err := os.Create(outputFilename)
  if err != nil {
  	panic(err)
  }
  defer outputFile.Close()

	record(sampleCount, outputFile, interval)
}

// It's a 10-bit sample: scale it appropriately; 0-1023 where 1023 equals Vref
func record(count int, file *os.File, interval time.Duration) {
	log.Printf("About to collect %d samples, at intervals of %s", count, interval.String())

	const Vref = 3.3
	log.Printf("Vref: %5.2f", Vref)

	bar := progressbar.New(count)
	for i := 1; i < count; i++ {
		voltage := float32(takeMCP3008Sample()) / 1023 * Vref
		fmt.Fprintf(file, "%s %5.2f\n", time.Now().Format("15:04:05.000"), voltage)
    bar.Add(1)
		time.Sleep(interval)
	}
	log.Printf("%d samples collected", count)
}

func setupSpi() {
	log.Print("setupSpi()...")

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

func teardownSpi() {
	rpio.SpiEnd(rpio.Spi0)
	log.Print("Called rpio.SpiEnd")

	rpio.Close()
	log.Print("Finished.")
}

//
//
//
func takeMCP3008Sample() uint16 {
	// Best explanation I have read is here
	// http://www.hertaville.com/interfacing-an-spi-adc-mcp3008-chip-to-the-raspberry-pi-using-c.html

	// The first byte starts it off. It is 7 x 0s, with a deciding final 1.
	// We don't care about the byte we get back in return.
	//
	// The next byte is SGL/DIFF D2 D1 D0 in the first four bits, and don't care for the remainder
	// Single-ended conversion has SGL=1,
	// and we're sampling CH0, so we set D2 = D1 = D0 = 0
	// In binary then, we send 10000000 (with 0's in the con't care bits).
	// In hex, that's 0x80.

	// Not so fast, though. This time we're interested in the byte we receive in return:
	// it contains (as its last two bits), the two most significant bits of the 10-bit sample.
	// The next byte we get back contains the remaining 8 bits of the sample.

	// Use the exchange function, sending 0's to get the third byte, like this
	buffer := []byte{ 0x01, 0x80, 0x00 }
	rpio.SpiExchange(buffer)

	// Do the maths to get the sample 10-bit value from the last two bytes
	// Mask out all the bits except the last two in the middle byte
	return uint16(buffer[1] & 0x03)<<8 | uint16(buffer[2])
}

// func httpHandler(w http.ResponseWriter, r *http.Request) {
//   fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }

// func mainx() {
//   http.HandleFunc("/", httpHandler)
//   log.Fatal(http.ListenAndServe(":8080", nil))
// }