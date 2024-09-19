package main

import "fmt"

func main() {
	numbers := []int{1,2}
	numbers = append(numbers, 3, 4, 5)
	fmt.Println(numbers)
	res := suma(numbers[1], numbers[3])
	fmt.Println(res)
	lista := impresion()
	fmt.Println(lista)
}

// Función que recibe dos parámetros y devuelve un entero
func suma(a int, b int) int {
    return a + b
}

func impresion() []int {
	var lista []int
	for i := 0; i < 10; i++ {
        lista = append(lista, i)
    }
	return lista
}