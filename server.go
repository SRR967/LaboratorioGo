package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

// Función para servir la página principal
func handler(w http.ResponseWriter, r *http.Request) {
	// Plantillas HTML
	templates := []string{"plantilla1.html", "plantilla2.html"}

	// Semilla para generar aleatoriedad
	rand.Seed(time.Now().UnixNano())

	// Seleccionar una plantilla aleatoriamente
	selectedTemplate := templates[rand.Intn(len(templates))]

	// Parsear la plantilla seleccionada
	tmpl, err := template.ParseFiles(selectedTemplate)
	if err != nil {
		http.Error(w, "Error al cargar la plantilla", http.StatusInternalServerError)
		return
	}

	// Renderizar la plantilla en la respuesta
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error al renderizar la plantilla", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Asignar la función handler a la ruta "/"
	http.HandleFunc("/", handler)

	// Servir archivos estáticos (imágenes, CSS, etc.)
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Iniciar el servidor en el puerto 8080
	fmt.Println("Servidor web en ejecución en http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
