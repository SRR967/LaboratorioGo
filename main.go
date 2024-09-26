package main

import (
	"encoding/base64"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

// Función para obtener el nombre del host
func getHostName() string {
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatal("Error obteniendo el nombre del host:", err)
	}
	return hostName
}

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
		log.Fatal("Error leyendo el directorio:", err)
	}

	var imageFiles []ImageData
	for _, file := range files {
		if !file.IsDir() && isImage(file.Name()) {
			imagePath := filepath.Join(dirPath, file.Name())
			encodedImage := encodeImageToBase64(imagePath)
			imageFiles = append(imageFiles, ImageData{
				Base64: template.URL("data:image/jpeg;base64," + encodedImage), // Aquí usamos el formato adecuado para imágenes Base64
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
	randomIndex := rand.Intn(len(imageFiles))
	return imageFiles[randomIndex]
}

// Función para leer y codificar una imagen en Base64
func encodeImageToBase64(imagePath string) string {
	imageData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		log.Fatal("Error leyendo el archivo de imagen:", err)
	}
	return base64.StdEncoding.EncodeToString(imageData)
}

// Función para seleccionar una plantilla aleatoria
func randomTemplate() string {
	templates := []string{"plantilla1.html", "plantilla2.html"}
	return templates[rand.Intn(len(templates))]
}

// Función para renderizar la plantilla con imágenes y datos
func renderTemplate(tmpl string, data PageData, w http.ResponseWriter) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		log.Fatal("Error cargando la plantilla:", err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal("Error renderizando la plantilla:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	dirPath := flag.String("dir", "", "Directorio que contiene las imágenes")
	theme := flag.String("theme", "Tema predeterminado", "El tema para la página")
	flag.Parse()

	// Verificar si se proporcionó un directorio
	if *dirPath == "" {
		log.Fatal("Debe proporcionar un directorio que contenga las imágenes. Use el flag -dir.")
	}

	imageFiles := getImageFiles(*dirPath)
	if len(imageFiles) == 0 {
		http.Error(w, "No se encontraron imágenes en la carpeta.", http.StatusNotFound)
		return
	}

	hostName := getHostName()

	// Selecciona 3 imágenes para mostrar
	selectedImages := selectNImages(imageFiles, 3)

	// Estructura que contiene los datos para la plantilla
	data := PageData{
		HostName: hostName,
		Images:   selectedImages,
		Theme:    *theme,
	}
	selectedTemplate := randomTemplate()

	// Renderiza la plantilla
	renderTemplate(selectedTemplate, data, w)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}