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

const cheerLightURL = "http://api.thingspeak.com/channels/1417/field/2/last.json"

type cheerlight struct {
  Colour string `json:"field2"`
}

var lightsOn = false

func startLights(theBlinkt *blinkt.Blinkt, currentColour colorful.Color) colorful.Color {
  theBlinkt.ShowInitialAnim()
  newColour := colorful.FastWarmColor()
  return Blend(theBlinkt, currentColour, newColour)
}

func stopLights(theBlinkt *blinkt.Blinkt, currentColour colorful.Color) colorful.Color {
  newColour, _ := colorful.Hex("#000000")
  return Blend(theBlinkt, currentColour, newColour)
}

func updateLights(theBlinkt *blinkt.Blinkt, currentColour colorful.Color) colorful.Color {
    var c, getErr = colorful.Hex(GetCheerlightColours())
    if getErr != nil {
      log.Println(getErr)
      c = colorful.HappyColor()
    }
    return Blend(theBlinkt, currentColour, c)
}

func main() {
  log.Println("Cheerlights started")

  theBlinkt := blinkt.NewBlinkt()
  theBlinkt.ShowAnimOnStart = false
  theBlinkt.CaptureExit = true
  theBlinkt.ShowAnimOnExit = true
  theBlinkt.ClearOnExit = true
//  theBlinkt.SetBrightness(0.5)

  theBlinkt.Setup()

  var currentColour, getErr = colorful.Hex("#000000")
  if getErr != nil {
    log.Println(getErr)
  }

  startTime := time.Date(0, 0, 0, 19, 30, 0, 0, time.Local)
  stopTime := time.Date(0, 0, 0, 23, 59, 0, 0, time.Local)

  //loop forever
  for {
    timeNow := time.Now()
    if(toTime(timeNow)).After(toTime(startTime)) && toTime(timeNow).Before(toTime(stopTime)) {
      //if(isTimeBefore(timeNow, startTime) && isTimeBefore(stopTime, timeNow)) {
      if lightsOn == false {
        currentColour = startLights(&theBlinkt, currentColour)
        time.Sleep(30 * time.Second)
      }
      lightsOn = true
      currentColour = updateLights(&theBlinkt, currentColour)
    } else {
      if lightsOn == true {
        currentColour = stopLights(&theBlinkt, currentColour)
      }
      lightsOn = false
    }

    //wait before checking for updated colour value
    time.Sleep(10 * time.Minute)
  }
}

/**
  Strips the year month and day off a time
 */
func toTime(t time.Time) time.Time {
  return time.Date(0, 0, 0, t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
}

func GetCheerlightColours() string {
  var netClient = &http.Client {
    Timeout: time.Second * 30,
  }
  resp, getErr := netClient.Get(cheerLightURL)

  if getErr != nil {
    log.Println(getErr)
  }

  if resp.Body != nil {
    defer resp.Body.Close()
  }

  body, readErr := ioutil.ReadAll(resp.Body)

  if readErr != nil {
    log.Println(getErr)
  }

  result := cheerlight{}

  parseErr := json.Unmarshal(body, &result)
  if parseErr != nil {
    log.Println("Can't parse response")
    return "#0000"
  }
  return result.Colour
}

/*
 Set all pixels to the colour c1
 */
func SetAll(theBlinkt *blinkt.Blinkt, c colorful.Color) {
  var r, g, b = c.RGB255()
  log.Printf("Setting pixels to %v, %v, %v", r, g, b)
  theBlinkt.SetAll(int(r), int(g), int(b))
  theBlinkt.Show()
}

/*
  Blend from colour c1 to colour c2
 */
func Blend(theBlink *blinkt.Blinkt, fromColour colorful.Color, toColour colorful.Color) colorful.Color {
  log.Printf("Blending from %s to %s", fromColour.Hex(), toColour.Hex())
  steps := 25
  for i := 0 ; i < steps; i++ {
    opColour := fromColour.BlendHsv(toColour, float64(i)/float64(steps - 1))
    SetAll(theBlink, opColour)
    time.Sleep(250 * time.Millisecond)
  }
  return toColour
}

