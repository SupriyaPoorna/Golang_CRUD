package main

import(
	"net/http"
	"html/template"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"	
	"log"
)

var tpl *template.Template

type User struct{
	Id int
	FName string
	LName string
	Email string
	Pwd string
}

func init(){
	tpl = template.Must(template.ParseGlob("html_template/*.html"))
}

//Opening the connection to the database
func db_conn()(db *sql.DB){
	db, err := sql.Open("mysql", "root:Password@(127.0.0.1:3306)/mydatabase")
	err = db.Ping()
  	 if err != nil {
	   log.Fatal(err)  
	}
	return db
}

// main is the entry point 
func main(){
	http.HandleFunc("/", index)
	http.HandleFunc("/nav_addUser",naviagteToAddpage)
	http.HandleFunc("/add_user",addUser)
	http.HandleFunc("/nav_displayUser",display_users)
	http.HandleFunc("/show",view_par_user)
	http.HandleFunc("/edit",edit_user)
	http.HandleFunc("/delete",delete_user)
	http.HandleFunc("/update_user",update_user)
	http.ListenAndServe(":8080",nil)
}

// function to display the home page
func index(w http.ResponseWriter, r *http.Request){
	tpl.ExecuteTemplate(w, "index_crud.html", nil)
} 

func naviagteToAddpage(w http.ResponseWriter, r *http.Request){
	tpl.ExecuteTemplate(w, "addUser_crud.html", nil)
}

func addUser(w http.ResponseWriter, r *http.Request){
	log.Println("inside add")
	if r.Method != "POST" {
		http.Redirect(w,r,"/",http.StatusSeeOther)
		return
	}
	
	usr := User{}

	usr.FName = r.FormValue("firstname")
	usr.LName = r.FormValue("lastname")
	usr.Email = r.FormValue("email")
	usr.Pwd = r.FormValue("password")
	
	db := db_conn()
	
	n, err := db.Prepare("insert into crud_new(firstName,lastName,email,password) values(?,?,?,?)")
	if err!=nil{
		log.Fatal(err)
	}else {
		n.Exec( usr.FName, usr.LName, usr.Email, usr.Pwd)	
		log.Println(n)
		display_users(w, r)	
	}

 	defer db.Close()	
}

func display_users(w http.ResponseWriter, r *http.Request){
	log.Println("inside display users")
	db := db_conn()
    selDB, err := db.Query("SELECT * FROM crud_new ORDER BY id ASC")
    if err != nil {
		log.Println("err")
		panic(err.Error())
		
    }
    usr := User{}
    res := []User{}
    for selDB.Next() {
        var id int
        var fname, lname, email, pwd string
        err = selDB.Scan(&id, &fname, &lname, &email, &pwd)
        if err != nil {
            panic(err.Error())
        }
        usr.Id = id
        usr.FName = fname
		usr.LName = lname
		usr.Email = email
        res = append(res, usr)
    }

	tpl.ExecuteTemplate(w, "displayUser_crud.html", res)
	defer db.Close()
	
}

func view_par_user(w http.ResponseWriter, r *http.Request){
	log.Println("inside view particular user details")

	nID := r.URL.Query().Get("id")
	usr := query_for_particular_row(nID)

	tpl.ExecuteTemplate(w, "showParticularUser_crud.html", usr)
}

func edit_user(w http.ResponseWriter, r *http.Request){
	log.Println("inside edit users ")
	
	nID := r.URL.Query().Get("id")	
	usr := query_for_particular_row(nID)
	
	tpl.ExecuteTemplate(w, "editUser_crud.html", usr)	
}

func update_user(w http.ResponseWriter, r *http.Request){
	log.Println("inside update users ")
	if r.Method != "POST" {
		http.Redirect(w,r,"/",http.StatusSeeOther)
		return
	}
	usr := User{}

	id :=  r.FormValue("id")
	usr.FName = r.FormValue("firstname")
	usr.LName = r.FormValue("lastname")
	usr.Email = r.FormValue("email")
	usr.Pwd = r.FormValue("password")
	
	db := db_conn()
	n, err := db.Prepare("UPDATE crud_new SET firstName=?, lastName=?,email=?, password =?  WHERE id=?")
	if err!=nil{
		log.Fatal(err)
	}else {
		n.Exec( usr.FName, usr.LName, usr.Email, usr.Pwd, id)	
		log.Println(n)
		display_users(w, r)
	}

 	defer db.Close()
}

//function to delete particular row based on primary key i.e ID
func delete_user(w http.ResponseWriter, r *http.Request){
	log.Println("inside delete users")
	db := db_conn()
	nID := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM crud_new WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(nID)
	log.Println("DELETE")
	display_users(w, r)
    defer db.Close()
}

// function to query for particular row in the database for the given ID.. returns single row
func query_for_particular_row(nId string)(usr1 User){
	
	usr := User{}
	
	var fname, lname, email, pwd string
	var id int
	db := db_conn()
	
	row := db.QueryRow("SELECT * FROM crud_new where id = ?", nId)
	err := row.Scan(&id, &fname, &lname, &email, &pwd)
	if err != nil && err == sql.ErrNoRows{
        panic(err)
	}else{
		usr.Id = id
		usr.FName = fname
		usr.LName = lname
		usr.Email = email
		usr.Pwd = pwd
		
		defer db.Close()
	}
	return usr
}