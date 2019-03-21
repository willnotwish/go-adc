package main

import (
  "errors"
  "context"
  "flag"
  "time"
  "log"
  "sort"
  "strconv"
  "os"
  "os/signal"
  "syscall"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/urfave/negroni"
  // "github.com/auth0/go-jwt-middleware"
  // "github.com/dgrijalva/jwt-go"
  dl "github.com/willnotwish/go-adc/internal/datalogger"
)

// type Jwks struct {
//   Keys []JSONWebKeys `json:"keys"`
// }

// type JSONWebKeys struct {
//   Kty string `json:"kty"`
//   Kid string `json:"kid"`
//   Use string `json:"use"`
//   N string `json:"n"`
//   E string `json:"e"`
//   X5c []string `json:"x5c"`
// }

func main() {
  log.SetPrefix("go-adc: ")

  var vref float64

  interval, err := time.ParseDuration("100ms")
  if err != nil {
    panic(err)
  }

  flag.DurationVar(&interval, "interval", interval, "clock interval between samples")
  flag.Float64Var(&vref, "vref", 3.3, "reference voltage applied to MCP3008")
  flag.Parse()

  output := make( chan dl.Data, 100 ) // buffered results channel
  stop := make( chan bool )

  dl.Start(interval, dl.Volts(vref), output)
  dl.Enable(0, 10)

  log.Println("About to start WS hub. Using Gorilla ws implementation")
  hub := newHub()
  go hub.run()

  go func() {
    buffer := dl.NewMovingWindow(1000, 2)

    for {
      select {
      case sample := <- output:
        log.Printf( "Storing %5.2fV (input %d)\n", sample.Voltage, sample.Input)
        buffer.PushBack( sample )
        js, err := json.Marshal( sample )
        if err != nil {
          panic(err)
        }
        hub.broadcast <- append(js, 0x0a)
      case req := <- slicer:
        log.Printf("Request received for sampled data since: %s\n", req.since)
        slice := buffer.Slice()
        index := sort.Search( len(slice), func( i int ) bool { return slice[i].Timestamp.After( req.since ) } )
        js, err := json.Marshal( slice[index:] )   // from 'index' to the end
        if err != nil {
          panic(err)
        }
        req.ch <- js
      case <- stop:
        log.Println("stop received: returning from results writer go routine")
        return
      }
    }
  }()

  // jwtValidator := jwtmiddleware.New(jwtmiddleware.Options {
  //   ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
  //     // Verify 'aud' claim
  //     aud := os.Getenv("AUTH0_AUDIENCE")
  //     checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
  //     if !checkAud {
  //       return token, errors.New("Invalid audience.")
  //     }
  //     // Verify 'iss' claim
  //     iss := "https://" + os.Getenv("AUTH0_DOMAIN") + "/"
  //     checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
  //     if !checkIss {
  //       return token, errors.New("Invalid issuer.")
  //     }

  //     cert, err := getPemCert(token)
  //     if err != nil {
  //       panic(err.Error())
  //     }

  //     result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
  //     return result, nil
  //   },
  //   SigningMethod: jwt.SigningMethodRS256,
  // })


  log.Println("About to start HTTP server. Using Gorilla mux router")

  // Use Gorilla mux router with negroni classic middleware + jwt validation
  router := mux.NewRouter()

  // nx := negroni.Classic()
  // nx.Use(negroni.HandlerFunc(jwtValidator.HandlerWithNext))

  // router.HandleFunc("/api/datax", dataHandler).Methods("GET")

  router.HandleFunc("/api/data",    dataHandler   ).Methods("GET")
  router.HandleFunc("/api/control", controlHandler).Methods("POST", "PUT")
  router.HandleFunc("/api/config",  configHandler ).Methods("GET")

  router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    log.Println("WS request received")
    serveWebsocket(hub, w, r)
  })

  n := negroni.Classic()
  // n.Use(negroni.HandlerFunc(jwtValidator.HandlerWithNext))
  n.UseHandler(router)

  srv := &http.Server{ Addr: ":8001", Handler: n }

  // Use a separate goroutine to handle graceful shutdown,
  // calling the new Shutdown func.
  idleConnsClosed := make(chan struct{})

  go func() {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    sig := <-sigCh
    log.Printf("Signal %s received from sigCh. About to call srv.Shutdown()\n", sig)

    if err := srv.Shutdown(context.Background()); err != nil {
      // Error from closing listeners, or context timeout:
      log.Printf("HTTP server Shutdown: %v", err)
    }
    close(idleConnsClosed)
  }()

  if err := srv.ListenAndServe(); err != http.ErrServerClosed {
    log.Printf("HTTP server ListenAndServe: %v", err)
  }

  log.Println("Waiting on idleConnsClosed")
  <- idleConnsClosed

  log.Println("HTTP server shut down. About to stop hub")
  hub.shutdown <- true

  log.Println("HTTP server shut down. About to stop data logger")
  dl.Stop()

  stop <- true
}

type SliceRequest struct {
  ch chan []byte
  since time.Time
}

var slicer chan SliceRequest = make(chan SliceRequest)

var dataHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  since, err := parseDataParams(r)
  if err != nil {
    http.Error(w, err.Error(), 422)
    return
  }
  ch := make(chan []byte) // a channel suitable for a byte slice

  // Pass the channel along with the since parameter to the go routine storing the samples (the slicer)
  slicer <- SliceRequest{ch: ch, since: since}

  // read the json the slicer generates; use as the http response
  w.Write( <-ch )
})

func configHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  cfg := dl.GetConfig()
  js, err := json.Marshal( cfg )
  if err != nil {
    panic(err)
  }
  w.Write( js )
}

func controlHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  input, divisor, err := parseControlParams(r)
  if err != nil {
    http.Error(w, err.Error(), 422)
    return
  }

  if divisor > 0 {
    dl.Enable(input, divisor)
  } else {
    dl.Disable(input)
  }
}

func parseControlParams(r *http.Request) (int, int, error) {
  q := r.URL.Query()

  input, err := strconv.Atoi(q.Get("i"))
  if err != nil {
    return -1, -1, err
  }

  divisor, err := strconv.Atoi(q.Get("d"))
  if err != nil {
    return input, -1, err
  }

  if input < 0 || input > 7 {
    return input, divisor, errors.New("Input must be between 0 and 7")
  }

  return input, divisor, nil
}

func parseDataParams(r *http.Request) (time.Time, error) {
  q := r.URL.Query()
  _, present := q["s"]
  if present != true {
    return time.Time{}, nil
  }

  return time.Parse(time.RFC3339Nano, q.Get("s"))
}

// func getPemCert(token *jwt.Token) (string, error) {
//   cert := ""
//   resp, err := http.Get("https://" + os.Getenv("AUTH0_DOMAIN") + "/.well-known/jwks.json")
//   if err != nil {
//     return cert, err
//   }
//   defer resp.Body.Close()

//   var jwks = Jwks{}
//   err = json.NewDecoder(resp.Body).Decode(&jwks)
//   if err != nil {
//     return cert, err
//   }

//   for k, _ := range jwks.Keys {
//     if token.Header["kid"] == jwks.Keys[k].Kid {
//       cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
//     }
//   }

//   if cert == "" {
//     err := errors.New("Unable to find appropriate key.")
//     return cert, err
//   }

//   return cert, nil
// }
