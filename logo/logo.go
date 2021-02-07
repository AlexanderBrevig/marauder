package main

import (
	"github.com/fogleman/gg"
	"log"
	"path/filepath"
)

func main() {
	dc := gg.NewContext(1200, 628)
	fontPath := filepath.Join("assets", "FiraCode-Regular.ttf")
	if err := dc.LoadFontFace(fontPath, 80); err != nil {
		log.Fatal(err)
	}
	s := "marauder"
	l := ">"
	textWidth, textHeight := dc.MeasureString(s)
	logoWidth, logoHeight := dc.MeasureString(l)
	dc.SetHexColor("89b482")
	dc.DrawString(l, 0, logoHeight)
	dc.SetHexColor("7daea3")
	dc.DrawString(s, logoWidth, textHeight)
	if err := dc.SavePNG("test.png"); err != nil {
		log.Fatal(err, textWidth, textHeight, logoWidth, logoHeight) //use vars
	}
}
