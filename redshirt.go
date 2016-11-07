// redshirt is a package for handling signals and/or dying gracefully
package redshirt

import (
	"os"
	"os/signal"
	"syscall"
)

const (
	SIGHUP  = syscall.SIGHUP
	SIGINT  = syscall.SIGINT
	SIGTERM = syscall.SIGTERM
	SIGQUIT = syscall.SIGQUIT
)

var receivers map[os.Signal][]Receiver

type Receiver interface {
	Signal(os.Signal) error
}

func Register(r Receiver, signals ...os.Signal) {
	register(r, signals)
}

// Anonymous function support
type ReceiverFunc func(os.Signal) error

// ReceiverFunc.Signal is used to fulfill Receiver requirements
func (r ReceiverFunc) Signal(sig os.Signal) error {
	err := r(sig)
	return err
}

// RegisterFunc allows for a function to be registered instead of a type.
func RegisterFunc(f func(os.Signal) error, signals ...os.Signal) {
	var r ReceiverFunc = f
	register(r, signals)
}

func register(r Receiver, signals []os.Signal) {
	for i := range signals {
		if receivers == nil {
			receivers = make(map[os.Signal][]Receiver)
		}
		// Check if key exists, create if it doesn't.
		if _, ok := receivers[signals[i]]; !ok {
			receivers[signals[i]] = []Receiver{r}
		} else {
			// Add Receiver type to map
			receivers[signals[i]] = append(receivers[signals[i]], r)
		}
	}
}

func Listen() {
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exit_chan := make(chan int)
	go func() {
		for {
			s := <-signal_chan
			switch s {
			// kill -SIGHUP XXXX
			case syscall.SIGHUP:
				notify(s)
			// kill -SIGINT XXXX or Ctrl+c
			case syscall.SIGINT:
				// kill -SIGTERM XXXX
				notify(s)
				exit_chan <- 0
			case syscall.SIGTERM:
				notify(s)
				exit_chan <- 0
			// kill -SIGQUIT XXXX
			case syscall.SIGQUIT:
				notify(s)
				exit_chan <- 0
			default:
				exit_chan <- 1
			}
		}
	}()

	code := <-exit_chan
	os.Exit(code)
}

func notify(sig os.Signal) {
	for i := range receivers[sig] {
		receivers[sig][i].Signal(sig)
	}
}
