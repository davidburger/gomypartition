package Model

import (
	"fmt"
	"database/sql"
)

type PartitionInfo struct {
	Name string
	SubName sql.NullString
	Position int
	Method string
	Expression sql.NullString
	Description sql.NullString
	Rows int
	AvgRowLength int
	DataLength int
	Created string
	Updated sql.NullString
}

func (m *PartitionInfo) GetQueryColumns() ([]string) {

	return []string{"PARTITION_NAME",
	"SUBPARTITION_NAME",
	"PARTITION_ORDINAL_POSITION",
	"PARTITION_METHOD",
	"PARTITION_EXPRESSION",
	"PARTITION_DESCRIPTION",
	"TABLE_ROWS",
	"AVG_ROW_LENGTH",
	"DATA_LENGTH",
	"CREATE_TIME",
	"UPDATE_TIME"}
}


func (m *PartitionInfo) Dump() {
	fmt.Println("----")
	fmt.Printf("Name:			%s\n", m.Name);

	if m.SubName.Valid {
		fmt.Printf("Subname:		%s\n", m.SubName.String)
	}
	fmt.Printf("Position:		%d\n", m.Position)
	fmt.Printf("Method:			%s\n", m.Method)

	if m.Expression.Valid {
		fmt.Printf("Expression:		%s\n", m.Expression.String)
	}

	if m.Description.Valid {
		fmt.Printf("Description:		%s\n", m.Description.String)
	}

	fmt.Printf("Rows:			%d\n", m.Rows)
	fmt.Printf("AvgRowLength:		%d\n", m.AvgRowLength)
	fmt.Printf("DataLength: 		%d\n", m.DataLength)
	fmt.Printf("Created:		%s\n", m.Created)

	if m.Updated.Valid {
		fmt.Printf("Updated:		%s\n", m.Updated.String)
	}
}
