package main

import (
	"github.com/fogleman/gg"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	fontPath := filepath.Join("assets", "FiraCode-Regular.ttf")

	size := 256
	if len(os.Args) > 1 {
		i, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Println(err)
			size = 256
		}
		size = i
	}
	S := float64(size)
	log.Println(S)
	dc := gg.NewContext(int(S), int(S))
	dc.SetHexColor("282828")
	dc.DrawRegularPolygon(6, S/2, S/2, S/2, 0)
	dc.Fill()
	s := "MARAUDER"
	l := ">"
	if err := dc.LoadFontFace(fontPath, S/6.4); err != nil {
		log.Fatal(err)
	}
	dc.Push()
	dc.Translate(0, S/15.5)
	textWidth, textHeight := dc.MeasureString(s)
	xOffset := float64(dc.Width())/2 - float64(textWidth)/2
	dc.SetHexColor("7daea3")
	dc.DrawString(s, xOffset, S/2+textHeight)

	if err := dc.LoadFontFace(fontPath, S/1.54); err != nil {
		log.Fatal(err)
	}
	logoWidth, logoHeight := dc.MeasureString(l)
	dc.SetHexColor("89b482")
	dc.DrawString(l, S/2-logoWidth/2, S/2)
	dc.Pop()
	if err := dc.SavePNG("assets/logo.png"); err != nil {
		log.Fatal(err, textWidth, textHeight, logoWidth, logoHeight) //use vars
	}
}
