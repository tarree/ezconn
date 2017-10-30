// Author Arsham Shirvani <arshamshirvani@gmail.com>
// license that can be found in the LICENSE file.

package main

import (
	"regexp"
	"time"

	"github.com/namsral/flag"

	"github.com/tarree/ezconn/connection"
	"github.com/tarree/ezconn/lib"
)

func main() {
	//bootstaping the application ☺
	var (
		err     error
		written int64
	)
	url, proxy, timeout, useragent, outfile, logLevel := initFlags()
	log := lib.GetLogger(logLevel)

	client, err := connection.NewClient(
		log,
		connection.SetUserAgent(useragent),
		connection.SetProxy(proxy),
		connection.SetTimeout(timeout),
	)
	if err != nil {
		log.Fatalf("creating new client: %s", err)
	}

	if outfile != "" {
		client.UpdateConf(connection.SetFileOutput(outfile))
	}
	defer client.Close()

	if written, err = client.Get(url); err != nil {
		log.Fatal(err)
	}
	log.Debugf("Received %d bytes", written)
}

func initFlags() (url string, proxy string, timeout time.Duration, userAgent string, outfile string, logLevel string) {
	flag.StringVar(&proxy, "proxy", "", "Sets the proxy. Direct connection if you don't set it.")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "Sets the timeout on request")
	flag.StringVar(&userAgent, "user-agent", "", "Sets the user-agent on client. Default is golang's default user-agent")
	flag.StringVar(&outfile, "o", "", "Saves the output to a file.")
	flag.StringVar(&logLevel, "logLevel", "error", "Log level. Values are: debug, info, warn, error. Default level is error")
	flag.Parse()

	// The last argument without specifying a flag is the url
	url = flag.Arg(0)
	if matched, _ := regexp.Match("https?://", []byte(url)); !matched {
		url = "http://" + url
	}
	return
}
