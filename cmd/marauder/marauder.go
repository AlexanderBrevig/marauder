package main

import (
	marauder "github.com/AlexanderBrevig/marauder/pkg"
)

func main() {
	var config marauder.Config
	config.Load()

	outStr, errStr := marauder.Exec()

	filename := marauder.FileName(config, outStr+errStr)

	marauder.WriteFile(config, filename+".txt", outStr)
	marauder.WriteFile(config, "error "+filename+".txt", errStr)
	marauder.DrawConsole(config, filename, config.UserName, config.HostName, config.Dir, outStr)
}
