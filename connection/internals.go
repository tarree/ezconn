// Author Arsham Shirvani <arshamshirvani@gmail.com>
// license that can be found in the LICENSE file.

package connection

import (
	"net/http"

	"golang.org/x/net/proxy"
)

func (c *Client) setProxy() error {
	dialer, err := proxy.SOCKS5("tcp", c.proxy, nil, proxy.Direct)
	if err != nil {
		return err
	}
	httpTransport := &http.Transport{}
	httpTransport.Dial = dialer.Dial
	c.Transport = httpTransport
	return nil
}

func (c *Client) updateHeaders(req *http.Request) {
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
}
