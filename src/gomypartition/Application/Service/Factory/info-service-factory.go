package Factory

import (
	"gomypartition/Application/Service"
	"gopkg.in/urfave/cli.v1"
	"gomypartition/Application/Database"
)


func NewInfoService(db *Database.DbDriver, c *cli.Context) (s Service.InfoService, err error) {

	s.Db = db

	s.TableName = c.String("table")
	s.SortBy = c.String("orderfld")

	err = Service.CheckRequiredFields(s, []string{"TableName"})

	return
}
