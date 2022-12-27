package main

import(
	"database/sql"
	"net/http"
	"log"
	"text/template"
	_"github.com/go-sql-driver/mysql"
)

func bdConnection() *sql.DB {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "gocrud"
	connection, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/" + Nombre)
	if err != nil {
		panic(err.Error())
	}
	return connection
}

var templates = template.Must(template.ParseGlob("templates/*"))

func main(){
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/insert", insertEmployee)
	http.HandleFunc("/remove", removeEmployee)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/update", updateEmployee)

	log.Println("Running server on port 8000")
	http.ListenAndServe(":8000", nil)
}

func updateEmployee(w http.ResponseWriter, r *http.Request){
	if(r.Method != "POST"){
		http.Error(w, "404 not found.", http.StatusNotFound)
	}
	bdConnectionOK := bdConnection()
	id:=r.FormValue("id")
	name:=r.FormValue("name")
	email:=r.FormValue("email")
		
	insertRecords,err := bdConnectionOK.Prepare("UPDATE employee SET name=?, email=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	
	insertRecords.Exec(name, email, id)
	
	http.Redirect(w, r, "/", 301)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	idEmployee := r.URL.Query().Get("id")
	if r.Method != "GET"{
		http.Error(w, "404 not found.", http.StatusNotFound)
	}
	bdConnectionOK := bdConnection()
	employeeRecord,err := bdConnectionOK.Query("SELECT * FROM employee WHERE id=?",idEmployee)
	if err != nil {
		panic(err.Error())
	}
	employee:=Employee{}
	var id, name, email string
	employeeRecord.Next()
	err =  employeeRecord.Scan(&id, &name, &email)
	if err != nil{
		panic(err.Error())
	}
	employee.Id = id
	employee.Name = name
	employee.Email = email
	templates.ExecuteTemplate(w, "edit", employee)
	
}

func removeEmployee(w http.ResponseWriter, r * http.Request) {
	idEmployee := r.URL.Query().Get("id")
	bdConnectionOK := bdConnection()
	query,err := bdConnectionOK.Prepare("DELETE FROM employee WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	query.Exec(idEmployee)
	http.Redirect(w, r, "/", 301)
}

func homeHandler(w http.ResponseWriter, r * http.Request){

	if r.Method != "GET"{
		http.Error(w, "404 not found.", http.StatusNotFound)
	}
	bdConnectionOK := bdConnection()
	employeeRecords,err := bdConnectionOK.Query("SELECT * FROM employee")
	if err != nil {
		panic(err.Error())
	}
	employee:=Employee{} // variable de tipo empleado
	employees := []Employee{} // array del tipo empleado
	for employeeRecords.Next(){
		var id, name, email string
		err =  employeeRecords.Scan(&id, &name, &email)
		if err != nil{
			panic(err.Error())
		}
		employee.Id = id
		employee.Name = name
		employee.Email = email
		employees = append(employees, employee)
	}
	templates.ExecuteTemplate(w, "inicio", employees)
}

func createHandler(w http.ResponseWriter, r * http.Request) {
	if r.Method != "GET"{
		http.Error(w, "404 not found.", http.StatusNotFound)
	}	
	templates.ExecuteTemplate(w, "create", nil)
}

func insertEmployee(w http.ResponseWriter, r * http.Request){

	if r.Method != "POST"{
		http.Error(w, "404 not found.", http.StatusNotFound)
	}
	bdConnectionOK := bdConnection()
	name:=r.FormValue("name")
	email:=r.FormValue("email")
		
	insertRecords,err := bdConnectionOK.Prepare("INSERT INTO employee(name, email) VALUES (?, ?)")
	if err != nil {
		panic(err.Error())
	}
	
	insertRecords.Exec(name, email)
	
	http.Redirect(w, r, "/", 301)	
}

type Employee struct {
	Id string
	Name string
	Email string
}