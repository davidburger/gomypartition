package Factory

import (
	"gomypartition/Application/Action"
	"strconv"
	"os"
	"log"
)

func NewDockerTestAction() (a Action.DockerTestAction) {
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))

	if err != nil {
		log.Fatal(err)
	}

	maxRowsPerWorker, err := strconv.Atoi(os.Getenv("MAX_ROWS_PER_WORKER"))

	if err != nil {
		log.Fatal(err)
	}

	var createTable bool
	if os.Getenv("CREATE_TABLE") == "yes" {
		createTable = true
	}

	a.Setup(concurrency, maxRowsPerWorker, createTable)

	return
}
