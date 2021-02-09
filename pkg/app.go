package marauder

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kennygrant/sanitize"
)

func FileName(content string) string {
	filename := strings.Join(os.Args[1:], " ")
	filename = filename + "-" + strconv.Itoa(len(content))
	filename = sanitize.Path(filename)
	filename = strings.ReplaceAll(filename, "/", "[slash]")
	return filename
}

func Exec() (string, string) {
	var cmd *exec.Cmd
	if len(os.Args) > 2 {
		cmd = exec.Command(os.Args[1], os.Args[2:]...)
	} else {
		cmd = exec.Command(os.Args[1])
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	cmd.Wait()
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	return outStr, errStr
}

func WriteFile(config Config, fileName string, content string) {
	path := filepath.Join(config.OutDir, fileName)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	fmt.Fprint(w, content)
	w.Flush()
}
