package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/nachorpaez/osquery-extensions/tables/chrome_extensions_dns"
	"github.com/nachorpaez/osquery-extensions/tables/chrome_preferences"
	osquery "github.com/osquery/osquery-go"
	"github.com/osquery/osquery-go/plugin/table"
)

var packageVersion string // This variable get's set by the build process

func main() {
	var (
		flSocketPath = flag.String("socket", "", "")
		flTimeout    = flag.Int("timeout", 0, "")
		_            = flag.Int("interval", 0, "")
		_            = flag.Bool("verbose", false, "")
	)
	flag.Parse()
	defer glog.Flush()

	// allow for osqueryd to create the socket path otherwise it will error
	time.Sleep(3 * time.Second)

	server, err := osquery.NewExtensionManagerServer(
		"osquery_extension",
		*flSocketPath,
		osquery.ServerTimeout(time.Duration(*flTimeout)*time.Second),
		osquery.ExtensionVersion(packageVersion),
	)
	if err != nil {
		log.Fatalf("Error creating extension: %s\n", err)
	}

	// Create and register a new table plugin with the server.
	// Adding a new table? Add it to the list and the loop below will handle
	// the registration for you.
	plugins := []osquery.OsqueryPlugin{
		table.NewPlugin("chrome_preferences", chrome_preferences.GoogleChromePreferencesColumns(), chrome_preferences.GoogleChromePreferencesGenerate),
		table.NewPlugin("chrome_extensions_dns", chrome_extensions_dns.ChromeExtensionsDNSColumns(), chrome_extensions_dns.ChromeExtensionsDNSGenerate),
	}

	// Platform specific tables
	// if runtime.GOOS == "windows" {
	// If there were windows only tables, they would go here
	// }

	// if runtime.GOOS == "darwin" {
	// 	darwinPlugins := []osquery.OsqueryPlugin{
	// 	}
	// 	plugins = append(plugins, darwinPlugins...)
	// }

	for _, p := range plugins {
		server.RegisterPlugin(p)
	}

	// Start the server. It will run forever unless an error bubbles up.
	if err := server.Run(); err != nil {
		glog.Errorln(err)
		os.Exit(1)
	}
}
