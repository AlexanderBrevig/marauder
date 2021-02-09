package marauder

import (
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"github.com/gobuffalo/packr"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func ShellLine(userName string, hostName string, dir string) string {
	shline := userName + "@" + hostName + ":" + dir + "$ " + strings.Join(os.Args[1:], " ")
	return shline
}

func DrawConsole(config Config, fileName string, userName string, hostName string, dir string, outStr string) {
	bgMargin := math.Max(0, config.TerminalMargin)
	textMargin := math.Max(0, config.TextMargin)
	buttonRadius := 6.0
	titleFontSize := 12.0
	fontSize := config.FontSize
	// TODO: But why?
	lineLimit := int(math.Min(80.0, float64(config.LineLimit)))
	// TODO: Serioiusly.. wtf is up with 0.58?????333
	lineWidth := fontSize * float64(lineLimit) * 0.58333
	lineSpacing := 1.0

	lines := 0
	for _, line := range strings.Split(outStr, "\n") {
		if len(line) != 0 {
			lines += (len(line) / lineLimit) + 1
		}
	}

	toolbarHeight := bgMargin + textMargin + buttonRadius
	contextHeight := lines*int(fontSize*lineSpacing) + int(bgMargin*2) + int(textMargin*3) + int(toolbarHeight)
	contextWidth := int(lineWidth) + int(bgMargin*2) + int(textMargin*2)

	dc := gg.NewContext(contextWidth, contextHeight)

	dc.SetHexColor(config.Colors.Background)
	dc.DrawRoundedRectangle(bgMargin, bgMargin, float64(contextWidth)-bgMargin*2, float64(contextHeight)-bgMargin*2, bgMargin)
	dc.Fill()

	box := packr.NewBox("../fonts")
	fontBytes, err := box.Find("FiraCode-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}
	font, err := LoadFontFace(fontBytes, titleFontSize)
	if err != nil {
		log.Fatal(err)
	}
	dc.SetFontFace(font)

	dc.SetHexColor(config.Colors.Title)
	title := os.Args[1] + " " + dir
	titleWidth, _ := dc.MeasureString(title)
	dc.DrawString(title, (float64(dc.Width())/2 - (titleWidth / 2)), bgMargin+textMargin+(buttonRadius*1.75))

	font, err = LoadFontFace(fontBytes, fontSize)
	if err != nil {
		log.Fatal(err)
	}
	dc.SetFontFace(font)

	dc.SetHexColor(config.Colors.Button1)
	dc.DrawCircle(bgMargin+textMargin+buttonRadius, bgMargin+textMargin+buttonRadius, buttonRadius)
	dc.Fill()

	dc.SetHexColor(config.Colors.Button2)
	dc.DrawCircle(bgMargin*2+textMargin+buttonRadius*2, bgMargin+textMargin+buttonRadius, buttonRadius)
	dc.Fill()

	dc.SetHexColor(config.Colors.Button3)
	dc.DrawCircle(bgMargin*3+textMargin+buttonRadius*3, bgMargin+textMargin+buttonRadius, buttonRadius)
	dc.Fill()

	renderOffsetX := bgMargin + textMargin
	renderOffsetY := (bgMargin * 2) + (buttonRadius * 2) + (textMargin * 2)

	dc.SetHexColor(config.Colors.UserName)
	dc.DrawString(userName, renderOffsetX, renderOffsetY)

	partialOffset, _ := dc.MeasureString(userName)
	renderOffsetX += partialOffset

	dc.SetHexColor(config.Colors.At)
	dc.DrawString(" @", renderOffsetX, renderOffsetY)

	partialOffset, _ = dc.MeasureString(" @")
	renderOffsetX += partialOffset

	dc.SetHexColor(config.Colors.HostName)
	dc.DrawString(" "+hostName, renderOffsetX, renderOffsetY)

	partialOffset, _ = dc.MeasureString(" " + hostName)
	renderOffsetX += partialOffset

	dc.SetHexColor(config.Colors.Dollar)
	dc.DrawString(" $", renderOffsetX, renderOffsetY)

	partialOffset, _ = dc.MeasureString(" $")
	renderOffsetX += partialOffset

	dc.SetHexColor(config.Colors.Command)
	dc.DrawString(" "+strings.Join(os.Args[1:], " "), renderOffsetX, renderOffsetY)

	dc.SetHexColor(config.Colors.Terminal)
	dc.DrawStringWrapped(outStr, bgMargin+textMargin, bgMargin+textMargin+toolbarHeight+(fontSize*0.58333*lineSpacing), 0.0, 0.0, lineWidth, lineSpacing, gg.AlignLeft)
	path := filepath.Join(config.OutDir, fileName+".png")
	if err := dc.SavePNG(path); err != nil {
		log.Fatal(err)
	}
}

func LoadFontFace(fontBytes []byte, points float64) (font.Face, error) {
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return face, nil
}
