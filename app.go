package main

import (
  MQTT "github.com/eclipse/paho.mqtt.golang"
  "github.com/go-co-op/gocron"
  "github.com/ikester/blinkt"
  "github.com/lucasb-eyer/go-colorful"
  "log"
  "time"
)
const startTime = "19:30:00"
const stopTime = "23:00:00"
const cheerLightURL = "tcp://mqtt.cheerlights.com:1883"
const CheerlightsTopic = "cheerlightsRGB"

var currentColour = colorful.LinearRgb(0,0,0)
var theBlinkt = blinkt.NewBlinkt()

func startLights(theBlinkt *blinkt.Blinkt, mqtt MQTT.Client) {
  log.Println("Lights started")
  theBlinkt.ShowInitialAnim()
  currentColour = Blend(theBlinkt, currentColour, colorful.FastWarmColor())

  //subscribe to the topic and request messages to be delivered
  //at a maximum qos of zero, wait for the receipt to confirm the subscription
  if token := mqtt.Subscribe(CheerlightsTopic, 0, colourChangedListener); token.Wait() && token.Error() != nil {
    log.Println(token.Error())
  }
}

func stopLights(theBlinkt *blinkt.Blinkt, mqtt MQTT.Client) {
  log.Println("Lights stopped")
  mqtt.Unsubscribe(CheerlightsTopic)
  newColour := colorful.LinearRgb(0,0,0)
  currentColour = Blend(theBlinkt, currentColour, newColour)
}

func updateLights(theBlinkt *blinkt.Blinkt, newColourStr string) {
   newColour, getErr := colorful.Hex(newColourStr)
   if getErr != nil {
     log.Println(getErr)
     newColour = colorful.HappyColor()
   }
  currentColour = Blend(theBlinkt, currentColour, newColour)
}



func main() {
  log.Println("Cheerlights started")

  theBlinkt.ShowAnimOnStart = false
  theBlinkt.CaptureExit = true
  theBlinkt.ShowAnimOnExit = true
  theBlinkt.ClearOnExit = true
  theBlinkt.Setup()
  log.Println("The Blinkt setup")

  mqttOpts := MQTT.NewClientOptions()
  mqttOpts.AddBroker(cheerLightURL)
  mqttOpts.SetClientID("go-CheerLights-NHS")
  //mqttOpts.SetConnectionLostHandler(connLost)
  mqttOpts.SetAutoReconnect(true)

  //create and start a client using the above ClientOptions
  mqtt := MQTT.NewClient(mqttOpts)
  if token := mqtt.Connect(); token.Wait() && token.Error() != nil {
    panic(token.Error())
  }
  log.Println("MQTT setup")

  timeNow := toTime(time.Now())
  st, _ := time.Parse("15:04:05", startTime)
  et, _ := time.Parse("15:04:05", stopTime)
  if timeNow.After(st) && timeNow.Before(et) {
    startLights(&theBlinkt, mqtt)
  }

  cronSched := gocron.NewScheduler(time.Local)
  cronSched.Every(1).Day().At(startTime).Do(startLights, &theBlinkt, mqtt)
  cronSched.Every(1).Day().At(stopTime).Do(stopLights, &theBlinkt, mqtt)
  cronSched.StartBlocking()
}

func colourChangedListener(client MQTT.Client, msg MQTT.Message) {
  log.Println("Colour: ", msg.Topic(), msg.MessageID(), msg.Payload())
  updateLights(&theBlinkt, string(msg.Payload()[:]))
}

/**
  Strips the year month and day off a time
 */
func toTime(t time.Time) time.Time {
  res, _ := time.Parse("15:04:05", t.Format("15:04:05"))
  return res
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

