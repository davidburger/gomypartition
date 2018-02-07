package Service

import (
	"fmt"
	"reflect"
	"errors"
)

func getErrorForEmptyField(field string) error {
	return errors.New(fmt.Sprintf("Required option `%s` is empty, please check command help", field))
}

func CheckRequiredFields(s interface{}, requiredFields []string) error {
	r := reflect.ValueOf(s)
	for _, fieldName := range requiredFields {
		f := reflect.Indirect(r).FieldByName(fieldName)

		switch ptype := f.Kind(); ptype {
		case reflect.String:
			if f.String() == "" {
				return getErrorForEmptyField(fieldName)
			}
		case reflect.Int:
			if f.Int() == 0 {
				return getErrorForEmptyField(fieldName)
			}
		default:
			return errors.New(fmt.Sprintf("Unsupported type %s for option %s", ptype.String()))
		}
	}

	return nil
}
