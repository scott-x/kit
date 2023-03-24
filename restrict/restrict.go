package restrict

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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
func NewDBField(c *gin.Context, arg string) *DBField {
	value := c.Query(arg)
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &DBField{
		Name: arg,
	}
}

func NewDBFieldWithSingleRestrict(c *gin.Context, param string, key RestrictKey, convert2int bool) *DBField {
	value := strings.TrimSpace(c.Query(param))
	if value == "" {
		return nil
	}
	field := NewDBField(c, param)
	if convert2int {
		_value, _ := strconv.Atoi(value)
		field.AddRestrict(key, _value)
	} else {
		field.AddRestrict(key, value)
	}
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
			}
			num_of_valid_condition++
			queryParams = append(queryParams, value)
		}
	}

	return queryParams, condition
}
