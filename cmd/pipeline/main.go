package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/TrueBlocks/trueblocks-math/internal/pipeline"
)

func main() {
	configPath := flag.String("config", pipeline.DefaultConfigPath(), "path to config.yaml")
	dryRun := flag.Bool("dry-run", false, "override config to force dry-run mode")
	once := flag.Bool("once", false, "run a single cycle and exit")
	port := flag.Int("port", 0, "override dashboard port")
	flag.Parse()

	cfg, err := pipeline.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if *dryRun {
		cfg.Pipeline.DryRun = true
	}
	if *port > 0 {
		cfg.Dashboard.Port = *port
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	logBuf := pipeline.NewLogBuffer(os.Stdout, 1000)
	runner := pipeline.NewRunner(cfg, cwd)
	runner.ConfigPath = *configPath
	runner.CLIDryRun = *dryRun
	runner.Log = log.New(logBuf, "", 0)

	var dash *pipeline.Dashboard
	runner.Client.OnRetry = func(attempt, maxAttempts int, err error, nextIn time.Duration) {
		if dash == nil {
			return
		}
		if err == nil {
			dash.SetRetryMessage("")
			runner.Log.Println("API connection restored")
			return
		}
		msg := fmt.Sprintf("API retry %d/%d \u2014 next attempt in %ds (%v)", attempt, maxAttempts, int(nextIn.Seconds()), err)
		dash.SetRetryMessage(msg)
	}

	if err := runner.DiscoverProjects(); err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering projects: %v\n", err)
		os.Exit(1)
	}

	if err := runner.LoadState(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading pipeline state: %v\n", err)
		os.Exit(1)
	}

	for _, ps := range runner.Projects {
		summary := ps.Summary()
		total := 0
		for _, v := range summary {
			total += v
		}
		runner.Log.Printf("[%s] %d essays (pending=%d research=%d outline=%d draft=%d factcheck=%d done=%d)",
			ps.Project, total, summary["pending"], summary["research"], summary["outline"],
			summary["draft"], summary["factcheck"], summary["done"])
	}

	if cfg.Pipeline.DryRun {
		runner.Log.Println("Mode: DRY RUN (no API calls)")
	} else {
		runner.Log.Println("Mode: LIVE")
	}

	interval := cfg.Pipeline.CycleInterval
	if interval <= 0 {
		interval = 15
	}

	dash = pipeline.NewDashboard(runner, cfg.Dashboard.Port, interval, *configPath, logBuf)

	go func() {
		if err := dash.Start(); err != nil {
			runner.Log.Printf("Dashboard error: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cycleRunning atomic.Bool

	runCycle := func() {
		if !cycleRunning.CompareAndSwap(false, true) {
			runner.Log.Println("Cycle already running, skipping")
			return
		}
		dash.SetCycleRunning()
		actions, err := runner.RunCycle(ctx)
		if err != nil && ctx.Err() == nil {
			runner.Log.Printf("Cycle error: %v", err)
		}
		dash.SetLastLog(actions)
		for _, a := range actions {
			runner.Log.Println(a)
		}
		dash.SetCycleFinished()
		dash.SetNextCycleAt(time.Now().Add(time.Duration(dash.CycleInterval) * time.Second))
		cycleRunning.Store(false)
	}

	if *once && cfg.Pipeline.Debug == "" {
		runCycle()
		return
	}

	sigCh := make(chan os.Signal, 4)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	var sigCount int32

	go func() {
		for range sigCh {
			n := atomic.AddInt32(&sigCount, 1)
			if n >= 3 {
				fmt.Fprintf(os.Stderr, "\nForce quit.\n")
				os.Exit(1)
			}
			if n == 1 {
				runner.Log.Println("Shutting down (cancel in-flight work)...")
				cancel()
			} else {
				runner.Log.Println("Interrupt again to force quit...")
			}
		}
	}()

	runner.Log.Printf("Dashboard: http://127.0.0.1:%d", cfg.Dashboard.Port)

	if cfg.Pipeline.Debug != "" {
		runCycle()
		runner.Log.Printf("DEBUG complete — dashboard still running at http://127.0.0.1:%d", cfg.Dashboard.Port)
		<-ctx.Done()
		for cycleRunning.Load() {
			time.Sleep(100 * time.Millisecond)
		}
		printFinalSummary(runner)
		return
	}

	runner.Log.Printf("Auto-cycling every %d seconds (step from dashboard to run immediately)", interval)

	dash.SetNextCycleAt(time.Now().Add(time.Duration(interval) * time.Second))
	go runCycle()

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			for cycleRunning.Load() {
				time.Sleep(100 * time.Millisecond)
			}
			printFinalSummary(runner)
			return

		case newInterval := <-dash.IntervalChannel():
			interval = newInterval
			ticker.Reset(time.Duration(interval) * time.Second)
			dash.SetNextCycleAt(time.Now().Add(time.Duration(interval) * time.Second))
			runner.Log.Printf("Cycle interval changed to %d seconds", interval)

		case <-dash.StepChannel():
			ticker.Reset(time.Duration(interval) * time.Second)
			go runCycle()

		case <-ticker.C:
			go runCycle()
		}
	}
}

func printFinalSummary(runner *pipeline.Runner) {
	runner.LoadState()
	for _, ps := range runner.Projects {
		s := ps.Summary()
		runner.Log.Printf("[%s] Final: pending=%d research=%d outline=%d draft=%d factcheck=%d done=%d (cost=$%.2f)",
			ps.Project, s["pending"], s["research"], s["outline"], s["draft"], s["factcheck"], s["done"], ps.TotalCost)
	}
}
