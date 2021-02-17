package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Tiendas struct{
	Nombre string
	Descripcion string
	Contacto string
	Calificacion int
}

type Departamentos struct{
	Nombre string
	Tiendas []Tiendas
}

type Datos struct{
	Indice string
	Departamentos []Departamentos
}

type Archivo struct{
	Datos []Datos
}

func main() {
	request()
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "EDD_VirtualMall_201901816")
}

func getArreglo(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "[1,2,3,4]")
}

var Matriz [][][]Lista
var Arreglo []Lista

func CargarTienda(w http.ResponseWriter, r *http.Request){
	body, _ := ioutil.ReadAll(r.Body)
	var Data Archivo
	json.Unmarshal(body, &Data)
	fmt.Println(Data)
	fmt.Fprint(w, Data.Datos[1].Indice)

	//Necesitamos obtener el numero de categorias
	lenInd := len(Data.Datos)
	lenDep := len(Data.Datos[0].Departamentos)
	lenCal := 5

	//Creamos la matriz
	Matriz = make([][][]Lista, lenInd)
	for i := 0; i < lenInd; i++{
		Matriz[i] = make([][]Lista, lenDep)
		for j := 0; j < lenDep; j++{
			Matriz[i][j] = make([]Lista, lenCal)
			for k := 0; k < lenCal; k++{
				Matriz[i][j][k] = Lista{nil, nil, 0}
			}
		}
	}


	fmt.Println("Indice: ",len(Matriz))
	fmt.Println("Categoria: ",len(Matriz[0]))
	fmt.Println("Calificacion: ",len(Matriz[0][0]))
	fmt.Println("**********************************")

	//ingresamos las tiendas a la matriz
	for i := 0; i < len(Data.Datos); i++{
		for j := 0; j < len(Data.Datos[i].Departamentos); j++{
			for k := 0; k < len(Data.Datos[i].Departamentos[j].Tiendas); k++{
				fmt.Println("Posicion: ", i, ",", j, ",", Data.Datos[i].Departamentos[j].Tiendas[k].Calificacion - 1)
				Matriz[i][j][Data.Datos[i].Departamentos[j].Tiendas[k].Calificacion - 1].Add(Data.Datos[i].Departamentos[j].Tiendas[k])
				fmt.Println("Agregado: ", Data.Datos[i].Departamentos[j].Tiendas[k].Nombre)
				fmt.Println("**********************************")
			}
		}
	}
	//INGRESAMOS LA MATRIZ AL ARREGLO
	for i := 0; i < len(Matriz); i++{
		for j := 0; j < len(Matriz[0]); j++{
			for k := 0; k < len(Matriz[0][0]); k++{
				Arreglo[i+len(Matriz)*(j+len(Matriz[0]))] = Matriz[i][j][k]
			}
		}
	}
	for i:=0; i < len(Arreglo); i++{
		Arreglo[i].Imprimir()
	}
}

func request(){
	Servidor := mux.NewRouter().StrictSlash(true)
	Servidor.HandleFunc("/", homePage)
	Servidor.HandleFunc("/GetArreglo", getArreglo).Methods("GET")
	Servidor.HandleFunc("/cargartienda", CargarTienda).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", Servidor))
}

//CLASE NODO
type Nodo struct{
	Siguiente *Nodo
	Anterior *Nodo
	Dato Tiendas
}

//CLASE LISTA
type Lista struct{
	Inicio *Nodo
	Fin *Nodo
	len int
}

//CREANDO LA LISTA PAPA
func NewList() *Lista{
	return &Lista{nil, nil, 0}
}

//INSERTAR UN NODO
func (l *Lista) Add(valor Tiendas){
	nuevo := &Nodo{nil,nil,valor}
	if l.Inicio == nil{
		l.Inicio = nuevo
		l.Fin = nuevo
	}else{
		l.Fin.Siguiente = nuevo
		nuevo.Anterior = l.Fin
		l.Fin = nuevo
	}
	l.len++
}

//IMPRIMIENDO LA LISTA PAPA
func (l *Lista) Imprimir(){
	Aux := l.Inicio
	for Aux != nil{
		fmt.Print("<--[", Aux.Dato.Nombre, "]-->")
		Aux = Aux.Siguiente
	}
	fmt.Println()
	fmt.Println("Tamaño de la lista: ", l.len)
}
