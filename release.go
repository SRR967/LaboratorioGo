package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Estructura que contiene la información de cada imagen
type ImageData struct {
	Base64 template.URL // URL amigable para HTML
	Name   string
}

// Datos para pasar a la plantilla
type PageData struct {
	HostName string
	Images   []ImageData
	Theme    string
}

// Variable fija para el nombre del host
var fixedHostName = "MiServidorFijo"

// Función para verificar si un archivo es una imagen
func isImage(fileName string) bool {
	imageExtensions := []string{".jpg", ".jpeg", ".png"}
	extension := strings.ToLower(filepath.Ext(fileName))

	for _, ext := range imageExtensions {
		if extension == ext {
			return true
		}
	}
	return false
}

// Función para leer las imágenes de un directorio y devolverlas codificadas en Base64
func getImageFiles(dirPath string) []ImageData {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Println("Error leyendo el directorio:", err)
		return nil
	}

	var imageFiles []ImageData
	for _, file := range files {
		if !file.IsDir() && isImage(file.Name()) {
			imagePath := filepath.Join(dirPath, file.Name())
			encodedImage := encodeImageToBase64(imagePath)
			imageFiles = append(imageFiles, ImageData{
				Base64: template.URL("data:image/jpeg;base64," + encodedImage),
				Name:   file.Name(),
			})
		}
	}
	return imageFiles
}

// Función para seleccionar n imágenes sin que se repitan
func selectNImages(imageFiles []ImageData, cantidad int) []ImageData {
	var imagesSelected []ImageData
	for len(imagesSelected) < cantidad {
		imagen := getRandomImage(imageFiles)
		if !imageInArray(imagesSelected, imagen) {
			imagesSelected = append(imagesSelected, imagen)
		}
	}
	return imagesSelected
}

// Función para verificar si una imagen ya existe en un arreglo de ImageData
func imageInArray(arr []ImageData, img ImageData) bool {
	for _, v := range arr {
		if v.Name == img.Name {
			return true
		}
	}
	return false
}

// Función para seleccionar una imagen al azar
func getRandomImage(imageFiles []ImageData) ImageData {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(imageFiles))
	return imageFiles[randomIndex]
}

// Función para leer y codificar una imagen en Base64
func encodeImageToBase64(imagePath string) string {
	imageData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		log.Println("Error leyendo el archivo de imagen:", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(imageData)
}

// Función para seleccionar una plantilla aleatoria
func randomTemplate() string {
	templates := []string{"plantilla1.html", "plantilla2.html"}
	rand.Seed(time.Now().UnixNano())
	return templates[rand.Intn(len(templates))]
}

// Función para obtener la subcarpeta (tema) especificada o una aleatoria
func getSubfolder(dir string, theme string) (string, string, error) {
	if theme != "" {
		subfolderPath := filepath.Join(dir, theme)
		if _, err := os.Stat(subfolderPath); os.IsNotExist(err) {
			return "", "", fmt.Errorf("el tema especificado no existe: %s", theme)
		}
		return subfolderPath, theme, nil
	} else {
		// Lee el contenido del directorio
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return "", "", err
		}

		// Filtra las subcarpetas
		var subfolders []fs.FileInfo
		for _, file := range files {
			if file.IsDir() {
				subfolders = append(subfolders, file)
			}
		}

		// Verifica si hay subcarpetas disponibles
		if len(subfolders) == 0 {
			return "", "", fmt.Errorf("no se encontraron subcarpetas en %s", dir)
		}

		// Genera un índice aleatorio y selecciona una subcarpeta
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(subfolders))

		// Obtiene la subcarpeta seleccionada
		subfolder := subfolders[randomIndex]
		subfolderName := subfolder.Name()
		subfolderPath := filepath.Join(dir, subfolderName)

		return subfolderPath, subfolderName, nil
	}
}

// Función para servir la página principal
func handler(imageDir string, theme string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		dirPath, themeName, err := getSubfolder(imageDir, theme)
		if err != nil {
			log.Printf("Error al obtener subcarpeta: %v", err)
			http.Error(w, "Error al obtener subcarpeta de imágenes.", http.StatusInternalServerError)
			return
		}

		imageFiles := getImageFiles(dirPath)
		if len(imageFiles) == 0 {
			http.Error(w, "No se encontraron imágenes en la carpeta.", http.StatusNotFound)
			return
		}

		hostName := fixedHostName

		// Selecciona 3 imágenes para mostrar
		selectedImages := selectNImages(imageFiles, 3)

		// Estructura que contiene los datos para la plantilla
		data := PageData{
			HostName: hostName,
			Images:   selectedImages,
			Theme:    themeName,
		}
		selectedTemplate := randomTemplate()

		// Parsear la plantilla seleccionada
		tmpl, err := template.ParseFiles(selectedTemplate)
		if err != nil {
			http.Error(w, "Error al cargar la plantilla", http.StatusInternalServerError)
			return
		}

		// Renderizar la plantilla en la respuesta
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error al renderizar la plantilla", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	// Definir las banderas
	port := flag.String("port", "8080", "Puerto en el que el servidor escuchará")
	imageDir := flag.String("imageDir", "./static/img/", "Directorio donde se encuentran las imágenes")
	theme := flag.String("theme", "", "Nombre del tema (subcarpeta de imágenes)")

	// Parsear las banderas
	flag.Parse()

	// Asignar la función handler a la ruta "/"
	http.HandleFunc("/", handler(*imageDir, *theme))

	// Servir archivos estáticos (imágenes, CSS, etc.)
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Mensaje indicando el puerto en uso
	fmt.Printf("Servidor ejecutándose en el puerto %s\n", *port)

	// Iniciar el servidor
	fmt.Printf("Servidor web en ejecución en http://localhost:%s\n", *port)
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
