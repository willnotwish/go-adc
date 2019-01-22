package main

import (
  "flag"
  "time"
  "log"
  "os"
  "github.com/schollz/progressbar"
  "github.com/willnotwish/go-adc/internal"
)

func main() {

  log.SetPrefix("go-adc: ")

  var count int
  var outputFilename string
  var vref float64
  var showProgress bool

  interval, err := time.ParseDuration("100ms")
  if err != nil {
    panic(err)
  }

  flag.IntVar(&count, "count", 1000, "number of samples to take")
  flag.DurationVar(&interval, "interval", interval, "interval between samples")
  flag.StringVar(&outputFilename, "output", "/results/output.txt", "output file name")
  flag.Float64Var(&vref, "vref", 3.3, "reference voltage applied to MCP3008")
  flag.BoolVar(&showProgress, "progress", true, "show progress bar")

  flag.Parse()

  ofile, err := os.Create(outputFilename)
  if err != nil {
    panic(err)
  }
  defer ofile.Close()

  var bar *progressbar.ProgressBar = nil
  if showProgress {
    bar = progressbar.New(count)
  }

  internal.Sample(count, interval, internal.Voltage(vref), ofile, bar)
}

// func httpHandler(w http.ResponseWriter, r *http.Request) {
//   fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }

// func mainx() {
//   http.HandleFunc("/", httpHandler)
//   log.Fatal(http.ListenAndServe(":8080", nil))
// }