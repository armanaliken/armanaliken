package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type server struct {
	db *sql.DB
}

type Car struct {
	Id        int
	Mark      string
	Model     string
	FirstName string
	LastName  string
	Price     string
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Arman!05"
	dbname   = "restyling"
)

func dbConnect() *server {
	dbconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbconn)
	fmt.Println("Opening database")
	if err != nil {
		log.Fatal("Open ", err)
	}

	//err = db.Ping()
	//if err != nil {
	//	log.Fatal(err)
	//}
	fmt.Println("Successfully connected to database")

	return &server{db: db}

}

func (s *server) cars(w http.ResponseWriter, r *http.Request) {
	var cars []Car
	res, err := s.db.Query("select * from cars;")
	if err != nil {
		log.Fatal("Query get car", err)
	}
	for res.Next() {
		var car Car
		res.Scan(&car.Id, &car.Mark, &car.Model, &car.FirstName, &car.LastName, &car.Price)
		cars = append(cars, car)
	}
	t, err := template.ParseFiles("./static/all.html")
	if err != nil {
		log.Fatal("parse", err)
	}
	t.Execute(w, cars)
}

func (s *server) new(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Fatal(err)
		}
		mark := r.FormValue("mark")
		model := r.FormValue("model")
		first := r.FormValue("first")
		last := r.FormValue("last")
		price := r.FormValue("price")
		_, _ = s.db.Exec("insert into cars(mark, model, firstname, lastname, price) values($1, $2, $3, $4, $5)", mark, model, first, last, price)
		http.Redirect(w, r, "/cars", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseFiles("./static/create.html")
	t.Execute(w, nil)
}

func (s *server) delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Fatal(err)
		}
		id := r.FormValue("id")
		_, _ = s.db.Exec("delete from cars where id=$1", id)
		http.Redirect(w, r, "/cars", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseFiles("./static/delete.html")
	t.Execute(w, nil)
}

func (s *server) update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Fatal(err)
		}
		id := r.FormValue("id")
		mark := r.FormValue("mark")
		model := r.FormValue("model")
		first := r.FormValue("first")
		last := r.FormValue("last")
		price := r.FormValue("price")
		_, _ = s.db.Exec("update cars set mark=$1, model=$2, firstname=$3, lastname=$4, price=$5 where id=$6", mark, model, first, last, price, id)
		http.Redirect(w, r, "/cars", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseFiles("./static/update.html")
	t.Execute(w, nil)
}

func main() {
	s := dbConnect()

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/cars", s.cars)
	http.HandleFunc("/new", s.new)
	http.HandleFunc("/delete", s.delete)
	http.HandleFunc("/update", s.update)
	defer s.db.Close()
	http.ListenAndServe(":8080", nil)
}
