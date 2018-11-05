package failover

import (
	"net"
	"sync"
	"time"
)

// App defines main struct for the program
type App struct {
	f    *Failover
	l    *net.Listener
	m    sync.Mutex
	c    *Config
	quit chan struct{}
}

func (a *App) Start() {
	t := time.NewTicker(time.Duration(a.c.Interval) * time.Millisecond)
	defer func() {
		t.Stop()
	}()

	for {
		select {
		case <-t.C:
			a.checkCluster()
		case <-a.quit:
			return
		}
	}
}

func (a *App) checkCluster() {

}
