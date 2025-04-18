package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func getAlumnos(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT name FROM student LIMIT 5")

	if err != nil {
		return nil, err
	}

	defer rows.Close() // Asegurarnos que los resultados de la consulta se cierren cuando la consulta acaba
	var alumnos []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		alumnos = append(alumnos, name)
	}
	return alumnos, nil
}

func studentHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	alumnos, err := getAlumnos(db)
	if err != nil {
		panic(err)
	}

	for i, alumno := range alumnos {
		fmt.Fprintln(w, "Alumno ", (i + 1), ": ", alumno)
	}

}

func main() {
	//DSN (Data Source Name) para conectarse a la base de datos
	dsn := "root:root@tcp(mysql-container:3306)/university"
	db, err := sql.Open("mysql", dsn) //Intenta abrir una conexión con university

	if err != nil {
		panic(err)
	}

	defer db.Close() //Aqui cierra la conexión cuando termina el programa

	//Verificar conexión
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Conexión exitosa a la base de datos")

	http.HandleFunc("/alumnos", func(w http.ResponseWriter, r *http.Request) {
		// Cabeceras CORS
		w.Header().Set("Access-Control-Allow-Origin", "*") //Permite que cualquier origen pueda hacer solicitudes
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") //Aqui ponemos que solicitudes están permitidas
		//GET para poder pedir los datos al servidor 
		//y Options es una petición que el navegador manda antes del GET para verificar si el servidor acepta ciertas reglas
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type") //

		// Preflight
		//r.Method dice que tipo de solicitud llegó al servidor, entonces si la solicitud que recibió es de tipo OPTIONS
		//Manda un código 200 de Ok
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		studentHandler(db, w, r)
		
	})

	// Iniciar el servidor
	fmt.Println("Servidor iniciado en http://localhost:3000/alumnos")
	err3 := http.ListenAndServe(":3000", nil)
	if err3 != nil {
		fmt.Println("Error iniciando el servidor:", err)
	}

}
