package helper

import "database/sql"

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func CommitAndRollbackError(tx *sql.Tx) {
	err := recover()
	if err != nil {
		err := tx.Rollback()
		PanicIf(err)
		panic(err)
	} else {
		err := tx.Commit()
		PanicIf(err)
	}
}
