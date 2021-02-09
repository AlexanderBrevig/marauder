package main

import (
	"fmt"
	"time"

	marauder "github.com/AlexanderBrevig/marauder/pkg"
)

func main() {
	var config marauder.Config
	config.Load()

	shline := marauder.ShellLine(config.UserName, config.HostName, config.Dir)
	outStr, errStr := marauder.Exec()

	dateStr := time.Now().Format("2006-01-02_15:04:05")
	filename := marauder.FileName(outStr + errStr)
	if config.DatePrefix {
		filename = fmt.Sprintf("%s %s", dateStr, filename)
	}

	marauder.WriteFile(config, filename+".txt", shline+"\n"+outStr)
	if len(errStr) > 0 {
		marauder.WriteFile(config, "error "+filename+".txt", shline+"\n"+errStr)
	}
	marauder.DrawConsole(config, filename, config.UserName, config.HostName, config.Dir, outStr)
}
