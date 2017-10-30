// Author Arsham Shirvani <arshamshirvani@gmail.com>
// license that can be found in the LICENSE file.

package lib

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

type timer struct {
	started time.Time
	logger  logrus.FieldLogger
}

// TimeIt is used for timing a function call.
// Usage: defer TimeIt(time.Now(), "function name or label")
func TimeIt(logger logrus.FieldLogger) func() {
	t := &timer{time.Now(), logger}
	return t.done
}

func (t *timer) done() {
	elapsed := time.Since(t.started)
	dur := fmt.Sprintf("%0.3fs", elapsed.Seconds())
	t.logger.WithField("duration", dur).Info("finished")
}
