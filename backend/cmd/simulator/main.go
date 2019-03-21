package main

import (
  // "flag"
  "time"
  // "log"
  "fmt"
  // "os"
  // "github.com/schollz/progressbar"
  // "github.com/willnotwish/go-adc/internal"
)

// func main() {

//   log.SetPrefix("go-adc: ")

//   var count int
//   var outputFilename string
//   var vref float64
//   var showProgress bool

//   interval, err := time.ParseDuration("100ms")
//   if err != nil {
//     panic(err)
//   }

//   flag.IntVar(&count, "count", 1000, "number of samples to take")
//   flag.DurationVar(&interval, "interval", interval, "interval between samples")
//   flag.StringVar(&outputFilename, "output", "/results/output.txt", "output file name")
//   flag.Float64Var(&vref, "vref", 3.3, "reference voltage applied to MCP3008")
//   flag.BoolVar(&showProgress, "progress", true, "show progress bar")

//   flag.Parse()

//   // ofile, err := os.Create(outputFilename)
//   // if err != nil {
//   //   panic(err)
//   // }
//   // defer ofile.Close()

//   // var bar *progressbar.ProgressBar = nil
//   // if showProgress {
//   //   bar = progressbar.New(count)
//   // }

//   // mcp3008.StreamSamples(count, interval, mcp3008.Volts(vref), ofile, bar)

//   fmt.Println("Finished. Gute Arbeit!")
// }

type Command struct {
  Count int
  Stop  bool
}

func periodic(ticker *time.Ticker, command <-chan Command, count0 int) {
  count := count0
  index := 0
  for {
    select {
    case cmd := <-command:
      if cmd.Stop {
        fmt.Println("periodic. shutdown request received.")
        return
      } else if cmd.Count > 0 {
        count = cmd.Count
      } else {
        fmt.Printf("periodic. unexpected command received: %s\n", cmd)
      }
    case <-ticker.C:
      if index >= count {
        fmt.Print(".")
        index = 0
      } else {
        index = index+1
      }
    }
  }
}

func main() {
  ticker := time.NewTicker(10 * time.Millisecond)

  //  shutdown := make(chan bool)
  cmd := make(chan Command)

  fmt.Println("main. starting periodic")
  go periodic(ticker, cmd, 50)

  fmt.Println("main. sleeping for 5s")
  time.Sleep(5 * time.Second)

  fmt.Println("main. sending speed up command")
  cmd <- Command{Count: 5}

  fmt.Println("main. sleeping for 5s")
  time.Sleep(5 * time.Second)

  fmt.Println("main. sending speed up command")
  cmd <- Command{Count: 3}

  fmt.Println("main. sleeping for 5s")
  time.Sleep(5 * time.Second)

  ticker.Stop()
  fmt.Println("main. ticker stopped. sleeping for 5s, during which time there should be no ticks...")
  time.Sleep(5 * time.Second)
  cmd <- Command{Stop: true}
  fmt.Println("main. shutdown command sent. sleeping for 5s")
  time.Sleep(5 * time.Second)
  fmt.Println("main. finished.")
}
