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

// New creates a new app
func New(c *Config) *App {
	f, err := NewRaft(c)
	if err != nil {
		panic(err)
	}
	return &App{
		f: f,
	}
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
