package models

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var index int = 0
var rowsLoaded int = 0

type Building struct {
	BASE_BBL    string    `db:"BASE_BBL"`
	MPLUTO_BBL  string    `db:"MPLUTO_BBL"`
	BIN         float64   `db:"BIN"`
	NAME        string    `db:"NAME"`
	LSTMODDATE  time.Time `db:"LSTMODDATE"`
	LSTSTATTYPE string    `db:"LSTSTATTYPE"`
	CNSTRCT_YR  float64   `db:"CNSTRCT_YR"`
	DOITT_ID    float64   `db:"DOITT_ID"`
	HEIGHTROOF  float64   `db:"HEIGHTROOF"`
	FEAT_CODE   float64   `db:"FEAT_CODE"`
	GROUNDELEV  float64   `db:"GROUNDELEV"`
	GEOM_SOURCE string    `db:"GEOM_SOURCE"`
}

func (db *DB) LoadData(row string) (int64, error) {
	maxcnt, err := strconv.Atoi(row)
	if err != nil {
		return 0, err
	}
	if maxcnt > 500 {
		return 0, err
	}
	maxcnt = maxcnt + rowsLoaded
	var cnt int = 0
	csvFile, _ := os.Open("building.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var first = true
	var buildings []Building
	for {
		line, error := reader.Read()
		if first {
			first = false
			continue
		}
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		if maxcnt != 0 {
			if cnt == maxcnt {
				break
			} else if rowsLoaded != 0 && cnt < rowsLoaded {
				cnt = cnt + 1
				continue
			}
		}
		var BASE_BBL = line[0]
		var MPLUTO_BBL = line[1]
		var BIN, err = strconv.ParseFloat(line[2], 64)
		if err != nil {
			return 0, err
		}
		var NAME = line[3]
		const dateformat = "01/02/2006 03:04:05 PM -0700"
		var LSTMODDATE, err1 = time.Parse(dateformat, line[4])
		if err1 != nil {
			log.Panic(err1)
		}
		var LSTSTATTYPE = line[5]
		var CNSTRCT_YR, err2 = strconv.ParseFloat(line[6], 64)
		if err2 != nil {
			log.Panic(err2)
		}
		var DOITT_ID, err3 = strconv.ParseFloat(line[7], 64)
		if err3 != nil {
			log.Panic(err3)
		}
		var HEIGHTROOF, err4 = strconv.ParseFloat(line[8], 64)
		if err4 != nil {
			log.Panic(err4)
		}
		var FEAT_CODE, err5 = strconv.ParseFloat(line[9], 64)
		if err5 != nil {
			log.Panic(err5)
		}
		var GROUNDELEV, err6 = strconv.ParseFloat(line[10], 64)
		if err6 != nil {
			log.Panic(err6)
		}
		var GEOM_SOURCE = line[11]

		buildings = append(buildings, Building{
			BASE_BBL:    BASE_BBL,
			MPLUTO_BBL:  MPLUTO_BBL,
			BIN:         BIN,
			NAME:        NAME,
			LSTMODDATE:  LSTMODDATE,
			LSTSTATTYPE: LSTSTATTYPE,
			CNSTRCT_YR:  CNSTRCT_YR,
			DOITT_ID:    DOITT_ID,
			HEIGHTROOF:  HEIGHTROOF,
			FEAT_CODE:   FEAT_CODE,
			GROUNDELEV:  GROUNDELEV,
			GEOM_SOURCE: GEOM_SOURCE,
		})
		cnt = cnt + 1
	}
	sqlStr := "INSERT INTO building(BASE_BBL, MPLUTO_BBL, BIN, NAME, LSTMODDATE, LSTSTATTYPE, CNSTRCT_YR, DOITT_ID, HEIGHTROOF, FEAT_CODE, GROUNDELEV, GEOM_SOURCE) VALUES "
	vals := []interface{}{}
	for _, row := range buildings {
		sqlStr += "($" + getIndex() + ", $" + getIndex() + ", $" + getIndex() + ", $" + getIndex() + ", $" + getIndex() + ", $" + getIndex() + ", $" + getIndex() + ", $" + getIndex() + ", $" + getIndex() + ", $" + getIndex() + ", $" + getIndex() + ", $" + getIndex() + "),"
		vals = append(vals, row.BASE_BBL, row.MPLUTO_BBL, row.BIN,
			row.NAME, row.LSTMODDATE, row.LSTSTATTYPE,
			row.CNSTRCT_YR, row.DOITT_ID, row.HEIGHTROOF,
			row.FEAT_CODE, row.GROUNDELEV, row.GEOM_SOURCE)
	}
	index = 0
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	//prepare the statement
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	//format all vals at once
	res, err := stmt.Exec(vals...)
	if err != nil {
		return 0, err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	rowsLoaded = rowsLoaded + int(rowCnt)
	return rowCnt, nil
}

func getIndex() string {
	index = index + 1
	return strconv.Itoa(index)
}

func (db *DB) GetAll() ([]*Building, error) {
	rows, err := db.Query("SELECT * FROM building")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bds := make([]*Building, 0)
	for rows.Next() {
		bd := new(Building)
		err := rows.Scan(&bd.BASE_BBL, &bd.MPLUTO_BBL, &bd.BIN, &bd.NAME, &bd.LSTMODDATE, &bd.LSTSTATTYPE,
			&bd.CNSTRCT_YR, &bd.DOITT_ID, &bd.HEIGHTROOF, &bd.FEAT_CODE, &bd.GROUNDELEV, &bd.GEOM_SOURCE)
		if err != nil {
			return nil, err
		}
		bds = append(bds, bd)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return bds, nil
}

func (db *DB) GetBuilding(parm string, val string) ([]*Building, error) {
	sqlStr := "SELECT * FROM building where "
	sqlStr += parm + " = $1"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if strings.ToLower(parm) == strings.ToLower("BIN") || strings.ToLower(parm) == strings.ToLower("CNSTRCT_YR") ||
		strings.ToLower(parm) == strings.ToLower("DOITT_ID") || strings.ToLower(parm) == strings.ToLower("HEIGHTROOF") ||
		strings.ToLower(parm) == strings.ToLower("GROUNDELEV") || strings.ToLower(parm) == strings.ToLower("FEAT_CODE") {
		val_f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			log.Panic(err)
		}
		rows, err = stmt.Query(val_f)
		if err != nil {
			return nil, err
		}
	} else if strings.ToLower(parm) == strings.ToLower("LSTMODDATE") {
		return db.GetAll()
	} else {
		val_s := val
		rows, err = stmt.Query(val_s)
		if err != nil {
			return nil, err
		}
	}

	bds := make([]*Building, 0)
	for rows.Next() {
		bd := new(Building)
		err := rows.Scan(&bd.BASE_BBL, &bd.MPLUTO_BBL, &bd.BIN, &bd.NAME, &bd.LSTMODDATE, &bd.LSTSTATTYPE,
			&bd.CNSTRCT_YR, &bd.DOITT_ID, &bd.HEIGHTROOF, &bd.FEAT_CODE, &bd.GROUNDELEV, &bd.GEOM_SOURCE)
		if err != nil {
			return nil, err
		}
		bds = append(bds, bd)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return bds, nil
}
