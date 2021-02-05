package main

//=============================================================================
// WORK IN PROGRESS
//=============================================================================

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
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
}
