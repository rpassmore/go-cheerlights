package main

import (
  "fmt"
  "time"
  "github.com/lucasb-eyer/go-colorful"
  "github.com/ikester/blinkt"
)

const (
  NUM_PIXELS = 8
)


var  spacing = 360 / 16
var  hue = 0


func main() {
  fmt.Println("Hello, playground")
    
  blinkt := blinkt.NewBlinkt()
  blinkt.ShowAnimOnStart = true
  blinkt.CaptureExit = true
  blinkt.ShowAnimOnExit = true
  blinkt.ClearOnExit = true

  blinkt.Setup()

  //loop forever
  for {
    hue = int(time.Now().Unix() * 100) % 360
    for i := 0; i < NUM_PIXELS; i++ {
      var offset = i * spacing
      var h = ((hue + offset) % 360) / 360.0
      var c = colorful.Hsv(float64(h), 1.0, 1.0)

      blinkt.SetPixel(i, int(c.R), int(c.G), int(c.B))
      blinkt.Show()

      fmt.Printf("Seting pixel %d to %v, %v, %v", i, c.R, c.G, c.B)
      time.Sleep(100 * time.Millisecond)
    }
  }
}

//func check(err error) {
//    if err != nil {
//        fmt.Println(err)
//        os.Exit(1)
//    }
//}
