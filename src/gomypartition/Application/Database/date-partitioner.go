package Database

import (
	"time"
	"fmt"
	"strings"
	"gomypartition/Application/Model"
)

const defaultMaxPartitions = 50
const defaultRangeDays = 30
const defaultPrefix = "to"
const maxPartitionName = "MAX"
const partitionDateFormat = "20060102"
const DatabaseDateFormat = "2006-01-02"

type DatePartitioner struct {
	Column string
	MinDate time.Time
	MaxDate time.Time
	MaxPartitions int
	RangeDays int
	Prefix string
	RetentionDays int
	CurrentPartitionCount int
	Initialized bool
}

func (d *DatePartitioner) getMaxPartitions() int {
	if d.MaxPartitions != 0 {
		return d.MaxPartitions
	}

	return defaultMaxPartitions
}

func (d *DatePartitioner) getPrefix() string {
	if "" != d.Prefix  {
		return d.Prefix
	}

	return defaultPrefix
}

func (d *DatePartitioner) getRangeDays() int {
	if 0 != d.RangeDays {
		return d.RangeDays
	}

	return defaultRangeDays
}

func (d *DatePartitioner) GetMaxPartitionName() string {
	return fmt.Sprintf("%s_%s",
		d.getPrefix(),
		maxPartitionName)
}

func (d *DatePartitioner) GetDatePartitionName(date time.Time) string {
	return fmt.Sprintf("%s_%s",
		d.getPrefix(),
		date.Format(partitionDateFormat))
}

func (d *DatePartitioner) GetDateFromPartitionName(name string) (time.Time, error) {
	fullPrefix := fmt.Sprintf("%s_", d.getPrefix())
	dateString := strings.TrimPrefix(name, fullPrefix)

	return time.Parse(partitionDateFormat, dateString)
}

func (d *DatePartitioner) GetPartitionSql(date time.Time) string {
	return fmt.Sprintf("PARTITION %s VALUES LESS THAN (TO_DAYS('%s')) ENGINE = InnoDB",
		d.GetDatePartitionName(date),
		date.Format(DatabaseDateFormat));
}

func (d *DatePartitioner) GetNewPartitionFromDate(date time.Time) (partition Model.NewPartition) {

	partition.Name = d.GetDatePartitionName(date)
	partition.Sql = d.GetPartitionSql(date)

	return
}

func (d *DatePartitioner) GetPartitionsQueryString() string {
	var query string

	query = fmt.Sprintf("PARTITION BY RANGE (TO_DAYS(%s))\n", d.Column)
	query += "("

	maxPartitions := d.getMaxPartitions()

	var currentDate time.Time
	if !d.MinDate.IsZero() {
		currentDate = d.MinDate
	} else {
		currentDate = time.Now()
	}

	var maxDate time.Time
	if !d.MaxDate.IsZero() {
		maxDate = d.MaxDate
	} else {
		maxDate = time.Now().AddDate(10, 0, 0);
	}

	rangeDays := d.getRangeDays()

	var i = 0

	for x := d.CurrentPartitionCount; x < maxPartitions; x++ {
		date := currentDate.Add(time.Hour * time.Duration(24 * rangeDays))

		if date.After(maxDate) {
			break
		}

		query += d.GetPartitionSql(date)
		i++

		if i > 0 {
			query += ",\n"
		}

		currentDate = date
	}

	query += fmt.Sprintf("PARTITION %s VALUES LESS THAN MAXVALUE ENGINE = InnoDB)",
		d.GetMaxPartitionName());

	return query
}