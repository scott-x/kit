package msql

import (
	"database/sql"
	"errors"
)

func HandleStmtExec(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected > 0 {
		return nil
	}
	return errors.New("0 row was affected")
}
