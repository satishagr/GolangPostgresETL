package main

import (
	"GoEtl/models"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

//structure for data with Db column names
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
	FEAT_CODE   int64     `db:"FEAT_CODE"`
	GROUNDELEV  float64   `db:"GROUNDELEV"`
	GEOM_SOURCE string    `db:"GEOM_SOURCE"`
}

//Env struct for db pointer
type Env struct {
	db models.Datastore
}

//Main function acting Controller
func main() {
	//postgres://username:password@host/DbName?sslmode=disable
	db, err := models.NewDB("postgres://postgres:Johncena1@localhost/buildings?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	env := &Env{db}

	http.HandleFunc("/getall", env.getAllExec)
	http.HandleFunc("/get", env.getBuildingExec)
	http.HandleFunc("/loaddata", env.loadDataExec)
	http.HandleFunc("/", env.homeExec)
	http.ListenAndServe(":8080", nil)
}

type PageVariables struct {
	PageTitle string
}

//Handle to execute index.html
func (env *Env) homeExec(w http.ResponseWriter, r *http.Request) {
	HomePageVars := PageVariables{
		PageTitle: "myApp",
	}

	t, err := template.ParseFiles("views/index.html") //parse the html file index.html
	if err != nil {                                   // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

//Handle to load data into Db
// row is the no. of rows you want to fetch. The exeution time depends on the system capacity. '0' is for all rows.
func (env *Env) loadDataExec(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	rows, ok := r.URL.Query()["row"]
	var row string = ""
	if !ok || len(rows[0]) < 1 {
		row = "0"
	} else {
		row = rows[0]
	}
	rowCnt, err := env.db.LoadData(row)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "%d Rows Loaded\n", rowCnt)
}

//Handle to fetch all records from Db
func (env *Env) getAllExec(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	bds, err := env.db.GetAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	for _, bd := range bds {
		fmt.Fprintf(w, "%s, %s, %f, %s, %s, %s, %f, %f, %f, %f, %f, %s\n", bd.BASE_BBL, bd.MPLUTO_BBL, bd.BIN, bd.NAME, bd.LSTMODDATE, bd.LSTSTATTYPE,
			bd.CNSTRCT_YR, bd.DOITT_ID, bd.HEIGHTROOF, bd.FEAT_CODE, bd.GROUNDELEV, bd.GEOM_SOURCE)
	}
}

//Handle to filter based on parameter and it's value.
//get?parm=doitt_id,10.36
func (env *Env) getBuildingExec(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	temp := r.URL.Query()
	parm_map := strings.Split(temp.Get("parm"), ",")
	bds, err := env.db.GetBuilding(parm_map[0], parm_map[1])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	for _, bd := range bds {
		fmt.Fprintf(w, "%s, %s, %f, %s, %s, %s, %f, %f, %f, %f, %f, %s\n", bd.BASE_BBL, bd.MPLUTO_BBL, bd.BIN, bd.NAME, bd.LSTMODDATE, bd.LSTSTATTYPE,
			bd.CNSTRCT_YR, bd.DOITT_ID, bd.HEIGHTROOF, bd.FEAT_CODE, bd.GROUNDELEV, bd.GEOM_SOURCE)
	}
}
