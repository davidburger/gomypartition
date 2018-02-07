package Factory

import (
	"gomypartition/Application/Service"
	"gopkg.in/urfave/cli.v1"
	"gomypartition/Application/Database"
)

func NewMaintenanceService(db *Database.DbDriver, c *cli.Context) (s Service.MaintenanceService, err error) {

	s.Db = db

	s.TableName = c.String("table")
	s.Column = c.String("column")
	s.MaxPartitionCount = c.Int("max-partitions")
	s.RangeInDays = c.Int("range")
	s.OldPartitionRetentionDays = c.Int("retention")
	s.PartitionNamePrefix = c.String("prefix")
	s.DryRun = c.Bool("dry-run")

	err = Service.CheckRequiredFields(s, []string{"TableName", "Column"})

	return
}
