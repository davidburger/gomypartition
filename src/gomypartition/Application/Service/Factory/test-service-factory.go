package Factory

import (
	"gomypartition/Application/Service"
	"os"
	"strconv"
	"log"
)

func NewTestService() (s Service.TestService) {
	var err error

	s.Host = os.Getenv("MYSQL_HOST")
	s.User = os.Getenv("MYSQL_USER")
	s.Password = os.Getenv("MYSQL_PASSWORD")
	s.Port, err = strconv.Atoi(os.Getenv("MYSQL_PORT"))
	s.Database = os.Getenv("MYSQL_SCHEMA")

	if err != nil {
		log.Fatal(err)
	}

	return
}
