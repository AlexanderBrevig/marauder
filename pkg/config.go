package marauder

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	UserName string `yaml:"userName", envconfig:"MARAUDER_USER_NAME"`
	HostName string `yaml:"hostName", envconfig:"MARAUDER_HOST_NAME"`
	OutDir   string `yaml:"outDir", envconfig:"MARAUDER_OUT_DIR"`
	Dir      string
	Colors   struct {
		Button1    string `yaml:"button1" envconfig:"MARAUDER_COLOR_BUTTON1"`
		Button2    string `yaml:"button2" envconfig:"MARAUDER_COLOR_BUTTON2"`
		Button3    string `yaml:"button3" envconfig:"MARAUDER_COLOR_BUTTON3"`
		Background string `yaml:"background" envconfig:"MARAUDER_COLOR_BACKGROUND"`
		Title      string `yaml:"title" envconfig:"MARAUDER_COLOR_TITLE"`
		UserName   string `yaml:"userName" envconfig:"MARAUDER_COLOR_USERNAME"`
		At         string `yaml:"at" envconfig:"MARAUDER_COLOR_AT"`
		HostName   string `yaml:"hostName" envconfig:"MARAUDER_COLOR_HOSTNAME"`
		Dollar     string `yaml:"dollar" envconfig:"MARAUDER_COLOR_DOLLAR"`
		Command    string `yaml:"command" envconfig:"MARAUDER_COLOR_COMMAND"`
		Terminal   string `yaml:"terminal" envconfig:"MARAUDER_COLOR_TERMINAL"`
	} `yaml:"colors"`
	DatePrefix     bool    `yaml:"datePrefix", envconfig:"MARAUDER_DATE_PREFIX"`
	LineLimit      uint16  `yaml:"lineLimit", envconfig:"MARAUDER_LINE_LIMIT"`
	FontSize       float64 `yaml:"fontSize", envconfig:"MARAUDER_FONT_SIZE"`
	TerminalMargin float64 `yaml:"terminalMargin", envconfig:"MARAUDER_MARGIN"`
	TextMargin     float64 `yaml:"textMargin", envconfig:"MARAUDER_MARGIN"`
}

func (c *Config) Load() {
	c.DatePrefix = true
	c.LineLimit = 80
	c.FontSize = 16.0
	c.TerminalMargin = 10.0
	c.TextMargin = 10.0
	c.Colors.Button1 = "#ff0000"
	c.Colors.Button2 = "#ffff00"
	c.Colors.Button3 = "#00ff00"
	c.Colors.Background = "#060606"
	c.Colors.Title = "#666666"
	c.Colors.UserName = "#89b482"
	c.Colors.At = "#ffffff"
	c.Colors.HostName = "#ea6962"
	c.Colors.Dollar = "#ffffff"
	c.Colors.Command = "#ff0000"
	c.Colors.Terminal = "#ffffff"

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

func readFile(cfg *Config) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	confPath := filepath.Join(user.HomeDir, ".marauder.yml")
	pathOverride := os.Getenv("MARAUDER_CONFIG")
	if pathOverride != "" {
		confPath = pathOverride
	}
	f, err := os.Open(confPath)
	if err != nil {
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
