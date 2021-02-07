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
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
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
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	dateStr := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s-%s", os.Args[1], dateStr)
	f, err := os.Create(filename + ".txt")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	shline := user.Username + "@" + dir + "$ " + strings.Join(os.Args[1:], " ") + "\n"
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
	cmd.Wait()
	time.Sleep(1 * time.Second)
	//TODO: ability to use other tool than scro
	scrot := exec.Command("scrot", "-u", filename+".png")
	err = scrot.Run()
	if err != nil {
		log.Fatalf("scrot failed with %s\n", err)
	}
	xclip := exec.Command("xclip", "-selection", "clipboard", "-t", "image/png", filename+".png")
	err = xclip.Run()
	if err != nil {
		log.Fatalf("xclip failed with %s\n", err)
	}

	dc := gg.NewContext(1200, 628)
	fontPath := filepath.Join("assets", "FiraCode-Regular.ttf")
	if err := dc.LoadFontFace(fontPath, 80); err != nil {
		log.Fatal(err)
	}
	dc.SetColor(color.White)
	s := "TEST text"
	marginX := 50.0
	marginY := -10.0
	textWidth, textHeight := dc.MeasureString(s)
	x := float64(dc.Width()) - textWidth - marginX
	y := float64(dc.Height()) - textHeight - marginY
	dc.DrawString(s, x, y)
	if err := dc.SavePNG("test.png"); err != nil {
		log.Fatal(err)
	}
}
