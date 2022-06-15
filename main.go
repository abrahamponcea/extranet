package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

//Todo lo que estre dentro de la carpeta plantillas se guatdara en la
//variable plantillas
var plantillas = template.Must(template.ParseGlob("plantillas/*"))

type Reserva struct {
	Id         int
	Nombre     string
	habitacion string
	numero     int
	precio     string
}

func main() {
	http.HandleFunc("/", Inicio) //Declarando funcion Inicio
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/update", Update)
	log.Println("Servidor trabajando...")
	http.ListenAndServe(":3000", nil)
}

func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contra := "root"
	Nombre := "extranet"

	conexion, err := sql.Open(Driver, Usuario+":"+Contra+"@tcp(127.0.0.1:3307)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("txtId")
		nombre := r.FormValue("txtNombreReserva")
		habitacion := r.FormValue("txtHabitacion")
		numero := r.FormValue("txtNumero")
		precio := r.FormValue("txtPrecio")

		conexion := conexionDB()
		modificar, err := conexion.Prepare("Update reservas SET nombreReserva=?, tipoHabitacion=?, numeroHabitacion=?, precioHabitacion=? Where id=?")
		if err != nil {
			panic(err.Error())
		}
		modificar.Exec(nombre, habitacion, numero, precio, id)
		http.Redirect(w, r, "/", 301)
	}
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idReserva := r.URL.Query().Get("id")
	conexion := conexionDB()
	editar, err := conexion.Query("Select * From reservas Where id=?", idReserva)
	_reserva := Reserva{}
	for editar.Next() {
		var id, numero int
		var nombre, habitacion, precio string
		err = editar.Scan(&id, &nombre, &habitacion, &numero, &precio)
		if err != nil {
			panic(err.Error())
		}
		_reserva.Id = id
		_reserva.Nombre = nombre
		_reserva.habitacion = habitacion
		_reserva.numero = numero
		_reserva.precio = precio
	}
	plantillas.ExecuteTemplate(w, "editar", _reserva)
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idReserva := r.URL.Query().Get("id")
	conexion := conexionDB()
	borrar, err := conexion.Prepare("Delete From reservas where id=?")
	if err != nil {
		panic(err.Error())
	}
	borrar.Exec(idReserva)
	http.Redirect(w, r, "/", 301)
}

//Ruta Para incertar Datos a BD
func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("txtNombre")
		habitacion := r.FormValue("txtHabitacion")
		numero := r.FormValue("txtNumero")
		precio := r.FormValue("txtPrecio")

		conexion := conexionDB()
		insertar, err := conexion.Prepare("insert into reservas(nombreReserva, tipoHabitacion, numeroHabitacion, precioHabitacion) values(?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insertar.Exec(nombre, habitacion, numero, precio)
		http.Redirect(w, r, "/", 301)
	}
}

func Inicio(w http.ResponseWriter, r *http.Request) {
	//platillas es el nombre de la carpeta+
	conexion := conexionDB()
	reserva := Reserva{}
	arregloReserva := []Reserva{} //arreglo de reservas
	//Seleccionando los datos de la base de datos para mostrar
	seleccionar, err := conexion.Query("Select * From reservas")
	if err != nil {
		panic(err.Error())
	}

	for seleccionar.Next() {
		var id, numero int
		var nombre, habitacion, precio string
		err = seleccionar.Scan(&id, &nombre, &habitacion, &numero, &precio)
		if err != nil {
			panic(err.Error())
		}
		reserva.Id = id
		reserva.Nombre = nombre
		reserva.habitacion = habitacion
		reserva.numero = numero
		reserva.precio = precio
		arregloReserva = append(arregloReserva, reserva)
	}

	//insertar.Exec()
	plantillas.ExecuteTemplate(w, "inicio", arregloReserva)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	//platillas es el nombre de la carpeta
	plantillas.ExecuteTemplate(w, "crear", nil)
}
