package main

import (
  "errors"
  "context"
  "flag"
  "fmt"
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
  "github.com/auth0/go-jwt-middleware"
  "github.com/dgrijalva/jwt-go"
  "github.com/willnotwish/go-adc/internal/adc"
)

type Volts float64

func (v Volts) MarshalJSON() ([]byte, error) {
  return []byte(fmt.Sprintf("%5.2f", v)), nil
}

type SampledVoltage struct {
  Input int
  Timestamp time.Time
  Voltage Volts
}

type SliceRequest struct {
  ch chan []byte
  since time.Time
}

type Jwks struct {
  Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
  Kty string `json:"kty"`
  Kid string `json:"kid"`
  Use string `json:"use"`
  N string `json:"n"`
  E string `json:"e"`
  X5c []string `json:"x5c"`
}

func main() {
  log.SetPrefix("go-adc: ")

  var vref float64
  vref = 3.3

  var scale float64
  scale = vref

  interval, err := time.ParseDuration("100ms")
  if err != nil {
    panic(err)
  }

  flag.DurationVar(&interval, "interval", interval, "clock interval between samples")
  flag.Float64Var(&scale, "scale", vref, "scaling factor to be applied to raw samples")
  flag.Parse()


  sampler := adc.NewSampler(interval, 100)
  sampler.Run()
  sampler.Enable(0, 10)

  log.Println("About to start WS hub (Gorilla implementation)")
  hub := newHub()
  go hub.run()

  // warehouse := newWarehouse( 1000, sampler.Out, scalingFactors, notifier )
  // warehouse.Run()

  stopCh := make( chan bool )
  slicerCh := make( chan SliceRequest )

  go func() {
    buffer := NewSlidingWindow(1000, 2)

    for {
      select {
      case sample := <- sampler.Out:
        v := float64(sample.Value) * scale / 1023
        log.Printf( "Voltage on channel %d: %5.2f\n", sample.Input, v )
        vs := SampledVoltage{ sample.Input, sample.Timestamp, Volts(v) }
        buffer.PushBack( vs )
        js, err := json.Marshal( vs )
        if err != nil {
          panic(err)
        }
        hub.broadcast <- append(js, 0x0a) // add newline
      case req := <- slicerCh:
        log.Printf("Slice request received for sampled data since: %s\n", req.since)
        slice := buffer.Slice()
        index := sort.Search( len(slice), func( i int ) bool { return slice[i].Timestamp.After( req.since ) } )
        js, err := json.Marshal( slice[index:] )   // from 'index' to the end
        if err != nil {
          panic(err)
        }
        req.ch <- js
      case <- stopCh:
        log.Println("stop received: returning from results writer go routine")
        return
      }
    }
  }()

  jwtValidator := jwtmiddleware.New(jwtmiddleware.Options {
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
      // Verify 'aud' claim
      aud := os.Getenv("AUTH0_AUDIENCE")
      checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
      if !checkAud {
        return token, errors.New("Invalid audience.")
      }
      // Verify 'iss' claim
      iss := "https://" + os.Getenv("AUTH0_DOMAIN") + "/"
      checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
      if !checkIss {
        return token, errors.New("Invalid issuer.")
      }

      cert, err := getPemCert(token)
      if err != nil {
        panic(err.Error())
      }

      result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
      return result, nil
    },
    SigningMethod: jwt.SigningMethodRS256,
  })

  log.Println("About to start HTTP server. Using Gorilla mux router")

  // Use Gorilla mux router with negroni classic middleware + jwt validation

  _data := func(w http.ResponseWriter, r *http.Request) {
    dataHandler(w, r, slicerCh)
  }

  _control := func(w http.ResponseWriter, r *http.Request) {
    controlHandler(w, r, sampler)
  }

  _config := func(w http.ResponseWriter, r *http.Request) {
    configHandler(w, r, sampler)
  }

  _ws := func(w http.ResponseWriter, r *http.Request) {
    log.Println("WS request received")
    serveWebsocket(hub, w, r)
  }

  // Routing

  // Public
  router := mux.NewRouter()
  router.HandleFunc("/data", _data).Methods("GET") // for testing
  router.HandleFunc("/ws",   _ws  ).Methods("GET")

  // API - for authorised clients only
  apiRoutes := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(true)

  apiRoutes.HandleFunc("/data",    _data   ).Methods("GET")
  apiRoutes.HandleFunc("/control", _control).Methods("POST", "PUT")
  apiRoutes.HandleFunc("/config",  _config ).Methods("GET")

  classic := negroni.Classic()

  classicWithJwtValidation := classic.With(
    negroni.HandlerFunc(jwtValidator.HandlerWithNext),
    negroni.Wrap(apiRoutes))

  router.PathPrefix("/api").Handler(classicWithJwtValidation)

  // Pay attention! This:
  //
  // classic.UseHandler(router)
  //
  // is the same as

  classic.Use(negroni.Wrap(router))

  // because
  //
  // func (n *Negroni) UseHandler(handler http.Handler) {
  //   n.Use(Wrap(handler))
  // }
  //

  srv := &http.Server{ Addr: ":8001", Handler: classic }

  // Use a separate goroutine to handle graceful shutdown,
  // calling the new Shutdown func.
  idleConnsClosed := make( chan struct{} )

  go func() {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    sig := <-sigCh
    log.Printf("Signal %s received. About to call srv.Shutdown()\n", sig)

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
  hub.stop()

  log.Println("HTTP server shut down. About to stop sampler")
  sampler.Stop()

  stopCh <- true
}



type DataHandler struct {
  slicer chan SliceRequest
}

func NewDataHandler(ch chan SliceRequest) *DataHandler {
  return &DataHandler{ slicer: ch }
}

func (h *DataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  dataHandler(w, r, h.slicer)
}

func dataHandler(w http.ResponseWriter, r *http.Request, slicer chan SliceRequest) {
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
}

func configHandler(w http.ResponseWriter, r *http.Request, sampler *adc.Sampler) {
  w.Header().Set("Content-Type", "application/json")

  cfg := sampler.GetConfig()
  js, err := json.Marshal( cfg )
  if err != nil {
    panic(err)
  }
  w.Write( js )
}

func controlHandler(w http.ResponseWriter, r *http.Request, sampler *adc.Sampler) {
  w.Header().Set("Content-Type", "application/json")

  input, divisor, err := parseControlParams(r)
  if err != nil {
    http.Error(w, err.Error(), 422)
    return
  }

  if divisor > 0 {
    sampler.Enable(input, divisor)
  } else {
    sampler.Disable(input)
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

func getPemCert(token *jwt.Token) (string, error) {
  cert := ""
  resp, err := http.Get("https://" + os.Getenv("AUTH0_DOMAIN") + "/.well-known/jwks.json")
  if err != nil {
    return cert, err
  }
  defer resp.Body.Close()

  var jwks = Jwks{}
  err = json.NewDecoder(resp.Body).Decode(&jwks)
  if err != nil {
    return cert, err
  }

  for k, _ := range jwks.Keys {
    if token.Header["kid"] == jwks.Keys[k].Kid {
      cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
    }
  }

  if cert == "" {
    err := errors.New("Unable to find appropriate key.")
    return cert, err
  }

  return cert, nil
}
