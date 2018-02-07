package Action

import (
	"fmt"
	"gomypartition/Application/Service"
	"gomypartition/Application/Service/Factory"
	"log"
	"sync"
)

const testTable = "test_table_partitioned"

type DockerTestAction struct {
	concurrency      int
	maxRowsPerWorker int
	createTable      bool
}

func (a *DockerTestAction) Setup(concurrency int, maxRowsPerWorker int, createTable bool) {
	a.concurrency = concurrency
	a.maxRowsPerWorker = maxRowsPerWorker
	a.createTable = createTable
}

func (a *DockerTestAction) Process() error {
	var err error

	test := Factory.NewTestService()

	fmt.Println("Action partition docker test is running ...");

	defer test.CloseConnection()

	if a.createTable {
		err = test.CreateTable(testTable)

		if err != nil {
			return err
		}
	}

	var wg sync.WaitGroup

	wg.Add(1)

	fmt.Printf("Trying to insert data with %d concurrent clients (each %d rows) ...\n", a.concurrency, a.maxRowsPerWorker)

	for i := 1; i <= a.concurrency; i++ {
		go func (tt *Service.TestService, max int, clientId int) {
			var counter int
			for x:=0; x < max; x++ {

				err = test.InsertRecord(testTable)

				if err != nil {
					log.Panic(err)
				}

				counter++
			}

			wg.Done()
		} (&test, a.maxRowsPerWorker, i)
	}

	wg.Wait() //wait for all runnning goroutines to finish the work

	fmt.Println("done.")

	return nil
}
