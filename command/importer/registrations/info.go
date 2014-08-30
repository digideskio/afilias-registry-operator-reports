package registrations

import (
	"code.google.com/p/gcfg"
	"fmt"
	"github.com/dothiv/afilias-registry-operator-reports/cli"
	"github.com/dothiv/afilias-registry-operator-reports/config"
	"github.com/wsxiaoys/terminal/color"
	"os"
)

const NAME = "import-registrations"

func Help() {
	c := newDefaultConfig()
	cli.HelpBanner(NAME + " @{g}<reportsdir>@{|} @{c}[flags]@{|}")
	os.Stdout.WriteString("Import new registrations into the local database.\n")
	os.Stdout.WriteString("\n")
	color.Fprintln(os.Stdout, "  @{g}reportsdir@{|}           is the directory containing the reports")
	color.Fprintln(os.Stdout, "  @{c}-q@{|}                   be quiet")
	color.Fprintln(os.Stdout, "  @{c}-c=<config-file>@{|}     config file location")
	color.Fprintln(os.Stdout, fmt.Sprintf("                         default: @{c}%s@{|}", c.ConfigFile))
}

type Config struct {
	config.ConfigFromFile
	Quiet      bool
	ReportsDir string
	Database   struct {
		config.Database
	}
}

func newDefaultConfig() (c *Config) {
	c = new(Config)
	c.Quiet = false
	c.ConfigFile = "./importer.ini"
	c.Database.Defaults()
	return
}

func NewConfig() (c *Config, err error) {
	c = newDefaultConfig()

	// Parse args
	if len(os.Args) < 3 {
		err = fmt.Errorf("Missing reportsdir argument")
		return
	}
	c.ReportsDir = os.Args[2]

	// Parse flags
	for i := range os.Args {
		if os.Args[i] == "-q" {
			c.Quiet = true
			continue
		}
		if len(os.Args[i]) > 3 && os.Args[i][0:3] == "-c=" {
			c.ConfigFile = os.Args[i][3:]
		}
	}

	err = gcfg.ReadFileInto(c, c.ConfigFile)
	if err != nil {
		return
	}

	// Check reports dir
	_, err = os.Stat(c.ReportsDir)
	if err != nil {
		return
	}

	return
}
