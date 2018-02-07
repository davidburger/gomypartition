package Service

import (
	"gomypartition/Application/Database"
	"gomypartition/Application/Model"
	"log"
	"fmt"
	"time"
	"errors"
)

type MaintenanceService struct {
	Db *Database.DbDriver
	TableName string
	Column string
	MaxPartitionCount int
	RangeInDays int
	OldPartitionRetentionDays int
	PartitionNamePrefix string
	DryRun bool
	partitions []Model.PartitionInfo
}

func (s *MaintenanceService) SetPartitions(partitions []Model.PartitionInfo) {
	s.partitions = partitions
}

func (s *MaintenanceService) getPartitionTool() (p Database.DatePartitioner) {
	p.Column = s.Column
	p.MaxPartitions = s.MaxPartitionCount
	p.RangeDays = s.RangeInDays
	p.Prefix = s.PartitionNamePrefix
	p.RetentionDays = s.OldPartitionRetentionDays

	return
}

func (s *MaintenanceService) Process() error {
	partitionCount := len(s.partitions)
	if partitionCount == 0 {
		log.Fatal(fmt.Sprintf("Table `%s` has no partition", s.TableName))
	}

	fmt.Printf("Found %d partitions.\n", partitionCount)

	dateFormat := Database.DatabaseDateFormat
	currentDate, err := time.Parse(dateFormat, time.Now().Format(dateFormat));

	if err != nil {
		return errors.New("Error while parsing current date")
	}

	ptool := s.getPartitionTool()

	var startDate time.Time
	var partitionsToDrop[]Model.PartitionInfo

	lastPartitionName := ptool.GetMaxPartitionName()

	for _, partitionInfo := range(s.partitions) {
		name := partitionInfo.Name

		if name == lastPartitionName {
			continue
		}

		partitionDate, err := ptool.GetDateFromPartitionName(name)

		if err != nil {
			return err
			//return errors.New(fmt.Sprintf("Error while parsing date from partition %s", name))
		}

		diff := currentDate.Sub(partitionDate)
		diffDays := int(diff.Hours() / 24)

		if (diffDays > s.OldPartitionRetentionDays) {
			partitionsToDrop = append(partitionsToDrop, partitionInfo)
			fmt.Printf("Partition %s will be removed\n", name)
		} else {
			startDate = partitionDate
		}
	}

	countWithoutDeleted := partitionCount - len(partitionsToDrop)

	date := startDate

	var partitionsToAdd []Model.NewPartition

	for x := countWithoutDeleted; x < s.MaxPartitionCount; x++ {
		date = date.AddDate(0, 0, s.RangeInDays)

		if !date.After(currentDate) {
			continue
		}

		partition := ptool.GetNewPartitionFromDate(date)
		partitionsToAdd = append(partitionsToAdd, partition)
	}

	return s.processQueries(&ptool, partitionsToDrop, partitionsToAdd)
}

func (s *MaintenanceService) processQueries(ptool *Database.DatePartitioner,
	toDrop []Model.PartitionInfo,
	toAdd []Model.NewPartition) error {

		dropLen := len(toDrop)
		addLen := len(toAdd)

		fmt.Printf("%d partitions will be dropped\n", dropLen)
		fmt.Printf("%d partitions will be created\n\n", addLen)

		conn, err := s.Db.GetConnection()

		if err != nil {
			return err
		}

		defer s.Db.CloseConnection()

		for _, drop := range (toDrop) {
			query := fmt.Sprintf("ALTER TABLE `%s` DROP PARTITION %s;", s.TableName, drop.Name)

			fmt.Println(query)
			if s.DryRun == false {
				_, err = conn.Exec(query)

				if err != nil {
					return err
				}

				fmt.Println("Partition dropped.")
			}
		}

		if addLen == 0 {
			return nil
		}

		maxName := ptool.GetMaxPartitionName()

		query := fmt.Sprintf("ALTER TABLE `%s` \n", s.TableName);
		query += fmt.Sprintf("REORGANIZE PARTITION %s INTO (\n", maxName)


		for _, padd := range(toAdd) {
			query += padd.Sql;
			query += ",\n"
		}

		query += fmt.Sprintf("PARTITION %s VALUES LESS THAN MAXVALUE);", maxName)

		fmt.Println("\n" + query)
		if s.DryRun == false {
			_, err = conn.Exec(query)

			if err != nil {
				return err
			}

			fmt.Println("Partitions reorganized.")
		}

		return nil
}
