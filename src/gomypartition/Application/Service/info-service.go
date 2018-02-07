package Service

import (
	"gomypartition/Application/Database"
	"fmt"
	"strings"
	"database/sql"
	"gomypartition/Application/Model"
)

const DefaultOrderField = "PARTITION_ORDINAL_POSITION"

type InfoService struct {
	Db *Database.DbDriver
	TableName string
	SortBy string
	cols []string
}

func (s *InfoService) getInfoQuery() string {
	var query, orderField string
	var model Model.PartitionInfo

	query = "SELECT \n";
	query += strings.Join(model.GetQueryColumns(), ",\n") + "\n"
	query += "FROM information_schema.PARTITIONS\n"
	query += "WHERE TABLE_SCHEMA = ?\n"
	query += "AND TABLE_NAME = ?\n"

	orderField = DefaultOrderField

	if s.SortBy != "" {
		orderField = s.SortBy
	}

	query += fmt.Sprintf("ORDER BY `%s`", orderField)

	return query
}


func (s *InfoService) GetPartitionInfo() (result []Model.PartitionInfo, err error) {
	var conn *sql.DB

	fmt.Printf("Gathering partition info for %s.%s.%s ...\n",
		s.Db.Host,
		s.Db.Database,
		s.TableName)

	conn, err = s.Db.GetConnection()

	if err != nil {
		return
	}

	stmt, err := conn.Prepare(s.getInfoQuery())

	if err != nil {
		return
	}

	defer stmt.Close()

	rows, err := stmt.Query(s.Db.Database, s.TableName)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {

		var info Model.PartitionInfo

		err = rows.Scan(&info.Name,
			&info.SubName,
			&info.Position,
			&info.Method,
			&info.Expression,
			&info.Description,
			&info.Rows,
			&info.AvgRowLength,
			&info.DataLength,
			&info.Created,
			&info.Updated)

		if err != nil {
			return
		}

		result = append(result, info)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}
