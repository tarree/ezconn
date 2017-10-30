// Author Arsham Shirvani <arshamshirvani@gmail.com>
// license that can be found in the LICENSE file.

package connection

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
)

func Test_UserAgent(t *testing.T) {
	var expected string
	handler := func(w http.ResponseWriter, r *http.Request) {
		expected = r.UserAgent()
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ua := "some useragent"
	log := logrus.New()
	log.Out = ioutil.Discard
	c, err := NewClient(log, SetUserAgent(ua))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if _, err = c.Get(server.URL); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if expected != ua {
		t.Errorf("Expected %s, received %s", expected, ua)
	}
}

func Test_Proxy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	p := "bad proxy"
	log := logrus.New()
	log.Out = ioutil.Discard
	c, err := NewClient(log, SetProxy(p))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if _, err = c.Get(server.URL); err == nil {
		t.Error("Expected error, returned nil")
	}
	// TODO: mock a socks5 proxy and write more tests
}

func Test_Timeout(t *testing.T) {
	var duration time.Duration
	log := logrus.New()
	log.Out = ioutil.Discard
	sleepT := time.Second
	started := time.Now()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(sleepT)
	}))
	defer server.Close()

	c, err := NewClient(log, SetTimeout(sleepT))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	select {
	case <-time.After(time.Duration(sleepT * 2)):
		t.Error("Didn't finish in time.")
	default:
		c.Get(server.URL)
		duration = time.Since(started)
	}

	upperbound := sleepT + (10 * time.Millisecond) //nanoseconds
	lowerbound := sleepT - (10 * time.Millisecond)
	if duration > upperbound {
		t.Errorf("It took more than expected: %d (want %d)", duration, upperbound)
	} else if duration < lowerbound {
		t.Errorf("It didn't wait for expected time: %d (want %d)", duration, upperbound)
	}
}
