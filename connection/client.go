// Author Arsham Shirvani <arshamshirvani@gmail.com>
// license that can be found in the LICENSE file.

package connection

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/arsham/ezconn/lib"
)

// Conf can be passed to NewClient for customizing the Client. (See helpers.go)
// You can set the following settings:
// proxy
// user-agent
// timeout
type Conf func(client *Client)

// Client uses `userAgent` and `proxy` on all nequests throughout its life.
// If you need to change them, simply create a new Client.
// Please note that your can use this Client concurrently as long as the writer is cincurrent safe.
// But updating it's configuration is not concurrent safe.
type Client struct {
	http.Client
	userAgent string
	proxy     string
	writer    io.WriteCloser
	logger    logrus.FieldLogger
}

// NewClient will use `config` to setup the Client.
// See helpers.go to find out what you can use.
// If no configuration has passes, it behaves like a normal client.
func NewClient(logger logrus.FieldLogger, config ...Conf) (httpClient *Client, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			var ok bool
			httpClient = nil
			if err, ok = rec.(error); !ok {
				err = fmt.Errorf("%v", rec)
			}
		}
	}()

	httpClient = &Client{
		writer: os.Stdout,
		logger: logger,
	}

	for _, conf := range config {
		httpClient.UpdateConf(conf)
	}

	if httpClient.proxy != "" {
		if err = httpClient.setProxy(); err != nil {
			return nil, fmt.Errorf("connecting to the proxy: %s", err)
		}
	}

	return
}

// UpdateConf lets you update the configuration at runtime
// It panics if the conf variable is nil
func (c *Client) UpdateConf(conf Conf) {
	conf(c)
}

// Get fetches the `url`. It uses a proxy and/or useragent if defined.
func (c *Client) Get(url string) (written int64, err error) {
	var (
		resp *http.Response
		req  *http.Request
	)
	defer func() {
		if rec := recover(); rec != nil {
			var ok bool
			if err, ok = rec.(error); !ok {
				err = fmt.Errorf("fetching url: %v", rec)
			}
		}
	}()
	c.logger = c.logger.WithField("url", url)
	defer lib.TimeIt(c.logger)()

	if req, err = c.NewRequest("GET", url, nil); err != nil {
		return 0, fmt.Errorf("Client req: %s", err)
	}

	if resp, err = c.Do(req); err != nil {
		return 0, fmt.Errorf("Client do: %s", err)
	}
	defer resp.Body.Close()

	return io.Copy(c.writer, resp.Body)
}

// NewRequest creates the request object and sets the user-agent if is set
func (c *Client) NewRequest(method, url string, body io.Reader) (req *http.Request, err error) {

	if req, err = http.NewRequest(method, url, body); err != nil {
		return nil, fmt.Errorf("Client NewRequst: %s", err)
	}

	c.updateHeaders(req)
	return
}

// Close closes the file
func (c *Client) Close() error {
	return c.writer.Close()
}
