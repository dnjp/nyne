/*
nstart is used for launching acme along with all of its
dependencies and helpers. An example of its usage can
be found in the (acme start script) ../../mac/Acme.app/Contents/MacOS/acme.

The programs that are launched can be configured in
(config.go) ../../config.go. Stderr and Stdout are
grouped together and written to the default log
file which can be found at $HOME/.config/acme/acme.log.
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dnjp/nyne"
)

var (
	timefmt = time.RFC3339
	deps    = nyne.AcmeDeps
	procs   = nyne.AcmeHelpers
	logloc  = homeloc
)

func homeloc() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(home, "/") {
		home += "/"
	}
	return home + ".local/var/log/acme/acme.log", nil
}

type procErr struct {
	err error
	cmd string
}

func (e procErr) Error() string {
	return fmt.Sprintf("%s failed: %+v", e.cmd, e.err)
}

type fanout struct {
	c        chan struct{}
	stop     chan struct{}
	children []chan struct{}
	mux      sync.RWMutex
}

func newfanout() *fanout {
	return &fanout{
		c:        make(chan struct{}),
		stop:     make(chan struct{}),
		children: make([]chan struct{}, 0),
	}
}

func (f *fanout) newchan() chan struct{} {
	c := make(chan struct{})
	f.mux.Lock()
	f.children = append(f.children, c)
	f.mux.Unlock()
	return c
}

func (f *fanout) send() {
	f.c <- struct{}{}
}

func (f *fanout) kill() {
	f.stop <- struct{}{}
}

func (f *fanout) listen() chan struct{} {
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-f.stop:
				done <- struct{}{}
				return
			case msg := <-f.c:
				f.mux.RLock()
				for _, child := range f.children {
					go func(c chan struct{}, msg struct{}) {
						c <- msg
					}(child, msg)
				}
				f.mux.RUnlock()
			}
		}
	}()
	return done
}

func sendout(cmd *exec.Cmd, pipe io.ReadCloser, errs chan error, out chan string) {
	defer pipe.Close()
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		out <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		errs <- procErr{
			err: err,
			cmd: cmd.Path,
		}
	}
}

func startproc(errs chan error, done chan struct{}, wg *sync.WaitGroup, stdout, stderr chan string, x string, args ...string) {
	sendErr := func(err error) {
		errs <- procErr{
			err: err,
			cmd: x,
		}
	}

	cmd := exec.Command(x, args...)
	cmd.Env = os.Environ()

	// stderr
	sep, err := cmd.StderrPipe()
	if err != nil {
		sendErr(err)
		return
	}
	go sendout(cmd, sep, errs, stderr)

	// stdout
	sop, err := cmd.StdoutPipe()
	if err != nil {
		sendErr(err)
		return
	}
	go sendout(cmd, sop, errs, stdout)

	err = cmd.Start()
	if err != nil {
		sendErr(err)
		return
	}

	go func() {
		err = cmd.Wait()
		if err != nil {
			sendErr(err)
		}
	}()

	stdout <- fmt.Sprintf("[%s] running...\n", x)
	<-done
	cmd.Process.Kill()
	stdout <- fmt.Sprintf("[%s] killed\n", x)
	if wg != nil {
		wg.Done()
	}
}

func newlog() (*os.File, error) {
	loc, err := logloc()
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(filepath.Base(loc), 0755)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(loc, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func usage(name string) {
	fmt.Fprintf(os.Stderr, "%s [acme] [acme args]\n", name)
}

func main() {
	var err error
	args := os.Args
	if len(args) < 2 {
		usage(args[0])
		os.Exit(1)
	}

	// communication
	killed := make(chan struct{}, 1)
	done := newfanout()
	errs := make(chan error)
	deperrs := make(chan error)
	stderr := make(chan string)
	stdout := make(chan string)

	// log stdout/stderr to file
	logfile, err := newlog()
	if err != nil {
		panic(err)
	}
	logwrite := func(fmtstr string, args ...interface{}) {
		s := fmt.Sprintf(fmtstr, args...)
		if s[len(s)-1] != '\n' {
			s += "\n"
		}
		ts := time.Now().Format(timefmt)
		entry := ts + " " + s
		print(entry)
		_, err := logfile.WriteString(entry)
		if err != nil {
			errs <- err
		}
	}

	// set namespace
	os.Setenv("NAMESPACE", nyne.Namespace())
	fmt.Println(os.Getenv("NAMESPACE"))

	// start all deps
	for _, dep := range deps {
		parts := strings.Split(dep, " ")
		if len(dep) >= 2 {
			go startproc(
				deperrs,
				done.newchan(),
				nil,
				stdout,
				stderr,
				parts[0],
				parts[1:]...)
		} else {
			go startproc(
				deperrs,
				done.newchan(),
				nil,
				stdout,
				stderr,
				parts[0])
		}
	}

	// start acme
	a := exec.Command(args[1], args[2:]...)
	a.Env = os.Environ()
	a.Stdin = os.Stdin
	a.Stdout = os.Stdout
	a.Stderr = os.Stderr
	err = a.Start()
	if err != nil {
		panic(err)
	}
	go func() {
		err := a.Wait()
		if err != nil {
			errs <- err
		}
		killed <- struct{}{}
	}()
	time.Sleep(500 * time.Millisecond) // wait for acme

	// start all sub-processes
	var wg sync.WaitGroup
	wg.Add(len(procs))
	for _, proc := range procs {
		parts := strings.Split(proc, " ")
		if len(parts) >= 2 {
			go startproc(
				errs,
				done.newchan(),
				&wg,
				stdout,
				stderr,
				parts[0],
				parts[1:]...)
		} else {
			go startproc(
				errs,
				done.newchan(),
				&wg,
				stdout,
				stderr,
				parts[0])
		}
	}

	// listen for events
	donestopped := done.listen()
	for {
		select {
		case err := <-deperrs:
			// ignore errors when starting dependencies
			// because they likely have already been
			// started before
			logwrite("error: %+v\n", err)
		case err := <-errs:
			if err != nil {
				go func(err error) {
					logwrite("error: %+v\n", err)
					a.Process.Kill()
				}(err)
			}
		case <-killed:
			go func() {
				logwrite("stopping...\n")
				done.send()
				wg.Wait()
				done.kill()
				<-donestopped
				logwrite("done\n")
				err := logfile.Close()
				if err != nil {
					log.Fatal(err)
				}
				os.Exit(0)
			}()
		case so := <-stdout:
			logwrite(so)
		case se := <-stderr:
			logwrite(se)
		}
	}
}
