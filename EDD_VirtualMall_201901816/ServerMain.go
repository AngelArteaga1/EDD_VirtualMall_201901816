package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

//VARIABLES GLOBALES
var Matriz [][][]Lista
var Arreglo []Lista

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

type BusquedaEspecifica struct{
	Departamento string
	Nombre string
	Calificacion int
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "EDD_VirtualMall_201901816")
}

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
	Arreglo = make([]Lista, lenDep*lenInd*lenCal)
	for i := 0; i < len(Matriz); i++{
		for j := 0; j < len(Matriz[0]); j++{
			for k := 0; k < len(Matriz[0][0]); k++{
				Arreglo[i+len(Matriz)*(j+len(Matriz[0])*k)] = Matriz[i][j][k]
			}
		}
	}
}

func TiendaEspecifica(w http.ResponseWriter, r *http.Request){
	if Arreglo == nil{
		fmt.Fprint(w, "Por favor, primero ingrese las tiendas")
	} else {
		body, _ := ioutil.ReadAll(r.Body)
		var Data BusquedaEspecifica
		json.Unmarshal(body, &Data)
		var existe int
		var tienda Tiendas
		Encontrado := false
		for i := 0; i < len(Arreglo); i++{
			if(Arreglo[i].Inicio != nil){
				existe = Arreglo[i].FindTienda(Data.Nombre, Data.Calificacion)
				if existe == 1{
					tienda = Arreglo[i].GetTienda(Data.Nombre, Data.Calificacion)
					Encontrado = true
				}
			}
		}
		if Encontrado == false{
			fmt.Fprint(w, "No se ha encontrado la tienda introducida :(")
		} else {
			TiendaEncontrada, err := json.Marshal(tienda)
			Errores(err)
			fmt.Fprint(w, string(TiendaEncontrada))
		}
	}
}

func BusquedaPosicion(w http.ResponseWriter, r *http.Request){
	if Arreglo == nil{
		fmt.Fprint(w, "Por favor, primero ingrese las tiendas")
	} else {
		vars := mux.Vars(r)
		var num string
		num = string(vars["num"])
		id, err := strconv.ParseInt(num,10,64)
		Errores(err)
		a := Arreglo[id].GetArray()
		json, err1 := json.Marshal(a)
		Errores(err1)
		fmt.Fprint(w, string(json))
	}
}

func Eliminar(w http.ResponseWriter, r *http.Request){
	if Arreglo == nil{
		fmt.Fprint(w, "Por favor, primero ingrese las tiendas")
	} else {
		body, _ := ioutil.ReadAll(r.Body)
		var Data BusquedaEspecifica
		json.Unmarshal(body, &Data)
		var existe int
		Encontrado := false
		for i := 0; i < len(Arreglo); i++{
			if(Arreglo[i].Inicio != nil){
				existe = Arreglo[i].FindTienda(Data.Nombre, Data.Calificacion)
				if existe == 1{
					Arreglo[i].DeleteTienda(Data.Nombre, Data.Calificacion)
					Encontrado = true
				}
			}
		}
		if Encontrado == false{
			fmt.Fprint(w, "No se ha encontrado la tienda introducida :(")
		} else {
			fmt.Fprint(w, "La tienda ha sido eliminada exitosamente!")
		}
	}
}

func GetArreglo(w http.ResponseWriter, r *http.Request){
	for i := 0; i< len(Arreglo); i++{
		fmt.Println("****************")
		Arreglo[i].Imprimir()
	}
	if Arreglo != nil{
		//ESCRIBIMOS LAS PRIMAS COSAS DL GRAFO
		f, err := os.Create("Grafica.dot")
		Errores(err)
		//w := bufio.NewWriter(f)
		f.WriteString("digraph structs {\n")
		f.WriteString("node [shape=record, fontname=\"Bookman Old Style\", " +
			"style=filled, fillcolor=lightpink];\nedge [dir=\"both\"]\nlabelloc=\"t\";" +
			"\nlabel=\"Estructura de Datos\";\nfontsize=30;\n")
		//CREAMOS LA ESTRUCTURA GENERAL
		struct1 := "struct [fillcolor=brown1, label=\""
		for i :=0; i < len(Arreglo); i++{
			if i == len(Arreglo)-1{
				struct1 = struct1 + "<f" + strconv.Itoa(i) + "> ___" + strconv.Itoa(i) + "___"
			} else {
				struct1 = struct1 + "<f" + strconv.Itoa(i) + "> ___" + strconv.Itoa(i) + "___|"
			}
		}
		struct1 = struct1 + "\"];\n"
		f.WriteString(struct1)
		//CREAMOS LOS NODOS
		for i := 0; i< len(Arreglo); i++ {
			if Arreglo[i].Inicio != nil{
				f.WriteString(Arreglo[i].GraphNodes(i))
			}
		}
		f.WriteString("}")
		//cmd := exec.Command("dot -Tpdf Grafica.dot -o Grafica.pdf")
		//cmd.Start()
		path, _ := exec.LookPath("dot")
		cmd, _ := exec.Command(path, "-Tsvg", "Grafica.dot").Output()
		mode := int(0777)
		ioutil.WriteFile("Grafica.svg", cmd, os.FileMode(mode))
		fmt.Fprint(w, "La gráfica fue realizada exitosamente!")
	} else {
		fmt.Fprint(w, "Por favor ingrese sus tiendas primero")
	}
}

