package Factory

import (
	"gomypartition/Application/Action"
	"gomypartition/Application/Service/Factory"
	"gomypartition/Application/Database"
	"gopkg.in/urfave/cli.v1"
)

func NewMaintenanceAction(db *Database.DbDriver, c *cli.Context) (a Action.MaintenanceAction, err error) {
	infoService, err := Factory.NewInfoService(db, c)

	if err != nil {
		return
	}

	maintenanceService, err := Factory.NewMaintenanceService(db, c)

	if err == nil {
		a.Create(infoService, maintenanceService)
	}
	return
}
