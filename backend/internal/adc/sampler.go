package adc

import (
  "log"
  "time"
  "github.com/willnotwish/go-adc/internal/mcp3008"
)

// Exports

type Sample struct {
  Input int
  Timestamp time.Time
  Value uint16 // 0-1023
}

type Config struct {
  Interval int64 // in ms
  Divisors []int
}

type Sampler struct {
  // The base interval between samples
  interval time.Duration

  // Where samples are written to
  Out chan Sample

  // Incoming requests
  configRequests chan chan Config
  thresholdUpdates chan thresholdUpdate

  shutdown chan bool
}

func NewSampler(interval time.Duration, bufferSize int) *Sampler {
  return &Sampler{
    interval:         interval,
    Out:              make(chan Sample, bufferSize),
    configRequests:   make(chan chan Config),
    thresholdUpdates: make(chan thresholdUpdate),
    shutdown:         make(chan bool),
  }
}

func (s *Sampler) Run() {
  ticker := time.NewTicker(s.interval)

  go func() {
    mcp3008.SpiSetup()
    defer mcp3008.SpiTeardown()

    thresholds := make([]int, 8)
    counters := make([]int, 8)
    for {
      select {
      case <-ticker.C:
        for input, counter := range counters {
          // When an input's counter reaches its threshold (if non zero),
          // that input's voltage is logged and the counter reset. Otherwise,
          // the counter is incremented.
          if thresholds[input] > 0 {
            if counter >= thresholds[input] {
              s.Out <- Sample{ input, time.Now(), mcp3008.SpiSample(input) }
              counters[input] = 0
            } else {
              counters[input] += 1
            }
          } else {
            counters[input] = 0
          }
        }

      case resp := <- s.configRequests:
        log.Println("Sampler. config request received")
        resp <- Config{ s.interval.Nanoseconds()/1000000, thresholds }

      case update := <- s.thresholdUpdates:
        log.Printf("Threshold update received: %v\n", update)
        thresholds[update.input] = update.threshold

      case <- s.shutdown:
        log.Println("Shutting down by returning from go routine")
        return
      }
    }
  }()
}

func (s *Sampler) Stop() {
  s.shutdown <- true
}

func (s *Sampler) Enable(input int, threshold int) {
  s.thresholdUpdates <- thresholdUpdate{ input, threshold }
}

func (s *Sampler) Disable(input int) {
  s.thresholdUpdates <- thresholdUpdate{ input, 0 }
}

func (s *Sampler) GetConfig() Config {
  ch := make( chan Config )
  s.configRequests <- ch
  return <- ch
}

// Private

type thresholdUpdate struct {
  input int
  threshold int
}