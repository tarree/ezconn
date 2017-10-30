// Author Arsham Shirvani <arshamshirvani@gmail.com>
// license that can be found in the LICENSE file.

package connection

import (
	"io"
	"os"
	"time"
)

// Proxy tells you what proxy is being used in this session
func (c *Client) Proxy() string {
	return c.proxy
}

// SetProxy can be passed to the NewClient to configure the proxy.
// Usage: NewClient(SetProxy(url))
func SetProxy(proxy string) Conf {
	return func(cc *Client) {
		if proxy != "" {
			cc.proxy = proxy
			cc.logger = cc.logger.WithField("proxy", proxy)
		}
	}
}

// UserAgent tells you what user-agent string is being used in this session
func (c *Client) UserAgent() string {
	return c.userAgent
}

// SetUserAgent can be passed to the NewClient to configure the user-agent string.
// Usage: NewClient(SetUserAgent(userAgent))
func SetUserAgent(userAgent string) Conf {
	return func(cc *Client) {
		if userAgent != "" {
			cc.userAgent = userAgent
			cc.logger = cc.logger.WithField("userAgent", userAgent)
		}
	}
}

// SetTimeout can be passed to the NewClient to configure connection timeout.
// Usage: NewClient(SetTimeout(timeout))
func SetTimeout(timeout time.Duration) Conf {
	return func(cc *Client) {
		cc.Timeout = timeout
		cc.logger = cc.logger.WithField("timeout", timeout)
	}
}

// SetWriter defined the output of the operation. The default is stdout.
// This client does not close your writer and you should do that yourself.
func SetWriter(w io.WriteCloser) Conf {
	return func(cc *Client) {
		cc.writer = w
	}
}

// SetFileOutput writes the output to a file. It exits the program if the file is not writable
// This client does not close your file and you should do that yourself.
func SetFileOutput(filename string) Conf {
	return func(cc *Client) {
		out, err := os.Create(filename)
		if err != nil {
			cc.logger.Warnf("createing file failed, falling back to stdout: %s", err)
			return
		}
		cc.logger = cc.logger.WithField("file", filename)
		cc.writer = out
	}
}
