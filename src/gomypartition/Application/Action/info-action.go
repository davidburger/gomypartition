package Action

import (
	"gomypartition/Application/Service"
)

type InfoAction struct {
	infoService Service.InfoService
}

func (a *InfoAction) Create(service Service.InfoService) {
	a.infoService = service
}

func (a *InfoAction) Process() error {
	records, err := a.infoService.GetPartitionInfo()

	defer a.infoService.Db.CloseConnection()
	if err != nil {
		return err
	}

	for _, info := range records {
		info.Dump()
	}

	return nil
}
