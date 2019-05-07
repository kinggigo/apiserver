package service

import (
	"errors"
	db2 "github.com/kinggigo/secret/server/db"
)

func GetRowNum() (int, error) {
	db := db2.GormDB
	rows, err := db.Raw(`SELECT FOUND_ROWS() as nRows `).Rows()
	if err != nil {
		return 0, errors.New("테이블 조회 오류 : " + err.Error())
	}
	var nRows int
	for rows.Next() {
		err = rows.Scan(&nRows)
		if err != nil {
			return 0, errors.New("테이블 조회 오류 : " + err.Error())
		}
	}
	return nRows, nil
}
