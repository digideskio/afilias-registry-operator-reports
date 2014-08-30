package server

import (
	"code.google.com/p/gcfg"
	"fmt"
	"github.com/dothiv/afilias-registry-operator-reports/cli"
	"github.com/dothiv/afilias-registry-operator-reports/config"
	"github.com/wsxiaoys/terminal/color"
	"os"
	"strconv"
)

const NAME = "server"

func Help() {
	c := newDefaultConfig()
	cli.HelpBanner(NAME + " @{c}[flags]@{|}")
	os.Stdout.WriteString("Run the server.\n")
	os.Stdout.WriteString("\n")
	color.Fprintln(os.Stdout, fmt.Sprintf("  @{c}-p=<port>@{|}     default: @{c}%d@{|}", c.Port))
}

type Config struct {
	config.ConfigFromFile
	Port     int
	Database struct {
		config.Database
	}
}

func newDefaultConfig() (c *Config) {
	c = new(Config)
	c.Port = 8666
	c.ConfigFile = "./importer.ini"
	c.Database.Defaults()
	return
}

func NewConfig() (c *Config, err error) {
	c = newDefaultConfig()

	// Parse flags
	for i := range os.Args {
		if len(os.Args[i]) > 3 && os.Args[i][0:3] == "-p=" {
			c.Port, _ = strconv.Atoi(os.Args[i][3:])
		}
	}

	err = gcfg.ReadFileInto(c, c.ConfigFile)
	if err != nil {
		return
	}

	return
}