func request(){
	Servidor := mux.NewRouter().StrictSlash(true)
	Servidor.HandleFunc("/", homePage)
	Servidor.HandleFunc("/cargartienda", CargarTienda).Methods("POST")
	Servidor.HandleFunc("/TiendaEspecifica", TiendaEspecifica).Methods("POST")
	Servidor.HandleFunc("/getArreglo", GetArreglo).Methods("GET")
	Servidor.HandleFunc("/Eliminar", Eliminar).Methods("DELETE")
	Servidor.HandleFunc("/id/{num:[0-9]+}", BusquedaPosicion).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", Servidor))
}


func Errores(e error) {
	if e != nil {
		panic(e)
	}
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

//BUSQUEDA POR NOMBRE Y CALIFICACION
func (l *Lista) GetTienda(nombre string, calificacion int) Tiendas{
	Aux := l.Inicio
	Encontrado := false
	var tienda Tiendas
	for Aux != nil || Encontrado != true{
		if nombre == Aux.Dato.Nombre && calificacion == Aux.Dato.Calificacion{
			Encontrado = true
			tienda = Aux.Dato
		}
		Aux = Aux.Siguiente
	}
	return tienda
}

//ELIMINACION POR NOMBRE Y CALIFICACION
func (l *Lista) DeleteTienda(nombre string, calificacion int){
	Aux := l.Inicio
	Encontrado := false
	if Aux.Dato.Nombre == nombre && Aux.Dato.Calificacion == calificacion{
		if Aux == l.Inicio && Aux == l.Fin{
			l.Inicio = nil
			l.Fin = nil
			l.len--
		} else {
			Aux.Siguiente.Anterior = nil
			l.Inicio = Aux.Siguiente
			l.len--
		}
	} else {
		for Aux != nil || Encontrado != true{
			if nombre == Aux.Dato.Nombre && calificacion == Aux.Dato.Calificacion{
				Encontrado = true
				if Aux == l.Fin{
					Aux.Anterior.Siguiente = nil
					l.len--
				} else {
					Aux.Anterior.Siguiente = Aux.Siguiente
					Aux.Siguiente.Anterior = Aux.Anterior
					l.len--
				}
			}
			Aux = Aux.Siguiente
		}
	}
}

//BUSQUEDA POR NOMBRE Y CALIFICACION
func (l *Lista) FindTienda(nombre string, calificacion int) int{
	Aux := l.Inicio
	resultado := -1
	for Aux != nil{
		if nombre == Aux.Dato.Nombre && calificacion == Aux.Dato.Calificacion{
			resultado = 1
		}
		Aux = Aux.Siguiente
	}
	return resultado
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

//OBTENER UN ARREGLO DE LA LISTA
func (l *Lista) GetArray() []Tiendas{
	a := make([]Tiendas, l.len)
	i := 0
	Aux := l.Inicio
	for Aux != nil{
		a[i] = Aux.Dato
		i++
		Aux = Aux.Siguiente
	}
	return a
}

//IMPRIMIENDO LOS NODOS EN GRAPHVIZ
func (l *Lista) GraphNodes(i int) string{
	Aux := l.Inicio
	nodos := ""
	j := 0
	for Aux != nil{
		nodos = nodos + "a" + strconv.Itoa(i) + "Node" + strconv.Itoa(j) + " [label=\""+ Aux.Dato.Nombre +"\"]\n"
		j++
		Aux = Aux.Siguiente
	}
	k := 0
	nodos = nodos + "struct:f" + strconv.Itoa(i)
	Aux = l.Inicio
	for Aux != nil{
		nodos = nodos + " -> a" + strconv.Itoa(i) + "Node" + strconv.Itoa(k)
		k++
		Aux = Aux.Siguiente
	}
	nodos = nodos + ";\n"
	return nodos
}

//ESTA VACÍA
func (l *Lista) IsEmpty() bool{
	if l.Inicio == nil{
		return true
	} else {
		return false
	}
}
