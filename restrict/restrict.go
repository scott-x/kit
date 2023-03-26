package restrict

import (
	"fmt"
	"strings"
)

type RestrictKey int

// restrict for one field in database
type Restrict struct {
	RestrictKey `json:"restrict_key"`
	Value       interface{} `json:"value"`
}

type DBField struct {
	Name      string     `json:"name"` //database field
	Restricts []Restrict `json:"restricts"`
}

const (
	EQ     RestrictKey = iota //equal
	LT                        //less than
	GT                        //greater than
	LIKE                      //like '%%'
	L_Like                    //left like
	R_Like                    //right like
	LOCATE                    //locate
)

func (field *DBField) AddRestrict(key RestrictKey, value interface{}) {
	//nil pointer will throw error
	if field == nil {
		return
	}
	field.Restricts = append(field.Restricts, Restrict{
		RestrictKey: key,
		Value:       value,
	})
}

// the praram get from gin is string only
func NewDBField(arg string, value interface{}) *DBField {
	switch res := value.(type) {
	case int:
		if res == 0 {
			return nil
		}
	case string:
		if strings.TrimSpace(res) == "" {
			return nil
		}
	}
	return &DBField{
		Name: arg,
	}
}

func NewDBFieldWithSingleRestrict(param string, value interface{}, key RestrictKey) *DBField {
	switch res := value.(type) {
	case int:
		if res == 0 {
			return nil
		}
	case string:
		if strings.TrimSpace(res) == "" {
			return nil
		}
	}
	field := NewDBField(param, value)
	field.AddRestrict(key, value)
	return field
}

func HandleDbFiles(db_fields []*DBField) ([]interface{}, string) {
	var num_of_valid_condition int //the number of valid conditions
	var queryParams []interface{}  //query params that stmt.Query() / stmt.QueryRow() will use
	var condition string           // sql sentence: eg where balabala....

	for _, field := range db_fields {
		if field == nil { //front may pass null data
			continue
		}
		if num_of_valid_condition == 0 {
			condition += " where "
		}
		fieldName := field.Name
		for _, restrict := range field.Restricts {
			key := restrict.RestrictKey
			value := restrict.Value
			switch key {
			case EQ:
				if num_of_valid_condition == 0 {
					condition += fmt.Sprintf(" %s=?", fieldName)
				} else {
					condition += fmt.Sprintf(" and %s=?", fieldName)
				}
			case LT:
				if num_of_valid_condition == 0 {
					condition += fmt.Sprintf(" %s<=?", fieldName)
				} else {
					condition += fmt.Sprintf(" and %s<=?", fieldName)
				}
			case GT:
				if num_of_valid_condition == 0 {
					condition += fmt.Sprintf(" %s>=?", fieldName)
				} else {
					condition += fmt.Sprintf(" and %s>=?", fieldName)
				}
			case LIKE:
				if num_of_valid_condition == 0 {
					condition += fieldName + ` like CONCAT('%',?,'%')`
				} else {
					condition += " and " + fieldName + ` like CONCAT('%',?,'%')`
				}
			case R_Like:
				if num_of_valid_condition == 0 {
					condition += fieldName + ` like CONCAT('%',?)`
				} else {
					condition += " and " + fieldName + ` like CONCAT('%',?)`
				}
			case L_Like:
				if num_of_valid_condition == 0 {
					condition += fieldName + ` like CONCAT(?,'%')`
				} else {
					condition += " and " + fieldName + ` like CONCAT(?,'%')`
				}
			case LOCATE:
				if num_of_valid_condition == 0 {
					condition += fmt.Sprintf(` locate(%s , %s)>0 `, value.(string), fieldName)
				} else {
					condition += fmt.Sprintf(` and locate(%s , %s)>0 `, value.(string), fieldName)
				}
			}
			num_of_valid_condition++
			queryParams = append(queryParams, value)
		}
	}

	return queryParams, condition
}
