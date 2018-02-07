package Action

import (
	"gomypartition/Application/Service"
)

type MaintenanceAction struct {
	infoService Service.InfoService
	maintenanceService Service.MaintenanceService
}

func (a *MaintenanceAction) Create(infoService Service.InfoService, maintenanceService Service.MaintenanceService) {
	a.infoService = infoService
	a.maintenanceService  = maintenanceService
}


func (a *MaintenanceAction) Process() error {

	partitions, err := a.infoService.GetPartitionInfo()

	defer a.infoService.Db.CloseConnection()

	if err != nil {
		return err
	}

	a.maintenanceService.SetPartitions(partitions)
	return a.maintenanceService.Process()
}
