package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	//TODO: load config
	//TODO: if fake shell line
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	userName := user.Username
	if fakeUser := os.Getenv("MARAUDER_FAKE_USER"); fakeUser != "" {
		userName = fakeUser
	}
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	if fakeHostname := os.Getenv("MARAUDER_FAKE_HOSTNAME"); fakeHostname != "" {
		hostname = fakeHostname
	}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if fakeDir := os.Getenv("MARAUDER_FAKE_DIR"); fakeDir != "" {
		dir = fakeDir
	}

	//TODO: verify args len
	//TODO: print usage
	cmd := exec.Command(os.Args[1], os.Args[2:]...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	cmd.Wait()
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	dateStr := time.Now().Format("2006-01-02_15:04:05")
	filename := fmt.Sprintf("%s-%s", dateStr, os.Args[1])
	f, err := os.Create(filename + ".txt")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	shline := userName + "@" + hostname + ":" + dir + "$ " + strings.Join(os.Args[1:], " ") + "\n"
	fmt.Fprint(w, shline)
	fmt.Fprint(w, outStr)
	w.Flush()

	if len(errStr) > 0 {
		errf, err := os.Create(fmt.Sprintf("err-%s", filename+".txt"))
		defer errf.Close()
		if err != nil {
			panic(err)
		}
		errw := bufio.NewWriter(errf)
		fmt.Fprint(errw, shline)
		fmt.Fprint(errw, errStr)
		errw.Flush()
	}

	bgMargin := 10.0
	textMargin := 10.0
	buttonRadius := 6.0
	titleFontSize := 12.0
	fontSize := 16.0
	// TODO: But why?
	adjustedFontsize := fontSize * 0.75
	// TODO: minium 80
	lineLimit := 80
	// TODO: Serioiusly.. wtf is up with 0.58?????333
	lineWidth := fontSize * float64(lineLimit) * 0.58333
	lineSpacing := 1.5
	fontPath := filepath.Join("assets", "FiraCode-Regular.ttf")

	lines := 0
	for _, line := range strings.Split(outStr, "\n") {
		if len(line) != 0 {
			lines += (len(line) / lineLimit) + 1
		}
	}

	toolbarHeight := bgMargin + textMargin + buttonRadius
	contextHeight := lines*int(adjustedFontsize*lineSpacing) + int(bgMargin*2) + int(textMargin*3) + int(toolbarHeight)
	contextWidth := int(lineWidth) + int(bgMargin*2) + int(textMargin*2)

	dc := gg.NewContext(contextWidth, contextHeight)

	dc.SetColor(color.Black)
	dc.DrawRoundedRectangle(bgMargin, bgMargin, float64(contextWidth)-bgMargin*2, float64(contextHeight)-bgMargin*2, bgMargin)
	dc.Fill()

	if err := dc.LoadFontFace(fontPath, titleFontSize); err != nil {
		log.Fatal(err)
	}

	dc.SetHexColor("666666")
	title := os.Args[1] + " " + dir
	titleWidth, _ := dc.MeasureString(title)
	dc.DrawString(title, (float64(dc.Width())/2 - (titleWidth / 2)), bgMargin+textMargin+(buttonRadius*1.75))

	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		log.Fatal(err)
	}

	dc.SetHexColor("ff0000")
	dc.DrawCircle(bgMargin+textMargin+buttonRadius, bgMargin+textMargin+buttonRadius, buttonRadius)
	dc.Fill()

	dc.SetHexColor("ffff00")
	dc.DrawCircle(bgMargin*2+textMargin+buttonRadius*2, bgMargin+textMargin+buttonRadius, buttonRadius)
	dc.Fill()

	dc.SetHexColor("00ff00")
	dc.DrawCircle(bgMargin*3+textMargin+buttonRadius*3, bgMargin+textMargin+buttonRadius, buttonRadius)
	dc.Fill()

	renderOffsetX := bgMargin + textMargin
	renderOffsetY := (bgMargin * 2) + (buttonRadius * 2) + (textMargin * 2)

	dc.SetHexColor("89b482")
	dc.DrawString(userName, renderOffsetX, renderOffsetY)

	partialOffset, _ := dc.MeasureString(userName)
	renderOffsetX += partialOffset

	dc.SetColor(color.White)
	dc.DrawString(" @", renderOffsetX, renderOffsetY)

	partialOffset, _ = dc.MeasureString(" @")
	renderOffsetX += partialOffset

	dc.SetHexColor("ea6962")
	dc.DrawString(" "+hostname, renderOffsetX, renderOffsetY)

	partialOffset, _ = dc.MeasureString(" " + hostname)
	renderOffsetX += partialOffset

	dc.SetColor(color.White)
	dc.DrawString(" $", renderOffsetX, renderOffsetY)

	partialOffset, _ = dc.MeasureString(" $")
	renderOffsetX += partialOffset

	dc.SetHexColor("ff0000")
	dc.DrawString(" "+strings.Join(os.Args[1:], " "), renderOffsetX, renderOffsetY)

	dc.SetColor(color.White)
	dc.DrawStringWrapped(outStr, bgMargin+textMargin, bgMargin+textMargin+toolbarHeight+(fontSize*0.58333*lineSpacing), 0.0, 0.0, lineWidth, lineSpacing, gg.AlignLeft)
	if err := dc.SavePNG(filename + ".png"); err != nil {
		log.Fatal(err)
	}
}
