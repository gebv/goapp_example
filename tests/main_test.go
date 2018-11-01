package tests

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

var (
	Ctx context.Context

	addressF = flag.String("address", "127.0.0.1:3030", "Address listen (host and port).")
)

func runMake(arg string) {
	args := []string{"-C", "..", arg}
	cmd := exec.Command("make", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Print(strings.Join(cmd.Args, " "))
	if err := cmd.Run(); err != nil {
		log.Panic(err)
	}
}

func runGo(bin string, args ...string) {
	cmd := exec.Command(filepath.Join("..", "bin", bin), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), `GORACE="halt_on_error=1"`)
	log.Print(strings.Join(cmd.Args, " "))
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	go func() {
		<-Ctx.Done()
		_ = cmd.Process.Signal(unix.SIGTERM)
	}()

	if err := cmd.Wait(); err != nil {
		log.Print(err)
	}
}

func TestMain(m *testing.M) {

	log.SetPrefix("testmain: ")
	log.SetFlags(0)

	flag.Parse()
	if testing.Short() {
		log.Print("-short flag is passed, skipping integration tests.")
		os.Exit(0)
	}

	var cancel context.CancelFunc
	Ctx, cancel = context.WithCancel(context.Background())

	// handle termination signals: first one cancels context, force exit on the second one
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, unix.SIGTERM, unix.SIGINT)
	go func() {
		s := <-signals
		log.Printf("Got %s, shutting down...", unix.SignalName(s.(unix.Signal)))
		cancel()

		s = <-signals
		log.Panicf("Got %s, exiting!", unix.SignalName(s.(unix.Signal)))
	}()

	var exitCode int
	defer func() {
		if p := recover(); p != nil {
			panic(p)
		}
		os.Exit(exitCode)
	}()

	runMake("build-race")
	go runGo("goapp", "--address=127.0.0.1:3031")

	time.Sleep(time.Second)

	exitCode = m.Run()
	cancel()
}
