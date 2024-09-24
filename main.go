package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

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

// Función para leer las imágenes de un directorio
func getImageFiles(dirPath string) []string {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal("Error leyendo el directorio:", err)
	}

	var imageFiles []string
	for _, file := range files {
		if !file.IsDir() && isImage(file.Name()) {
			imageFiles = append(imageFiles, file.Name())
		}
	}
	return imageFiles
}

// Función para seleccionar una imagen al azar
func getRandomImage(imageFiles []string) string {
	randomIndex := rand.Intn(len(imageFiles))
	return imageFiles[randomIndex]
}

// Función para leer y codificar una imagen en Base64
func encodeImageToBase64(imagePath string) string {
	imageData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		log.Fatal("Error leyendo el archivo de imagen:", err)
	}

	encodedImage := base64.StdEncoding.EncodeToString(imageData)
	return encodedImage
}

func main() {
	// Argumento de la línea de comandos para la ruta del directorio
	dirPath := flag.String("dir", "", "Directorio que contiene las imágenes")
	flag.Parse()

	// Verificar si se proporcionó un directorio
	if *dirPath == "" {
		log.Fatal("Debe proporcionar un directorio que contenga las imágenes. Use el flag -dir.")
	}

	// Obtener y mostrar el nombre del host
	hostName := getHostName()
	fmt.Printf("Nombre del host: %s\n", hostName)

	// Obtener la lista de imágenes del directorio
	imageFiles := getImageFiles(*dirPath)
	if len(imageFiles) == 0 {
		fmt.Println("No se encontraron imágenes en la carpeta.")
		return
	}

	// Seleccionar una imagen al azar
	randomImage := getRandomImage(imageFiles)
	fmt.Printf("Imagen seleccionada al azar: %s\n", randomImage)

	// Codificar la imagen seleccionada en Base64
	imagePath := filepath.Join(*dirPath, randomImage)
	encodedImage := encodeImageToBase64(imagePath)

	// Imprimir la imagen codificada en Base64
	fmt.Println("Imagen codificada en Base64:")
	fmt.Println(encodedImage)

}
