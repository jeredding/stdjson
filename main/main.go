package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/nkvoll/stdjson"
	"github.com/nkvoll/stdjson/config"
	"github.com/nkvoll/stdjson/rewriter"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

var (
	configFile = flag.String("config", "", "config file to load configuration from")
	debug      = flag.Bool("debug", false, "if enabled, shows debug logs")
	_          = flag.String("-", "", "anything after the first -- is executed as a subprocess")
)

func main() {
	var command []string

	for i, arg := range os.Args {
		if arg == "--" {
			command = os.Args[i+1:]
			os.Args = os.Args[:i]
			break
		}
	}

	flag.Parse()
	if *configFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	if len(command) == 0 {
		log.Errorln("No command to run.")
		os.Exit(2)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(&stdjson.LogFormatter{})

	config, err := config.LoadConfigFromFile(*configFile)
	if err != nil {
		log.Panicln(err)
	}

	log.Debugln("Running subprocess", command)
	p := stdjson.NewProcess(command[0], command[1:]...)

	var (
		stdout rewriter.Rewriter = rewriter.NewNoop(os.Stdout)
		stderr rewriter.Rewriter = rewriter.NewNoop(os.Stderr)
	)

	if config.Stdout != nil {
		log.Debugln("Rewriting stdout")
		o, err := rewriter.RewriterForStreamConfig(config.Stdout, os.Stdout)
		if err != nil {
			log.Panicln(err)
		}
		stdout = o
	}

	if config.Stderr != nil {
		log.Debugln("Rewriting stderr")
		o, err := rewriter.RewriterForStreamConfig(config.Stderr, os.Stderr)
		if err != nil {
			log.Panicln(err)
		}
		stderr = o
	}

	wg := sync.WaitGroup{}
	wg.Add(2) // stdout + stderr

	go func() {
		log.Debugln("Starting stdout rewriter")
		if err := stdout.Run(); err != nil {
			log.Errorln("Error from running out rewriter:", err)
		}
		log.Debugln("Closed stdout rewriter")
		wg.Done()
	}()

	go func() {
		log.Debugln("Starting stderr rewriter")
		if err := stderr.Run(); err != nil {
			log.Errorln("Error from running err rewriter:", err)
		}
		log.Debugln("Closed stderr rewriter")
		wg.Done()
	}()

	p.Stdout = stdout
	p.Stderr = stderr

	runErr := p.Run()

	if err := stdout.Close(); err != nil {
		log.Errorln("Error from closing out rewriter:", err)
	}
	if err := stderr.Close(); err != nil {
		log.Errorln("Error from closing out rewriter:", err)
	}
	wg.Wait()

	if runErr != nil {
		if exiterr, ok := runErr.(*exec.ExitError); ok {
			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			} else {
				os.Exit(-1)
			}
		}
	}

	os.Exit(0)
}
