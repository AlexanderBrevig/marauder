package marauder

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	UserName string `yaml:"userName", envconfig:"MARAUDER_USER_NAME"`
	HostName string `yaml:"hostName", envconfig:"MARAUDER_HOST_NAME"`
	Dir      string
	OutDir   string `yaml:"outDir", envconfig:"MARAUDER_OUT_DIR"`
	Colors   struct {
		Button1    string `yaml:"button1"`
		Button2    string `yaml:"button2"`
		Button3    string `yaml:"button3"`
		Background string `yaml:"background"`
		Title      string `yaml:"title"`
		UserName   string `yaml:"userName"`
		At         string `yaml:"at"`
		HostName   string `yaml:"hostName"`
		Dollar     string `yaml:"dollar"`
		Command    string `yaml:"command"`
		Terminal   string `yaml:"terminal"`
	} `yaml:"colors"`
	DatePrefix bool   `yaml:"datePrefix", envconfig:"MARAUDER_DATE_PREFIX"`
	LineLimit  uint16 `yaml:"lineLimit", envconfig:"MARAUDER_LINE_LIMIT"`
}

func (c *Config) Load() {
	c.DatePrefix = true
	c.LineLimit = 80
	c.Colors.Button1 = "ff0000"
	c.Colors.Button2 = "ffff00"
	c.Colors.Button3 = "00ff00"
	c.Colors.Background = "060606"
	c.Colors.Title = "666666"
	c.Colors.UserName = "89b482"
	c.Colors.At = "ffffff"
	c.Colors.HostName = "ea6962"
	c.Colors.Dollar = "ffffff"
	c.Colors.Command = "ff0000"
	c.Colors.Terminal = "ffffff"

	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	c.UserName = user.Username
	c.HostName = hostName
	c.Dir = dir

	readFile(c)
	readEnv(c)
}

func fakeIfEnv(real string, fakeEnv string) string {
	if fakeName := os.Getenv(fakeEnv); fakeName != "" {
		return fakeName
	}
	return real
}
func readFile(cfg *Config) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	confPath := filepath.Join(user.HomeDir, ".marauder.yml")
	f, err := os.Open(confPath)
	if err != nil {
		log.Printf("No config found at %s", confPath)
		return
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		panic(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err)
	}
}
