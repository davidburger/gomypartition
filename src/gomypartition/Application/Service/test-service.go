package Service

import (
	"fmt"
	"gomypartition/Application/Database"
	"time"
	"log"
	"github.com/satori/go.uuid"
	"math/rand"
)

type TestService struct {
	Database.DbDriver
	partitioner Database.DatePartitioner
	insertTsMin int64
	insertTsMax int64
}

const DOCKER_DB_HOST = "db"

func (t *TestService) getPartitionTool() Database.DatePartitioner {

	if t.partitioner.Initialized == false {
		var partitioner Database.DatePartitioner
		var err error

		currentYear := time.Now().Year();
		partitioner.Column = "event_date"

		partitioner.MinDate, err = time.Parse(time.RFC3339,
			fmt.Sprintf("%d-01-01T00:00:00+00:00", currentYear-1))

		if (err != nil) {
			log.Panic(err)
		}

		partitioner.MaxDate, err = time.Parse(time.RFC3339,
			fmt.Sprintf("%d-12-21T23:59:59+00:00", currentYear))

		if (err != nil) {
			log.Panic(err)
		}

		partitioner.Initialized = true

		t.partitioner = partitioner

		t.insertTsMin = partitioner.MinDate.Unix();
		t.insertTsMax = partitioner.MaxDate.Unix();
	}

	return t.partitioner
}

func (t *TestService) getCreateTableQuery(name string) string {
	var query string

	tool := t.getPartitionTool()

	query = fmt.Sprintf("CREATE TABLE `%s` (\n", name);
	query += "`id` int(10) unsigned NOT NULL AUTO_INCREMENT,\n";
	query += "`event_id` char(36) COLLATE utf8_unicode_ci NOT NULL,\n";
	query += "`event_date` datetime(6) NOT NULL,\n";
	query += "PRIMARY KEY (`id`,`event_date`)\n";
	query += ") ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci\n"
	query += tool.GetPartitionsQueryString()

	return query
}

func (t *TestService) dropTable(name string) error {

	var query string
	query = fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", name)

	if t.Host != DOCKER_DB_HOST {
		msg := "Drop table is allowed for docker env only,"
		msg += "you can run query manualy if you know what you are doing: \n%s\n"
		fmt.Printf(msg, query)
		return nil
	}

	conn, err := t.GetConnection()

	if err != nil {
		return err
	}

	_, err = conn.Exec(query)

	return err
}


func (t *TestService) CreateTable(name string) error {

	fmt.Printf("Dropping table %s if exists\n", name)

	err := t.dropTable(name)

	if err != nil {
		return err
	}

	query := t.getCreateTableQuery(name)
	fmt.Printf("Processing query:\n%s\n", query)

	conn, err := t.GetConnection()

	if err != nil {
		return err
	}

	_, err = conn.Exec(query)

	fmt.Println("Table created.")

	return err
}


func (t *TestService) InsertRecord(table string) error {

	if t.insertTsMin == 0 {
		t.getPartitionTool()
	}

	uu := uuid.NewV4()

	tsUnix := random(t.insertTsMin, t.insertTsMax)
	ts := time.Unix(tsUnix, 0)

	query := fmt.Sprintf("INSERT INTO `%s` (event_id, `%s`) VALUES(?, ?)",
		table,
		t.getPartitionTool().Column);

	conn, err := t.GetConnection()

	if err != nil {
		return err
	}

	_, err = conn.Exec(query, uu, ts.Format("2006-01-02 15:04:05"))

	if err != nil {
		return err
	}

	return nil
}

func random(min, max int64) int64 {
	rand.Seed(time.Now().Unix())
	return rand.Int63n(max - min) + min
}
