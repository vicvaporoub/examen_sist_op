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

	defer rows.Close()
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
	//DSN (Data Source Name) para conectarse a la base de datos sql
	dsn := "root:root@tcp(mysql-container:3306)/university"
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	//Verificar conexión
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Conexión exitosa a la base de datos")

	http.HandleFunc("/alumnos", func(w http.ResponseWriter, r *http.Request) {
		// Cabeceras CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Preflight
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
