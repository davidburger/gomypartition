package Action

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ActionWaitMode struct {

}

func (a *ActionWaitMode) Process() error {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done
	fmt.Println("SIGTERM or SIGKILL received, exiting ...")
	time.Sleep(500*time.Millisecond)
	return nil
}
