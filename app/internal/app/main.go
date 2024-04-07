package app

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"

	"github.com/ilya-mezentsev/versiond/app/internal/services"
)

const pidFilename = "pid"

func Main() {
	cfg := mustParseConfig(os.Getenv("CONFIG_PATH"))
	ss := services.New(cfg)

	mustSavePID(cfg.Cache.Dir)

	done := make(chan struct{})
	errs := make(chan error)

	go ss.Run(done, errs)
	go waitForSignal(done)

	for err := range errs {
		fmt.Println("got error:", err)
	}

	fmt.Println("exiting")
}

func waitForSignal(demonStop chan<- struct{}) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{}, 1)

	go func() {
		sig := <-sigs

		fmt.Printf("Got signal: %v. Starting graceful shutdown\n", sig)

		demonStop <- struct{}{}
		done <- struct{}{}
	}()

	<-done
}

func mustSavePID(versiondDir string) {
	err := os.WriteFile(
		path.Join(versiondDir, pidFilename),
		[]byte(strconv.FormatInt(int64(os.Getpid()), 10)),
		0755,
	)
	if err != nil {
		panic(fmt.Errorf("unable to save demon pid to dir: %s, %v", versiondDir, err))
	}
}
