package main

import (
  "encoding/json"
  "github.com/ikester/blinkt"
  "github.com/lucasb-eyer/go-colorful"
  "io/ioutil"
  "log"
  "net/http"
  "time"
)

const numPixels = 8
const cheerLightURL = "http://api.thingspeak.com/channels/1417/field/2/last.json"

type cheerlight struct {
  Colour string `json:"field2"`
}

func main() {
  log.Println("Cheerlights started")

  theBlinkt := blinkt.NewBlinkt()
  theBlinkt.ShowAnimOnStart = true
  theBlinkt.CaptureExit = true
  theBlinkt.ShowAnimOnExit = true
  theBlinkt.ClearOnExit = true
//  theBlinkt.SetBrightness(0.5)

  theBlinkt.Setup()

  var prevColour, getErr = colorful.Hex("#000000")
  if getErr != nil {
    log.Println(getErr)
  }
  newColour := colorful.FastWarmColor()
  //newColour := colorful.FastHappyColour()
  prevColour = Blend(theBlinkt, prevColour, newColour)
  time.Sleep(30 * time.Second)

  //loop forever
  for {
    var c, getErr = colorful.Hex(GetCheerlightColours())
    if getErr != nil {
      log.Println(getErr)
      c = colorful.HappyColor()
    }
    prevColour = Blend(theBlinkt, prevColour, c)

    //wait before checking for updated colour value
    time.Sleep(10 * time.Minute)
  }
}

func GetCheerlightColours() string {
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

/*
 Set all pixels to the colour c1
 */
func SetAll(theBlinkt blinkt.Blinkt, c colorful.Color) {
  var r, g, b = c.RGB255()
  log.Printf("Setting pixels to %v, %v, %v", r, g, b)
  theBlinkt.SetAll(int(r), int(g), int(b))
  theBlinkt.Show()
}

/*
  Blend from colour c1 to colour c2
 */
func Blend(theBlink blinkt.Blinkt, fromColour colorful.Color, toColour colorful.Color) colorful.Color {
  log.Printf("Blending from %s to %s", fromColour.Hex(), toColour.Hex())
  steps := 25
  for i := 0 ; i < steps; i++ {
    opColour := fromColour.BlendHsv(toColour, float64(i)/float64(steps - 1))
    SetAll(theBlink, opColour)
    time.Sleep(250 * time.Millisecond)
  }
  return toColour
}

