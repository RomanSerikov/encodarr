package job_health

import (
	"context"
	"time"

	"github.com/BrenekH/encodarr/controller"
)

func NewChecker(ds controller.HealthCheckerDataStorer, ss controller.SettingsStorer) Checker {
	return Checker{
		ds: ds,
		ss: ss,

		lastCheckTime: time.Unix(0, 0),
		nowSincer:     TimeNowSince{},
	}
}

type Checker struct {
	ds controller.HealthCheckerDataStorer
	ss controller.SettingsStorer

	lastCheckTime time.Time
	nowSincer     NowSincer
}

// Run loops through the provided slice of dispatched jobs and checks if any have
// surpassed the allowed time between updates, if the Health Check timing interval has expired.
func (c *Checker) Run() (uuidsToNull []controller.UUID) {
	if c.nowSincer.Since(c.lastCheckTime) >= time.Duration(c.ss.HealthCheckInterval()) {
		c.lastCheckTime = c.nowSincer.Now()

		djs := c.ds.DispatchedJobs()

		for _, v := range djs {
			if c.nowSincer.Since(v.LastUpdated) >= time.Duration(c.ss.HealthCheckTimeout()) {
				// Since DeleteJob may be blocked by an IO error of some sort (SQLiteDB.SetMaxOpenConns should fix this issue but just in case)
				//   attempt to delete the job up to a hundred times.
				jobDeleted := false
				for i := 0; i < 100; i++ {
					if err := c.ds.DeleteJob(v.UUID); err == nil {
						jobDeleted = true
						break
					}
					time.Sleep(time.Microsecond * 2)
				}

				if jobDeleted {
					uuidsToNull = append(uuidsToNull, v.UUID)
				}
			}
		}
	}
	return
}

// Start just satisfies the controller.HealthChecker interface.
// There is no implemented functionality.
func (c *Checker) Start(ctx *context.Context) {}
