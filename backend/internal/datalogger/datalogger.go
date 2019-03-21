package datalogger

import (
  "fmt"
  "log"
  "time"
  "github.com/willnotwish/go-adc/internal/mcp3008"
)

type Volts float64

func (v Volts) MarshalJSON() ([]byte, error) {
  return []byte(fmt.Sprintf("%5.2f", v)), nil
}

type Data struct {
  Input int
  Timestamp time.Time
  Voltage Volts
}

type Command struct {
  Shutdown bool
  Input int
  Threshold int
}

type Config struct {
  Interval int64 // in ms
  Divisors []int
}

type Request struct {
  response chan Config
}

var cmdCh = make(chan Command)
var reqCh = make(chan Request)

func Start( interval time.Duration, vref Volts, dataChannel chan Data ) {
  log.Printf("Start. Interval: %s\n", interval)

  ticker := time.NewTicker(interval)

  go func() {
    mcp3008.SpiSetup()
    defer mcp3008.SpiTeardown()

    thresholds := make([]int, 8)
    counters := make([]int, 8)
    for {
      select {
      case <-ticker.C:
        // When an input's counter reaches its threshold (if non zero),
        // that input's voltage is logged and the counter reset. Otherwise,
        // the counter is incremented.
        for input, counter := range counters {
          if thresholds[input] > 0 {
            if counter >= thresholds[input] {
              dataChannel <- Data{ input, time.Now(), Volts(mcp3008.SpiSample(input)) / 1023 * vref }
              counters[input] = 0
            } else {
              counters[input] += 1
            }
          } else {
            counters[input] = 0
          }
        }
      case cmd := <-cmdCh:
        if cmd.Shutdown {
          log.Println("Shutdown command received. Go routine returning...")
          return
        } else {
          thresholds[cmd.Input] = cmd.Threshold
        }
      case req := <-reqCh:
        log.Println("config request received")
        req.response <- Config{ Interval: interval.Nanoseconds()/1000000, Divisors: thresholds }
      }
    }
  }()
}

func Stop() {
  log.Println("Sending Shutdown command")
  cmdCh <- Command{ Shutdown: true }
}

func Enable(input int, threshold int) {
  log.Printf("EnableInput %d with threshold %d\n", input, threshold)
  cmdCh <- Command{ Input: input, Threshold: threshold }
}

func Disable(input int) {
  log.Printf("Disable input %d\n", input)
  cmdCh <- Command{ Input: input, Threshold: 0 }
}

func GetConfig() Config {
  log.Println("Sending GetConfig request")
  ch := make( chan Config )
  reqCh <- Request{ response: ch }
  return <- ch
}