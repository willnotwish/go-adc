package main

import (
  // "html/template"
  "net/http"
  // "os"
  // "fmt"
  // "time"
  // "github.com/gorilla/mux"
  // "github.com/willnotwish/go-adc/internal/mcp3008"
)

// var indexTemplate *template.Template

type ViewData struct {
  Title string
  // Samples []mcp3008.VoltageSample
}

func main() {
  // var err error
  // indexTemplate, err = template.ParseFiles("web/index.gohtml")
  // if err != nil {
  //   panic(err)
  // }

  http.HandleFunc("/", home)
  http.ListenAndServe(":3000", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")

  // interval, err := time.ParseDuration("100ms")
  // if err != nil {
  //   panic(err)
  // }

  // var count int = 100

  // results := make([]mcp3008.VoltageSample, count)
  w.Header().Set("Content-Type", "application/json")
  // mcp3008.StreamSamplesAsJson(count, interval, 3.3, w)

  // err = indexTemplate.Execute(w, ViewData{
  //   Title: "Voltage Samples",
  //   Samples: results,
  // })

  // if err != nil {
  //   http.Error(w, err.Error(), http.StatusInternalServerError)
  //   return
  // }
}