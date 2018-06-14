package main

import (
  "fmt"
  "time"
  "github.com/lucasb-eyer/go-colorful"
  "github.com/ikester/blinkt"
  "io/ioutil"
  "log"
  "encoding/json"
  "net/http"
)

const numPixels = 8
const cheerLightURL = "http://api.thingspeak.com/channels/1417/field/2/last.json"

type cheerlight struct {
  Colour string `json:"field2"`
}

func main() {
  fmt.Println("Cheerlights started")

  theBlinkt := blinkt.NewBlinkt()
  theBlinkt.ShowAnimOnStart = true
  theBlinkt.CaptureExit = true
  theBlinkt.ShowAnimOnExit = true
  theBlinkt.ClearOnExit = true

  theBlinkt.Setup()
  time.Sleep(30 * time.Second)

  //loop forever
  for {

    var c, getErr = colorful.Hex(getCheerlightColours())
    if getErr != nil {
      log.Println(getErr)
    }
    var r, g, b = c.RGB255()

    //theBlinkt.SetPixel(i, int(c.R), int(c.G), int(c.B))
    theBlinkt.SetAll(int(r), int(g), int(b))
    theBlinkt.Show()

    log.Printf("Setting pixels to %v, %v, %v", r, g, b)
    time.Sleep(100 * time.Second)
  }
}

func getCheerlightColours() (string) {
  var netClient = &http.Client{
    Timeout: time.Second * 3,
  }
  resp, getErr := netClient.Get(cheerLightURL)

  if getErr != nil {
    log.Panic(getErr)
  }

  if resp.Body != nil {
    defer resp.Body.Close()
  }

  body, readErr := ioutil.ReadAll(resp.Body)

  if readErr != nil {
    log.Panic(getErr)
  }

  result := cheerlight{}

  parseErr := json.Unmarshal(body, &result)
  if parseErr != nil {
    log.Panic("Can't parse response")
    return "#0000"
  }
  return result.Colour
}

//func check(err error) {
//    if err != nil {
//        fmt.Println(err)
//        os.Exit(1)
//    }
//}
