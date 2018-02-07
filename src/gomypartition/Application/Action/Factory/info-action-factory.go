package Factory

import (
	"gomypartition/Application/Action"
	"gomypartition/Application/Service/Factory"
	"gomypartition/Application/Database"
	"gopkg.in/urfave/cli.v1"
)

func NewInfoAction(db *Database.DbDriver, c *cli.Context) (a Action.InfoAction, err error) {
	service, err := Factory.NewInfoService(db, c)

	if err == nil {
		a.Create(service)
	}

	return
}
