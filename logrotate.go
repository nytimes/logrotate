package logrotate

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type File struct {
	*os.File
	me     sync.Mutex
	path   string
	sighup chan os.Signal
}

func NewFile(path string) (*File, error) {

	lr := &File{
		me:     sync.Mutex{},
		path:   path,
		sighup: make(chan os.Signal, 1),
	}

	if err := lr.reopen(); err != nil {
		return nil, err
	}

	go func() {
		signal.Notify(lr.sighup, syscall.SIGHUP)

		for _ = range lr.sighup {
			fmt.Fprintf(os.Stderr, "%s: Reopening %q\n", time.Now(), lr.path)
			if err := lr.reopen(); err != nil {
				fmt.Fprintf(os.Stderr, "%s: Error reopening: %s\n", time.Now(), err)
			}
		}
	}()

	return lr, nil

}

func (lr *File) reopen() (err error) {
	lr.me.Lock()
	defer lr.me.Unlock()
	lr.File.Close()
	lr.File, err = os.OpenFile(lr.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	return
}

func (lr *File) Write(b []byte) (int, error) {
	lr.me.Lock()
	defer lr.me.Unlock()
	return lr.File.Write(b)
}

func (lr *File) Close() error {
	lr.me.Lock()
	defer lr.me.Unlock()
	signal.Stop(lr.sighup)
	close(lr.sighup)
	return lr.File.Close()
}
