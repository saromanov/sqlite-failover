package failover

import (
	"net"
	"sync"
	"time"
)

// App defines main struct for the program
type App struct {
	f    *Failover
	l    net.Listener
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
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		panic(err)
	}
	return &App{
		f:    f,
		c:    c,
		l:    lis,
		quit: make(chan struct{}),
	}
}

// Start provides init of the app
func (a *App) Start() {
	t := time.NewTicker(time.Duration(a.c.Interval) * time.Second)
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

// Close provides closing of the app
func (a *App) Close() {
	select {
	case <-a.quit:
		return
	default:
		break
	}

	close(a.quit)
}

func (a *App) checkCluster() {

}
