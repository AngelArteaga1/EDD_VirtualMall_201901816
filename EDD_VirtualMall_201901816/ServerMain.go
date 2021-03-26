package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	_ "strings"
)

//VARIABLES GLOBALES
var Matriz [][][]Lista
var Arreglo []Lista
var MPosiciones [][][]Posicion
var APosiciones []Posicion
var lenInd int
var lenDep int
var lenCal int
var ListaInventario ListaProducto
var ListaProductos ListaPro
var ArbolPedidos AVL
var Carrito ListaItem

type XTiendas struct {
	Nombre string
	Descripcion string
	Contacto string
	Calificacion int
	Logo string
	Departamento string
	Indice string
}

type Tiendas struct{
	Nombre string
	Descripcion string
	Contacto string
	Calificacion int
	Logo string
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

type Posicion struct {
	Indice string
	Departamento string
}

//ESTO ES PARA LA CARGA DE PRODUCTOS
type ArchivoProductos struct {
	Inventarios []Inventario
}
type Inventario struct {
	Tienda string
	Departamento string
	Calificacion int
	Productos []Producto
}
type Producto struct {
	Nombre string
	Codigo int
	Descripcion string
	Precio float32
	Cantidad int
	Imagen string
}

type ArchivoPedidos struct {
	Pedidos []Pedido
}
type Pedido struct {
	Fecha string
	Tienda string
	Departamento string
	Calificacion int
	Productos []CodigoX
}
type CodigoX struct {
	Codigo int
}

type item struct {
	Tienda string
	Departamento string
	Calificacion int
	Producto string
	Descripcion string
	Imagen string
	Precio int
	Codigo int
}

type Cad struct {
	Cadena string
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
	//fmt.Println(body)

	//Necesitamos obtener el numero de categorias
	lenInd = len(Data.Datos)
	lenDep = len(Data.Datos[0].Departamentos)
	lenCal = 5

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

	//Creamos la matriz posiciones
	MPosiciones = make([][][]Posicion, lenInd)
	for i := 0; i < lenInd; i++{
		MPosiciones[i] = make([][]Posicion, lenDep)
		for j := 0; j < lenDep; j++{
			MPosiciones[i][j] = make([]Posicion, lenCal)
		}
	}

	//ingresamos las tiendas a la matriz
	for i := 0; i < len(Data.Datos); i++{
		for j := 0; j < len(Data.Datos[i].Departamentos); j++{
			for k := 0; k < len(Data.Datos[i].Departamentos[j].Tiendas); k++{
				//fmt.Println("Posicion: ", i, ",", j, ",", Data.Datos[i].Departamentos[j].Tiendas[k].Calificacion - 1)
				Matriz[i][j][Data.Datos[i].Departamentos[j].Tiendas[k].Calificacion - 1].Add(Data.Datos[i].Departamentos[j].Tiendas[k])
				//fmt.Println("Agregado: ", Data.Datos[i].Departamentos[j].Tiendas[k].Nombre)
				//fmt.Println("**********************************")
			}
		}
	}

	//DEJAMOS UNA REPLICA CON LAS COLUMNAS Y FILAS
	for i := 0; i <lenInd; i++{
		for j := 0; j < lenDep; j++{
			for k := 0; k < lenCal; k++{
				MPosiciones[i][j][k] = Posicion{Data.Datos[i].Indice,Data.Datos[i].Departamentos[j].Nombre}
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

	//INGRESAMOS LA MATRIZ AL ARREGLO DE POSICIONES
	APosiciones = make([]Posicion, lenDep*lenInd*lenCal)
	for i := 0; i < len(MPosiciones); i++{
		for j := 0; j < len(MPosiciones[0]); j++{
			for k := 0; k < len(MPosiciones[0][0]); k++{
				APosiciones[i+len(MPosiciones)*(j+len(MPosiciones[0])*k)] = MPosiciones[i][j][k]
			}
		}
	}

	//IMPRIMIR LAS POSICIONES
	for i := 0; i < len(MPosiciones); i++{
		for j := 0; j < len(MPosiciones[0]); j++{
			for k := 0; k < len(MPosiciones[0][0]); k++{
				fmt.Println("*********************************")
				fmt.Println("Posicion: ", i, ", ", j, ", ", k)
				fmt.Println("Indice: ",MPosiciones[i][j][k].Indice, "Departamento: ",MPosiciones[i][j][k].Departamento)
			}
		}
	}
	fmt.Fprint(w, "Los datos han sido cargados exitosamente!")
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
			if APosiciones[i].Departamento == Data.Departamento{
				if(Arreglo[i].Inicio != nil){
					existe = Arreglo[i].FindTienda(Data.Nombre, Data.Calificacion)
					if existe == 1{
						tienda = Arreglo[i].GetTienda(Data.Nombre, Data.Calificacion)
						Encontrado = true
					}
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
		if int(id) >= len(Arreglo) || int(id)< 0{
			fmt.Fprint(w, "Esa posición del arreglo no existe")
		} else {
			a := Arreglo[id].GetArray()
			json, err1 := json.Marshal(a)
			Errores(err1)
			fmt.Fprint(w, string(json))
		}
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

func GuardarDatos(w http.ResponseWriter, r *http.Request){
	if Arreglo == nil{
		fmt.Fprint(w, "Por favor, ingrese sus datos primero")
	} else {

		var tiendas [][]Lista
		var ArrayTienda []Tiendas

		//INICIALIZAMOS LAS TIENDAS
		tiendas = make([][]Lista, lenInd)
		for i := 0; i < lenInd; i++{
			tiendas[i] = make([]Lista, lenDep)
			for j := 0; j < lenDep; j++{
				tiendas[i][j] = Lista{nil, nil, 0}
			}
		}

		//OBTENEMOS TODAS LAS TIENDAS EN UN ARREGLO
		for i := 0; i < len(Arreglo); i++{
			if(Arreglo[i].Inicio != nil){
				ArrayTienda = Arreglo[i].GetArray()
				//fmt.Println(tiendas)
				for j := 0; j < len(ArrayTienda); j++{
					tiendas[GetPosicionIndice(i)][GetPosicionDepartamento(i)].Add(ArrayTienda[j])
				}
			}
		}

		//JUNTAMOS EL STRUCT
		var data Archivo
		var indice []Datos
		indice = make([]Datos, lenInd)
		var departamentos []Departamentos
		for i := 0; i < len(tiendas); i++{
			departamentos = make([]Departamentos, lenDep)
			for j := 0; j < len(tiendas[0]); j++{
				ArrayTienda = tiendas[i][j].GetArray()
				departamentos[j].Tiendas = ArrayTienda
				departamentos[j].Nombre = MPosiciones[i][j][0].Departamento
			}
			//indice[i].Departamentos = make(departamentos)
			indice[i].Departamentos = departamentos
			indice[i].Indice = MPosiciones[i][0][0].Indice
		}
		data.Datos = indice
		json, err := json.Marshal(data)
		Errores(err)
		fmt.Println(string(json))
		fmt.Fprint(w, string(json))
	}
}

func convertMonth(mes int) string{
	if mes == 1{
		return "ENERO"
	} else if mes == 2 {
		return "FEBRERO"
	} else if mes == 3 {
		return "MARZO"
	} else if mes == 4 {
		return "ABRIL"
	} else if mes == 5 {
		return "MAYO"
	} else if mes == 6 {
		return "JUNIO"
	} else if mes == 7 {
		return "JULIO"
	} else if mes == 8 {
		return "AGOSTO"
	} else if mes == 9 {
		return "SEPTIEMBRE"
	} else if mes == 10 {
		return "OCTUBRE"
	} else if mes == 11 {
		return "NOVIEMBRE"
	} else if mes == 12 {
		return "DICIEMBRE"
	}
	return "EQUIS DE"
}

func request(){
	Servidor := mux.NewRouter().StrictSlash(true)
	Servidor.HandleFunc("/", homePage)
	//FASE 1
	Servidor.HandleFunc("/cargartienda", CargarTienda).Methods("POST")
	Servidor.HandleFunc("/TiendaEspecifica", TiendaEspecifica).Methods("POST")
	Servidor.HandleFunc("/getArreglo", GetArreglo).Methods("GET")
	Servidor.HandleFunc("/Eliminar", Eliminar).Methods("DELETE")
	Servidor.HandleFunc("/id/{num:[0-9]+}", BusquedaPosicion).Methods("GET")
	Servidor.HandleFunc("/guardar", GuardarDatos).Methods("GET")
	//FASE 2
	Servidor.HandleFunc("/cargarinventario", CargarInventario).Methods("POST")
	Servidor.HandleFunc("/getTiendas", GetTiendas).Methods("GET")
	Servidor.HandleFunc("/getProductos",GetProductos).Methods("POST")
	Servidor.HandleFunc("/cargarPedidos",CargarPedidos).Methods("POST")
	Servidor.HandleFunc("/graficarPedidos", GraphPedidos).Methods("POST")
	Servidor.HandleFunc("/setItemCarrito", setItemCarrito).Methods("POST")
	Servidor.HandleFunc("/DeleteItemCarrito", DeleteItemCarrito).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(Servidor)))
}

func TrueItem(data item)item{

	for i := 0; i < ListaInventario.len; i++{
		bb := ListaInventario.Get(i)
		if bb.calificacion == data.Calificacion && bb.departamento == data.Departamento && bb.tienda == data.Tienda{
			arbol := bb.Dato
			if arbol.exist(arbol.root, data.Codigo) == true{
				data.Imagen = arbol.get(arbol.root, data.Codigo).value.Imagen
			}
		}
	}
	return data
}

func DeleteItemCarrito(w http.ResponseWriter, r *http.Request){
	fmt.Println("*****************DESDE ACA PAPA*************")
	body, _ := ioutil.ReadAll(r.Body)
	var Data item
	json.Unmarshal(body, &Data)
	fmt.Println(Data)
	//AQUI EMPEZAMOS
	if &Carrito == nil{
		if Data.Codigo != -1  && Data.Calificacion != -1 && Data.Precio != -1{
			Carrito = ListaItem{nil, nil, 0}
			Data = TrueItem(Data)
			Carrito.Delete(Data)
		}
	} else {
		if Data.Codigo != -1  && Data.Calificacion != -1 && Data.Precio != -1{
			Data = TrueItem(Data)
			Carrito.Delete(Data)
		}
	}
	CarritoArray := Carrito.GetArray()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CarritoArray)
}

func setItemCarrito(w http.ResponseWriter, r *http.Request){
	fmt.Println("*****************DESDE ACA PAPA*************")
	body, _ := ioutil.ReadAll(r.Body)
	var Data item
	json.Unmarshal(body, &Data)
	fmt.Println(Data)
	//AQUI EMPEZAMOS
	if &Carrito == nil{
		if Data.Codigo != -1  && Data.Calificacion != -1 && Data.Precio != -1{
			Carrito = ListaItem{nil, nil, 0}
			Data = TrueItem(Data)
			Carrito.Add(Data)
		}
	} else {
		if Data.Codigo != -1  && Data.Calificacion != -1 && Data.Precio != -1{
			Data = TrueItem(Data)
			Carrito.Add(Data)
		}
	}
	CarritoArray := Carrito.GetArray()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CarritoArray)
}

func GraphPedidos(w http.ResponseWriter, r *http.Request){
	body, _ := ioutil.ReadAll(r.Body)
	var Cadena Cad
	json.Unmarshal(body, &Cadena)
	fmt.Println("DATA: ", Cadena)
	fecha := strings.Split(Cadena.Cadena, "-")
	//obtenemos los datos necesarios:
	year, _ := strconv.ParseInt(fecha[2], 10, 64)
	month, _ := strconv.ParseInt(fecha[1], 10, 64)
	day, _ := strconv.ParseInt(fecha[0], 10, 64)
	dep := fecha[3]
	fmt.Println("YEAR: ", year)
	if ArbolPedidos.raiz != nil{
		// *************************GRAFICAMOS LOS AÑOS *************************
		f, err := os.Create("Dot/years.dot")
		Errores(err)
		f.WriteString("digraph years {\n")
		f.WriteString("rankdir=UD\n")
		f.WriteString("node[shape=box, fontname=\"Bookman Old Style\", style=filled, fillcolor=lightpink]\n")
		f.WriteString("concentrate=true\n")
		f.WriteString("labelloc=\"t\";\nlabel=\"Estructura de los Años\";\nfontsize=30;\n")
		f.WriteString(graph(ArbolPedidos.raiz))
		f.WriteString("}\n")

		path, _ := exec.LookPath("dot")
		cmd, _ := exec.Command(path, "-Tpng", "Dot/years.dot").Output()
		mode := int(0777)
		ioutil.WriteFile("C:\\Users\\Angel Arteaga\\Documents\\GitHub\\EDD_VirtualMall_201901816\\Web\\VirtualMall\\src\\assets\\reportes\\years.png", cmd, os.FileMode(mode))
		fmt.Fprint(w, "La gráfica fue realizada exitosamente!")

		//************************* GRAFICAMOS LOS MESES *************************
		if existe(ArbolPedidos.raiz, int(year)) == true{
			listaMeses := get(ArbolPedidos.raiz, int(year))
			f, err = os.Create("Dot/month.dot")
			Errores(err)
			f.WriteString("digraph month {\n")
			f.WriteString("node [shape=Msquare, fontname=\"Bookman Old Style\", width=1.5, style=filled, " +
				"fillcolor=cadetblue];\nedge [dir=\"both\"]\nlabelloc=\"t\";\nlabel=\"" +
				"Estructura de Meses de "+ strconv.Itoa(int(year)) +"\";\nfontsize=30;\nrankdir=\"LR\";\n")
			f.WriteString(listaMeses.meses.GraphNodes())
			f.WriteString("}\n")
			path, _ := exec.LookPath("dot")
			cmd, _ := exec.Command(path, "-Tpng", "Dot/month.dot").Output()
			mode := int(0777)
			ioutil.WriteFile("C:\\Users\\Angel Arteaga\\Documents\\GitHub\\EDD_VirtualMall_201901816\\Web\\VirtualMall\\src\\assets\\reportes\\month.png", cmd, os.FileMode(mode))
			fmt.Fprint(w, "La gráfica fue realizada exitosamente!")
			//************************* GRAFICAMOS LOS MESES *************************
			if listaMeses.meses.existe(convertMonth(int(month))) == true {
				matrizDay := listaMeses.meses.get(convertMonth(int(month)))
				f, err = os.Create("Dot/days.dot")
				Errores(err)
				f.WriteString("digraph days {\n")
				f.WriteString("rankdir = TB;\nnode [shape=rectangle, height=0.5, width=1.5, style = filled];\ngraph[ nodesep = 0.5];\n")
				f.WriteString("labelloc=\"t\";\nlabel=\"Estructura de "+ strconv.Itoa(int(year)) +
					", " + strings.ToLower(matrizDay.mes) + "\";\nfontsize=30;\n")
				f.WriteString(matrizDay.m.graphMatrix())
				f.WriteString("}")

				path, _ := exec.LookPath("dot")
				cmd, _ := exec.Command(path, "-Tpng", "Dot/days.dot").Output()
				mode := int(0777)
				ioutil.WriteFile("C:\\Users\\Angel Arteaga\\Documents\\GitHub\\EDD_VirtualMall_201901816\\Web\\VirtualMall\\src\\assets\\reportes\\days.png", cmd, os.FileMode(mode))
				fmt.Fprint(w, "La gráfica fue realizada exitosamente!")
				//************************* GRAFICAMOS EL DIA *************************
				if matrizDay.m.existe(int(day), GetDepartamento(dep)) == true{
					listaDay := matrizDay.m.get(int(day), GetDepartamento(dep))
					f, err = os.Create("Dot/exaxtday.dot")
					Errores(err)
					f.WriteString("digraph exaxtday {\n")
					f.WriteString("node [shape=record, style = filled, fillcolor=darkslategray];\nrankdir = LR\n")
					f.WriteString("labelloc=\"t\";\nlabel=\"Estructura, "+ strconv.Itoa(int(day)) +
						" de " + strings.ToLower(matrizDay.mes) + " del " + strconv.Itoa(int(year)) +
						", " + dep + "\";\nfontsize=30;\n")
					f.WriteString(listaDay.productos.GraphNodes())
					f.WriteString("}")
					path, _ := exec.LookPath("dot")
					cmd, _ := exec.Command(path, "-Tpng", "Dot/exaxtday.dot").Output()
					mode := int(0777)
					ioutil.WriteFile("C:\\Users\\Angel Arteaga\\Documents\\GitHub\\EDD_VirtualMall_201901816\\Web\\VirtualMall\\src\\assets\\reportes\\exaxtday.png", cmd, os.FileMode(mode))
					fmt.Fprint(w, "La gráfica fue realizada exitosamente!")
				} else {
				//NO EXISTE EL DIA CON SU DEPARTAMENTO
					f, err = os.Create("Dot/exaxtday.dot")
					f.WriteString("digraph month {\n")
					f.WriteString("node [shape=Msquare, fontname=\"Bookman Old Style\", style=filled, " +
						"fillcolor=lightpink];\nedge [dir=\"both\"]\nlabelloc=\"t\";\nlabel=\"" +
						"No Existen Pedidos Dentro de ese Dia y Departamento :c\";\nfontsize=30;\nrankdir=\"LR\";\n")
					f.WriteString("}\n")
					path, _ := exec.LookPath("dot")
					cmd, _ := exec.Command(path, "-Tpng", "Dot/exaxtday.dot").Output()
					mode := int(0777)
					ioutil.WriteFile("C:\\Users\\Angel Arteaga\\Documents\\GitHub\\EDD_VirtualMall_201901816\\Web\\VirtualMall\\src\\assets\\reportes\\exaxtday.png", cmd, os.FileMode(mode))
					fmt.Fprint(w, "La gráfica fue realizada exitosamente!")
				}
			} else {
				//NO EXISTE EL MES
				f, err = os.Create("Dot/days.dot")
				f.WriteString("digraph month {\n")
				f.WriteString("node [shape=Msquare, fontname=\"Bookman Old Style\", style=filled, " +
					"fillcolor=lightpink];\nedge [dir=\"both\"]\nlabelloc=\"t\";\nlabel=\"" +
					"No Existen Pedidos Dentro de ese Mes :c\";\nfontsize=30;\nrankdir=\"LR\";\n")
				f.WriteString("}\n")
				path, _ := exec.LookPath("dot")
				cmd, _ := exec.Command(path, "-Tpng", "Dot/days.dot").Output()
				mode := int(0777)
				ioutil.WriteFile("C:\\Users\\Angel Arteaga\\Documents\\GitHub\\EDD_VirtualMall_201901816\\Web\\VirtualMall\\src\\assets\\reportes\\days.png", cmd, os.FileMode(mode))
				fmt.Fprint(w, "La gráfica fue realizada exitosamente!")
				//NO EXISTE EL DIA CON SU DEPARTAMENTO
				f, err = os.Create("Dot/exaxtday.dot")
				f.WriteString("digraph month {\n")
				f.WriteString("node [shape=Msquare, fontname=\"Bookman Old Style\", style=filled, " +
					"fillcolor=lightpink];\nedge [dir=\"both\"]\nlabelloc=\"t\";\nlabel=\"" +
					"No Existen Pedidos Dentro de ese Dia y Departamento :c\";\nfontsize=30;\nrankdir=\"LR\";\n")
				f.WriteString("}\n")
				path, _ = exec.LookPath("dot")
				cmd, _ = exec.Command(path, "-Tpng", "Dot/exaxtday.dot").Output()
				mode = int(0777)
				ioutil.WriteFile("C:\\Users\\Angel Arteaga\\Documents\\GitHub\\EDD_VirtualMall_201901816\\Web\\VirtualMall\\src\\assets\\reportes\\exaxtday.png", cmd, os.FileMode(mode))
				fmt.Fprint(w, "La gráfica fue realizada exitosamente!")
			}
		} else {
			//NO EXISTE EL YEAR
			f, err = os.Create("Dot/month.dot")
			f.WriteString("digraph month {\n")
			f.WriteString("node [shape=Msquare, fontname=\"Bookman Old Style\", style=filled, " +
				"fillcolor=lightpink];\nedge [dir=\"both\"]\nlabelloc=\"t\";\nlabel=\"" +
				"No Existen Pedidos Dentro de ese Año :c\";\nfontsize=30;\nrankdir=\"LR\";\n")
			f.WriteString("}\n")
			path, _ := exec.LookPath("dot")
			cmd, _ := exec.Command(path, "-Tpng", "Dot/month.dot").Output()
			mode := int(0777)
			ioutil.WriteFile("C:\\Users\\Angel Arteaga\\Documents\\GitHub\\EDD_VirtualMall_201901816\\Web\\VirtualMall\\src\\assets\\reportes\\month.png", cmd, os.FileMode(mode))
			fmt.Fprint(w, "La gráfica fue realizada exitosamente!")
			//NO EXISTE EL MES
			f, err = os.Create("Dot/days.dot")
			f.WriteString("digraph month {\n")
			f.WriteString("node [shape=Msquare, fontname=\"Bookman Old Style\", style=filled, " +
				"fillcolor=lightpink];\nedge [dir=\"both\"]\nlabelloc=\"t\";\nlabel=\"" +
				"No Existen Pedidos Dentro de ese Mes :c\";\nfontsize=30;\nrankdir=\"LR\";\n")
			f.WriteString("}\n")
			path, _ = exec.LookPath("dot")
			cmd, _ = exec.Command(path, "-Tpng", "Dot/days.dot").Output()
			mode = int(0777)
			ioutil.WriteFile("C:\\Users\\Angel Arteaga\\Documents\\GitHub\\EDD_VirtualMall_201901816\\Web\\VirtualMall\\src\\assets\\reportes\\days.png", cmd, os.FileMode(mode))
			fmt.Fprint(w, "La gráfica fue realizada exitosamente!")
			//NO EXISTE EL DIA CON SU DEPARTAMENTO
			f, err = os.Create("Dot/exaxtday.dot")
			f.WriteString("digraph month {\n")
			f.WriteString("node [shape=Msquare, fontname=\"Bookman Old Style\", style=filled, " +
				"fillcolor=lightpink];\nedge [dir=\"both\"]\nlabelloc=\"t\";\nlabel=\"" +
				"No Existen Pedidos Dentro de ese Dia y Departamento :c\";\nfontsize=30;\nrankdir=\"LR\";\n")
			f.WriteString("}\n")
			path, _ = exec.LookPath("dot")
			cmd, _ = exec.Command(path, "-Tpng", "Dot/exaxtday.dot").Output()
			mode = int(0777)
			ioutil.WriteFile("C:\\Users\\Angel Arteaga\\Documents\\GitHub\\EDD_VirtualMall_201901816\\Web\\VirtualMall\\src\\assets\\reportes\\exaxtday.png", cmd, os.FileMode(mode))
			fmt.Fprint(w, "La gráfica fue realizada exitosamente!")
		}
	}
}

func CargarPedidos(w http.ResponseWriter, r *http.Request){
	body, _ := ioutil.ReadAll(r.Body)
	var Data ArchivoPedidos
	json.Unmarshal(body, &Data)
	pedidos := Data.Pedidos
	fmt.Println("********************* CARRITO ***********************")
	fmt.Println(Data)
	fmt.Println("HASTA AQUI: ", "0")
	//Reseteamos el carrito
	Carrito.Inicio = nil
	Carrito.len = 0
	Carrito.Fin = nil
	//RECORREMOS LOS PEDIDOS
	for i :=0; i < len(pedidos); i++{
		fmt.Println("************************* NUEVA ITERACION *************************")
		//obtenemos el array de productos
		productos := pedidos[i].Productos
		//obtenemos la fecha
		fecha := strings.Split(pedidos[i].Fecha, "-")
		//obtenemos los datos necesarios:
		year, _ := strconv.ParseInt(fecha[2], 10, 64)
		month, _ := strconv.ParseInt(fecha[1], 10, 64)
		day, _ := strconv.ParseInt(fecha[0], 10, 64)
		//datos de la tienda:
		tienda := pedidos[i].Tienda
		dep := pedidos[i].Departamento
		calificacion := pedidos[i].Calificacion
		fmt.Println("HASTA AQUI: ", "1")
		//AQUI EMPEZAMOS EL ALGORITMO DE INSERCION
		//Revisamos si existe el año dentro del arbol:
		if existe(ArbolPedidos.raiz, int(year)) == false{
			//INSERTAMOS COMO SI FUERA PRIMERA VEZ
			//Creamos la lista de codigos:
			ListaCodigos := &ListaInt{nil, nil, 0}
			ListaCodigos.Add(productos, tienda, dep, calificacion)
			//Creamos la matriz de listas de codigos:
			matrix := NewMatriz()
			fmt.Println("INDICE DEL DEPARTAMENTO ", dep , ": ", GetDepartamento(dep))
			matrix.Insert(ListaCodigos, int(day), GetDepartamento(dep))
			//Creamos la lista de matrices:
			listaMatrices := &ListaM{nil, nil, 0}
			listaMatrices.Add(matrix, convertMonth(int(month)))
			//Insertamos el nodo nuevo al arbol
			ArbolPedidos.Insertar(int(year), listaMatrices)
			fmt.Println("HASTA AQUI: ", "2")
		} else {
			//INSERTAMOS PARA EL NODO  DEL ARBOL EXISTENTE
			//Obtenemos el nodo donde se encuentre el anio:
			fmt.Println("HASTA AQUI: ", "3")
			nodoArbol := get(ArbolPedidos.raiz, int(year))
			if nodoArbol.meses.existe(convertMonth(int(month))) == false{
				//Creamos la lista de codigos:
				ListaCodigos := &ListaInt{nil, nil, 0}
				ListaCodigos.Add(productos, tienda, dep, calificacion)
				//Creamos la matriz de listas de codigos:
				matrix := NewMatriz()
				fmt.Println("INDICE DEL DEPARTAMENTO ", dep , ": ", GetDepartamento(dep))
				matrix.Insert(ListaCodigos, int(day), GetDepartamento(dep))
				//Creamos la lista de matrices:
				nodoArbol.meses.Add(matrix, convertMonth(int(month)))
				fmt.Println("HASTA AQUI: ", "4")
			} else {
				fmt.Println("HASTA AQUI: ", "5")
				//INSERTAMOS PARA EL NODO DE LA LISTA EXISTENTE
				//Obtenemos el nodo donde se encuentre el mes:
				nodoMes := nodoArbol.meses.get(convertMonth(int(month)))
				fmt.Println("HASTA AQUI: ", "5.1")
				fmt.Println("INDICE DEL DEPARTAMENTO ", dep , ": ", GetDepartamento(dep))
				if nodoMes.m.existe(int(day), GetDepartamento(dep)) == false {
					fmt.Println("HASTA AQUI: ", "5.5")
					//Creamos la lista de codigos:
					ListaCodigos := &ListaInt{nil, nil, 0}
					ListaCodigos.Add(productos, tienda, dep, calificacion)
					//Creamos la matriz de listas de codigos:
					fmt.Println("INDICE DEL DEPARTAMENTO ", dep , ": ", GetDepartamento(dep))
					nodoMes.m.Insert(ListaCodigos, int(day), GetDepartamento(dep))
					fmt.Println("HASTA AQUI: ", "6")
				} else {
					fmt.Println("HASTA AQUI: ", "7")
					//INSERTAMOS PARA EL NODO DE LA POSICION DE LA MATRIZ EXISTENTE
					//Obtenemos el nodo donde se encuentre el dia:
					fmt.Println("INDICE DEL DEPARTAMENTO ", dep , ": ", GetDepartamento(dep))
					nodoDia := nodoMes.m.get(int(day), GetDepartamento(dep))
					//Agregamos el pedido
					nodoDia.productos.Add(productos, tienda, dep, calificacion)
					fmt.Println("HASTA AQUI: ", "8")
				}
			}
		}
	}
	fmt.Fprint(w, "La carga fue realizada exitosamente!")
}

func GetProductos(w http.ResponseWriter, r *http.Request){
	body, _ := ioutil.ReadAll(r.Body)
	var Data BusquedaEspecifica
	var Array []Producto
	json.Unmarshal(body, &Data)
	fmt.Println("DATOS: ")
	fmt.Println(Data)
	if &ListaInventario != nil{
		for i := 0; i < ListaInventario.len; i++{
			temp := ListaInventario.Get(i)
			fmt.Println(temp.tienda, Data.Nombre)
			fmt.Println(temp.departamento, Data.Departamento)
			fmt.Println(temp.calificacion, Data.Calificacion)
			if (temp.tienda == Data.Nombre && temp.departamento == Data.Departamento && temp.calificacion == Data.Calificacion){
				fmt.Println("Encontrado")
				ListaProductos = ListaPro{nil, nil, 0}
				arbol := temp.Dato
				arbol.createList(arbol.root)
				Array = ListaProductos.GetArray()
				break
			}
		}
	}
	fmt.Println(Array)
	//fmt.Println(json.NewEncoder(w).Encode(Array))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Println("Aqui todo bien")
	json.NewEncoder(w).Encode(Array)
}

func GetTiendas(w http.ResponseWriter, r *http.Request){
	//ESTA FUNCION SOLO OBTIENE UNA ARREGLO DE TODAS LAS TIENDAS CARGADAS
	var ArregloTiendas []XTiendas
	ListaTiendas := &ListaXTienda{nil, nil, 0}
	if Arreglo != nil {
		for i := 0; i < len(Arreglo); i++{
			if Arreglo[i].len != 0{
				tienditas := Arreglo[i].GetArray()
				for j := 0; j < len(tienditas); j++{
					tiendita := tienditas[j]
					//CONSEGUIMOS LA CATEGORIA Y TODA LA ONDA
					departamento := APosiciones[i].Departamento
					indice := APosiciones[i].Indice
					xTienda := XTiendas{tiendita.Nombre, tiendita.Descripcion, tiendita.Contacto, tiendita.Calificacion, tiendita.Logo, departamento, indice}
					ListaTiendas.Add(xTienda)
				}
			}
		}
		ArregloTiendas = ListaTiendas.GetArray()
		//salida, _ := json.Marshal(ArregloTiendas)
		fmt.Println(ArregloTiendas)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ArregloTiendas)
	}
}

func CargarInventario(w http.ResponseWriter, r *http.Request){
	body, _ := ioutil.ReadAll(r.Body)
	var Data ArchivoProductos
	json.Unmarshal(body, &Data)
	fmt.Println(Data)
	if &ListaInventario == nil{
		ListaInventario = ListaProducto{nil, nil, 0}
		//METEMOS LOS PRODUCTOS EN EL ARREGLO
		for i := 0; i < len(Data.Inventarios); i++{
			tienda := Data.Inventarios[i].Tienda
			departamento := Data.Inventarios[i].Departamento
			calificacion := Data.Inventarios[i].Calificacion
			arbol := BST{}
			for j := 0; j < len(Data.Inventarios[i].Productos); j++{
				producto := Data.Inventarios[i].Productos[j]
				arbol.add(producto)
			}
			ListaInventario.Add(arbol, tienda, departamento, calificacion)
		}
	} else {
		for i := 0; i < len(Data.Inventarios); i++{
			tienda := Data.Inventarios[i].Tienda
			departamento := Data.Inventarios[i].Departamento
			calificacion := Data.Inventarios[i].Calificacion
			existe := false
			//AHORA BUSCAMOS SI EXISTE ESE ARBOL
			for i := 0; i < ListaInventario.len; i++{
				temp := ListaInventario.Get(i)
				if temp.tienda == tienda && temp.departamento == departamento && temp.calificacion == calificacion{
					for j := 0; j < len(Data.Inventarios[i].Productos); j++{
						producto := Data.Inventarios[i].Productos[j]
						temp.Dato.add(producto)
					}
					existe = true
				}
			}
			//SI NO EXISTE EL ARBOL ENTONCES
			if existe == false{
				arbol := BST{}
				for j := 0; j < len(Data.Inventarios[i].Productos); j++{
					producto := Data.Inventarios[i].Productos[j]
					arbol.add(producto)
				}
				ListaInventario.Add(arbol, tienda, departamento, calificacion)
			}
		}
	}
	for i := 0; i < ListaInventario.len; i++{
		temp := ListaInventario.Get(i)
		fmt.Println("ESTE ES EL ARBOL DE:", temp.tienda)
		temp.Dato.imprimir(temp.Dato.root)
		fmt.Println()
	}

}

func GetIndice(indice string)int{
	for i:=0; i< lenInd; i++{
		for j:=0; j< lenDep; j++{
			for k:=0; k< lenCal; k++{
				if indice == MPosiciones[i][j][k].Indice{
					return i
				}
			}
		}
	}
	return -1
}

func GetDepartamento(departamento string)int{
	for i:=0; i< lenInd; i++{
		for j:=0; j< lenDep; j++{
			for k:=0; k< lenCal; k++{
				if departamento == MPosiciones[i][j][k].Departamento{
					return j
				}
			}
		}
	}
	return -1
}

func GetPosicionIndice(posicion int)int{
	indice := APosiciones[posicion].Indice
	return GetIndice(indice)
}

func GetPosicionDepartamento(posicion int)int{
	departamento := APosiciones[posicion].Departamento
	return GetDepartamento(departamento)
}

func Errores(e error) {
	if e != nil {
		panic(e)
	}
}

/*****************************************************************************/
//CLASE NODO ITEM
type NodoItem struct {
	Siguiente *NodoItem
	Anterior *NodoItem
	Dato item
}
//CLASE LISTA ITEM
type ListaItem struct{
	Inicio *NodoItem
	Fin *NodoItem
	len int
}
//INSERTAR
func (l *ListaItem) Add(valor item){
	nuevo := &NodoItem{nil,nil,valor}
	if l.Inicio == nil{
		l.Inicio = nuevo
		l.Fin = nuevo
	}else{
		l.Fin.Siguiente = nuevo
		nuevo.Anterior = l.Fin
		l.Fin = nuevo
	}
	l.len++
	fmt.Println("add:")
	fmt.Println("")
	l.print()
}
//ELIMINAR
func (l *ListaItem) Delete(dato item){
	Aux := l.Inicio
	Encontrado := false
	if Aux.Dato.Codigo == dato.Codigo{
		if l.Inicio == l.Fin{
			l.Inicio = nil
			l.Fin = nil
			l.len--
		} else {
			l.Inicio.Siguiente.Anterior = nil
			l.Inicio = Aux.Siguiente
			l.len--
		}
	} else {
		for Aux != nil || Encontrado != true{
			if dato.Codigo == Aux.Dato.Codigo{
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
	fmt.Println("delete:", l.len)
	fmt.Println("")
	l.print()
}
//ELIMINAR2
func (l*ListaItem) Delete2(dato item){
	tmp := l.Inicio
	for tmp != nil {
		if tmp.Dato.Codigo == dato.Codigo{
			if tmp.Anterior == nil{
				l.Inicio = tmp.Siguiente
				l.len--
			} else if tmp.Siguiente == nil{
				l.Fin = tmp.Anterior
				l.len--
			} else {
				tmp.Anterior.Siguiente = tmp.Siguiente
				tmp.Siguiente.Anterior = tmp.Anterior
				l.len--
			}
			break
		}
		tmp = tmp.Siguiente
	}
}
//GET ARRAY
func (l *ListaItem) GetArray() []item{
	a := make([]item, l.len)
	i := 0
	Aux := l.Inicio
	for Aux != nil{
		a[i] = Aux.Dato
		i++
		Aux = Aux.Siguiente
	}
	return a
}
//PRINT
func (l *ListaItem) print(){
	tmp := l.Inicio
	for tmp != nil{
		fmt.Print(tmp.Dato.Codigo, "<->")
		tmp = tmp.Siguiente
	}
	fmt.Println("")
}

/*****************************************************************************/
//CLASE NODO INT
type NodoInt struct {
	Siguiente *NodoInt
	Anterior *NodoInt
	Dato []CodigoX
	tienda string
	departamento string
	calificacion int
}
//CLASE LISTA INT
type ListaInt struct{
	Inicio *NodoInt
	Fin *NodoInt
	len int
}
//ADD
//INSERTAR UN NODO INT
func (l *ListaInt) Add(valor []CodigoX, tienda, departamento string, calificacion int){
	nuevo := &NodoInt{nil,nil,valor, tienda, departamento, calificacion}
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
//GRAFICAR NODOS
func (l *ListaInt) GraphNodes()string{
	cadena := ""
	cont := 0
	temp := l.Inicio
	for temp != nil{
		cadena = cadena + "struct" + strconv.Itoa(cont)  + "[label=\""
		//fmt.Println("HOLA BB")
		fmt.Println(temp.Dato)
		for i := 0; i < len(temp.Dato); i++{
			//OBTENEMOS LA INFORMACION DEL PRODUCTO
			//fmt.Println("HOLA BB2")
			name := temp.tienda
			dep := temp.departamento
			cal := temp.calificacion
			cod := temp.Dato[i].Codigo
			fmt.Println("Codigo actual: ",cod)
			//Buscamos el codigo dentro de los productos:
			for j := 0; j < ListaInventario.len; j++{
				//fmt.Println("HOLA BB3")
				NodoAct := ListaInventario.Get(j)
				if NodoAct.tienda == name && NodoAct.departamento == dep && NodoAct.calificacion == cal{
					ArbolAct := NodoAct.Dato
					//fmt.Println("HOLA BB4")
					if ArbolAct.exist(ArbolAct.root, cod) ==  true {
						//fmt.Println("HOLA BB5")
						product := ArbolAct.get(ArbolAct.root, cod)
						cadena = cadena + "Producto: " + product.value.Nombre + "\\n"
						cadena = cadena + "Precio: " + strconv.Itoa(int(product.value.Precio)) + "\\n"
						cadena = cadena + "Codigo: " + strconv.Itoa(product.value.Codigo) + "\\n"
						if i != len(temp.Dato)-1{
							cadena = cadena + "|"
						}
					}
				}
			}
		}
		cadena = cadena + "\", fontcolor=\"aliceblue\"];\n"
		cont++
		temp = temp.Siguiente
	}
	cont = 0
	temp = l.Inicio
	for temp != nil{
		if cont == l.len-1{
			cadena = cadena + "struct" + strconv.Itoa(cont)
		} else {
			cadena = cadena + "struct" + strconv.Itoa(cont) + " ->"
		}
		cont++
		temp = temp.Siguiente
	}
	return cadena
}
/*****************************************************************************/
//CLASE NODO XTIENDA
type NodoXTienda struct {
	Siguiente *NodoXTienda
	Anterior *NodoXTienda
	Dato XTiendas
}
//CLASE LISTA
type ListaXTienda struct{
	Inicio *NodoXTienda
	Fin *NodoXTienda
	len int
}
//ADD
//INSERTAR UN NODO
func (l *ListaXTienda) Add(valor XTiendas){
	nuevo := &NodoXTienda{nil,nil,valor}
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
//OBTENER UN ARREGLO DE LA LISTA
func (l *ListaXTienda) GetArray() []XTiendas{
	a := make([]XTiendas, l.len)
	i := 0
	Aux := l.Inicio
	for Aux != nil{
		a[i] = Aux.Dato
		i++
		Aux = Aux.Siguiente
	}
	return a
}

/*****************************************************************************/
//CLASE LISTA MATRIZ
type ListaM struct{
	Inicio *NodoM
	Fin *NodoM
	len int
}
//CLASE NODO MATRIZ
type NodoM struct {
	Siguiente *NodoM
	Anterior *NodoM
	m *matriz
	mes string
}
//ADD
//INSERTAR UN NODO MATRIZ
func (l *ListaM) Add(valor *matriz, mes string){
	nuevo := &NodoM{nil,nil,valor, mes}
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
//EXISTE EL NODO
func (l *ListaM) existe(mes string) bool{
	Aux := l.Inicio
	for Aux != nil{
		if mes == Aux.mes{
			return true
		}
		Aux = Aux.Siguiente
	}
	return false
}
//RETURNAR NODO
func (l *ListaM) get(mes string) *NodoM{
	Aux := l.Inicio
	for Aux != nil{
		if mes == Aux.mes{
			return Aux
		}
		Aux = Aux.Siguiente
	}
	return nil
}
func (l *ListaM) GraphNodes() string{
	cadena := ""
	temp := l.Inicio
	for temp != nil{
		cadena = cadena + "nodo" + temp.mes + "[ label =\"" + temp.mes + "\"]\n"
		if temp.Anterior != nil{
			cadena = cadena + "nodo" + temp.Anterior.mes + "->nodo" + temp.mes + "\n"
		}
		temp = temp.Siguiente
	}
	return cadena
}

/*****************************************************************************/
//CLASE NODO PRODUCTO
type NodoProducto struct{
	Siguiente *NodoProducto
	Anterior *NodoProducto
	Dato BST
	tienda string
	departamento string
	calificacion int
}
//CLASE LISTA
type ListaProducto struct{
	Inicio *NodoProducto
	Fin *NodoProducto
	len int
}
//INSERTAR UN NODO
func (l *ListaProducto) Add(valor BST, tienda, departamento string, calificacion int){
	nuevo := &NodoProducto{nil,nil,valor, tienda, departamento, calificacion}
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
//FUNCION DE GET
func (l *ListaProducto) Get(indice int) NodoProducto{
	temp := l.Inicio
	i := 0
	for temp != nil{
		if indice == i{
			return *temp
		}
		i++
		temp = temp.Siguiente
	}
	return *l.Inicio
}

/*****************************************************************************/
//CLASE NODO PRODUCTO
type NodoPro struct{
	Siguiente *NodoPro
	Anterior *NodoPro
	Dato Producto
}
//CLASE LISTA
type ListaPro struct{
	Inicio *NodoPro
	Fin *NodoPro
	len int
}
//INSERTAR UN NODO
func (l *ListaPro) Add(valor Producto){
	nuevo := &NodoPro{nil,nil,valor}
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
//FUNCION DE GET
func (l *ListaPro) Get(indice int) NodoPro{
	temp := l.Inicio
	i := 0
	for temp != nil{
		if indice == i{
			return *temp
		}
		i++
		temp = temp.Siguiente
	}
	return *l.Inicio
}
//OBTENER UN ARREGLO DE LA LISTA
func (l *ListaPro) GetArray() []Producto{
	a := make([]Producto, l.len)
	i := 0
	Aux := l.Inicio
	for Aux != nil{
		a[i] = Aux.Dato
		i++
		Aux = Aux.Siguiente
	}
	return a
}

/*****************************************************************************/
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


/*****************************************************************************/
//ESTA ES LA CLASE DEL NODO
// Node is a value with two pointers
type Node struct {
	value Producto
	left  *Node
	right *Node
}
// BST is a set of sorted Nodes
type BST struct {
	root *Node
}
//FUNCION PARA AGREGAR
func (bst *BST) add(value Producto) {
	bst.root = bst._add(value, bst.root)
}
//FUNCION PARA AGREGAR X2
func (bst BST) _add(value Producto, tmp *Node) *Node {
	if tmp == nil {
		return &Node{value: value}
	} else if value.Codigo == tmp.value.Codigo{
		tmp.value.Cantidad = tmp.value.Cantidad + value.Cantidad
	} else if value.Codigo > tmp.value.Codigo {
		tmp.right = bst._add(value, tmp.right)
	} else {
		tmp.left = bst._add(value, tmp.left)
	}
	return tmp
}
//FUNCION PARA IMPRIMIR EL ARBOL EQUIS DE
func (bst BST) imprimir(tmp *Node) {
	if tmp != nil {
		fmt.Println("********************")
		fmt.Println("NOMBRE: ",tmp.value.Nombre)
		fmt.Println("DESCRIPCION: ",tmp.value.Descripcion)
		fmt.Println("CANTIDAD: ",tmp.value.Cantidad)
		fmt.Println("CODIGO: ",tmp.value.Codigo)
		fmt.Println("PRECIO: ",tmp.value.Precio)
		bst.imprimir(tmp.left)
		bst.imprimir(tmp.right)
	}
}
//FUNCION PA METER
func (bst BST) createList(tmp *Node){
	if tmp != nil{
		ListaProductos.Add(tmp.value)
		bst.createList(tmp.left)
		bst.createList(tmp.right)
	}
}
//FUNCION PARA REVISAR SI EXISTE
func (bst BST) exist(tmp *Node, value int) bool{
	if tmp != nil {
		if tmp.value.Codigo == value{
			return true
		} else {
			if (bst.exist(tmp.left, value) || bst.exist(tmp.right, value)){
				return true
			}
		}
	}
	return false
}
func (bst BST)get(tmp *Node, value int)*Node{
	if tmp != nil{
		if tmp.value.Codigo == value{
			return tmp
		} else {
			if (bst.get(tmp.left, value) != nil){
				return bst.get(tmp.left, value)
			} else {
				return bst.get(tmp.right, value)
			}
		}
	}
	return nil
}

/*****************************************************************************/
//CLASE PARA LA MATRIZ
type nodo struct {
	//Estos atributos son especificos para la matrix
	x,y int //Saber en que cabecera estoy
	productos *ListaInt //tipo de objeto
	izquierda, derecha, arriba, abajo *nodo //nodos con los que nos desplazamos dentro de la matriz
	//Estos atributos son especificos para la lista
	header int//tipo interno de la cabecera
	siguiente, anterior *nodo // nodos con los que nos vamos a desplazar dentro de las listas
}
type lista struct {
	first, last *nodo
}
type matriz struct {
	lst_h, lst_v *lista
}
func nodoMatriz(x int, y int, producto *ListaInt) *nodo {
	return &nodo{x,y,producto, nil,nil,nil,nil,0,nil,nil}
}
func nodoLista(header int) *nodo{
	return &nodo{0,0,nil,nil,nil,nil,nil,header,nil, nil}
}
func newLista() *lista{
	return &lista{nil,nil}
}
//Se cambio a primer letra mayuscula para poder acceder
func NewMatriz() *matriz{
	return &matriz{newLista(),newLista()}
}
func (n *nodo) headerX() int { return n.x }
func (n *nodo) headerY() int { return n.y }
func (l *lista ) ordenar(nuevo *nodo)  {
	aux := l.first
	for(aux != nil){
		if nuevo.header > aux.header{
			aux = aux.siguiente
		}else{
			if aux == l.first{
				nuevo.siguiente = aux
				aux.anterior = nuevo
				l.first = nuevo
			}else{
				nuevo.anterior = aux.anterior
				aux.anterior.siguiente = nuevo
				nuevo.siguiente = aux
				aux.anterior = nuevo
			}
			return
		}
	}
	l.last.siguiente = nuevo
	nuevo.anterior = l.last
	l.last = nuevo
}
func (l *lista) insert(header int) {
	nuevo := nodoLista(header)
	if l.first == nil{
		l.first = nuevo
		l.last = nuevo
	}else{
		l.ordenar(nuevo)
	}
}
func (l *lista) search(header int) *nodo{
	temp := l.first
	for temp != nil{
		if temp.header == header{
			return temp
		}
		temp = temp.siguiente
	}
	return nil
}
func (l *lista) print() {
	temp := l.first
	for temp != nil{
		fmt.Println("Cabecera:", temp.header)
		temp = temp.siguiente
	}
}
func (m *matriz) Insert(producto *ListaInt, x int, y int){
	h := m.lst_h.search(x)
	v := m.lst_v.search(y)

	if h==nil && v==nil{
		m.noExisten(producto, x,y)
	}else if h==nil && v!=nil{
		m.existeVertical(producto, x, y)
	}else if h!=nil && v==nil{
		m.existeHorizontal(producto, x, y)
	}else{
		m.existen(producto,x,y)
	}
}
func (m *matriz)noExisten(producto *ListaInt, x int, y int) {
	m.lst_h.insert(x)//insertamos en la lista que emula la cabecera horizontal
	m.lst_v.insert(y)//insertamos en la lista que emula la cabecera vertical

	h := m.lst_h.search(x)//vamos a buscar el nodo que acabos de insertar para poder enlazarlo
	v := m.lst_v.search(y)//vamos a buscar el nodo que acabos de insertar para poder enlazarlo

	nuevo := nodoMatriz(x,y,producto)//creamos nuevo nodo tipo matriz

	h.abajo = nuevo//enlazamos el nodo horizontal hacia abajo
	nuevo.arriba = h//enlazmos el nuevo nodo hacia arriba

	v.derecha = nuevo//enlazamos el nodo vertical hacia la derecha
	nuevo.izquierda = v//enlazamos el nuevo nodo hacia la izquierda
}
func (m *matriz) existeVertical(producto *ListaInt, x int, y int) {
	m.lst_h.insert(x)
	h := m.lst_h.search(x)
	v := m.lst_v.search(y)
	nuevo := nodoMatriz(x,y,producto)
	agregado := false
	aux := v.derecha
	var cabecera int
	for aux != nil {
		cabecera = aux.headerX()
		if cabecera < x {
			aux = aux.derecha
		} else {
			nuevo.derecha = aux
			nuevo.izquierda = aux.izquierda
			aux.izquierda.derecha = nuevo
			aux.izquierda = nuevo
			agregado = true
			break
		}
	}
	if agregado == false {
		aux = v.derecha
		for aux.derecha != nil {
			aux = aux.derecha
		}
		nuevo.izquierda = aux
		aux.derecha = nuevo
	}
	nuevo.arriba = h
	h.abajo = nuevo
}
func (m *matriz) existeHorizontal(producto *ListaInt, x int, y int) {
	m.lst_v.insert(y)
	h := m.lst_h.search(x)
	v := m.lst_v.search(y)
	nuevo := nodoMatriz(x,y,producto)
	agregado := false
	aux := h.abajo
	var cabecera int
	for aux != nil {
		cabecera = aux.headerY()
		if cabecera < y {
			aux = aux.abajo
		} else {
			nuevo.abajo = aux
			nuevo.arriba = aux.arriba
			aux.arriba.abajo = nuevo
			aux.arriba = nuevo
			agregado = true
			break
		}
	}
	if agregado == false {
		aux = h.abajo
		for aux.abajo != nil {
			aux = aux.abajo
		}
		nuevo.arriba = aux
		aux.abajo = nuevo
	}
	nuevo.izquierda = v
	v.derecha = nuevo
}
func (m *matriz) existen(producto *ListaInt, x int, y int) {
	h := m.lst_h.search(x)
	v := m.lst_v.search(y)
	nuevo := nodoMatriz(x,y,producto)
	agregado := false
	aux := v.derecha
	var cabecera int
	for aux != nil {
		cabecera = aux.headerX()
		if cabecera < x {
			aux = aux.derecha
		} else {
			nuevo.derecha = aux
			nuevo.izquierda = aux.izquierda
			aux.izquierda.derecha = nuevo
			aux.izquierda = nuevo
			agregado = true
			break
		}
	}
	if agregado == false {
		aux = v.derecha
		for aux.derecha != nil {
			aux = aux.derecha
		}
		nuevo.izquierda = aux
		aux.derecha = nuevo
	}
	agregado = false
	aux = h.abajo
	for aux != nil {
		cabecera = aux.headerY()
		if cabecera < y {
			aux = aux.abajo
		} else {
			nuevo.abajo = aux
			nuevo.arriba = aux.arriba
			aux.arriba.abajo = nuevo
			aux.arriba = nuevo
			agregado = true
			break
		}
	}
	if agregado == false {
		aux = h.abajo
		for aux.abajo != nil {
			aux = aux.abajo
		}
		nuevo.arriba = aux
		aux.abajo = nuevo
	}
}
func (m *matriz) print_vertical() {
	cabecera := m.lst_v.first
	for cabecera != nil {
		aux := cabecera.derecha
		for aux != nil {
			aux.print()
			aux = aux.derecha
		}
		cabecera = cabecera.siguiente
	}
}
func (n *nodo) print(){
	fmt.Println("x: ", n.x, "y: ", n.y)
}
func (m *matriz) print_horizontal() {
	cabecera := m.lst_h.first
	for cabecera != nil {
		aux := cabecera.abajo
		for aux != nil {
			aux.print()
			aux = aux.abajo
		}
		cabecera = cabecera.siguiente
	}
}
func (m *matriz) existe(x int, y int)bool{
	cabecera := m.lst_h.first
	for cabecera != nil {
		aux := cabecera.abajo
		for aux != nil {
			if aux.x == x && aux.y == y{
				return true
			}
			aux = aux.abajo
		}
		cabecera = cabecera.siguiente
	}
	return false
}
func (m *matriz) get(x, y int) *nodo{
	cabecera := m.lst_h.first
	for cabecera != nil {
		aux := cabecera.abajo
		for aux != nil {
			if aux.x == x && aux.y == y{
				return aux
			}
			aux = aux.abajo
		}
		cabecera = cabecera.siguiente
	}
	return nil
}
func (m *matriz) getSizeCol() int{
	cabecera := m.lst_v.first
	size := 0
	for cabecera != nil{
		size++
	}
	return size
}

func (m *matriz) graphMatrix() string{
	cadena := "node0 [label=\"Calendario\", fillcolor = brown1];\n"
	//DEFINIMOS LA PRIMERA FILA
	tempH := m.lst_h.first
	for tempH != nil{
		cadena = cadena + "nodex" + strconv.Itoa(tempH.header) + " [label=\"Dia " + strconv.Itoa(tempH.header) + "\", fillcolor = burlywood1];\n"
		tempH = tempH.siguiente
	}
	//DEFINIMOS LA PRIMERA COLUMNA
	tempV := m.lst_v.first
	for tempV != nil{
		cadena = cadena + "nodey" + strconv.Itoa(tempV.header) + " [label=\"" + MPosiciones[0][tempV.header][0].Departamento + "\", fillcolor = burlywood1];\n"
		tempV = tempV.siguiente
	}
	//DEFINIMOS LOS PEDIDOS
	cabecera := m.lst_h.first
	for cabecera != nil{
		temp := cabecera.abajo
		for temp != nil{
			cadena = cadena + "node" + strconv.Itoa(temp.x) + "_" + strconv.Itoa(temp.y) + " [label=\"Pedidos\", fillcolor = cornflowerblue];\n"
			temp = temp.abajo
		}
		cabecera = cabecera.siguiente
	}
	//RECORREMOS HORIZONTALMENTE:
	cabecera = m.lst_h.first
	cadena = cadena + "node0 -> nodex" + strconv.Itoa(cabecera.header) + "[ dir=both];\n"
	for cabecera != nil{
		if cabecera.siguiente != nil{
			cadena = cadena + "nodex" + strconv.Itoa(cabecera.header) + " -> nodex" + strconv.Itoa(cabecera.siguiente.header) + "[dir=both];\n"
		}
		cadena = cadena + "nodex" + strconv.Itoa(cabecera.header)
		temp := cabecera.abajo
		for temp != nil{
			cadena = cadena + " -> node" + strconv.Itoa(temp.x) + "_" + strconv.Itoa(temp.y)
			temp = temp.abajo
		}
		cadena = cadena + "[dir=both];\n"
		cabecera = cabecera.siguiente
	}
	//RECORREMOS VERTICALMENTE:
	cabecera = m.lst_v.first
	cadena = cadena + "node0 -> nodey" + strconv.Itoa(cabecera.header) + "[ dir=both];\n"
	for cabecera != nil{
		if cabecera.siguiente != nil{
			cadena = cadena + "nodey" + strconv.Itoa(cabecera.header) + " -> nodey" + strconv.Itoa(cabecera.siguiente.header) + "[dir=both];\n"
		}
		cadena = cadena + "nodey" + strconv.Itoa(cabecera.header)
		temp := cabecera.derecha
		for temp != nil{
			cadena = cadena + " -> node" + strconv.Itoa(temp.x) + "_" + strconv.Itoa(temp.y)
			temp = temp.derecha
		}
		cadena = cadena + "[constraint=false, dir=both];\n"
		cabecera = cabecera.siguiente
	}
	//POR ULTIMO SOLO RANKEAMOS HORIZONTAL
	tempH = m.lst_h.first
	cadena = cadena + "{ rank=same; node0; "
	for tempH != nil{
		cadena = cadena + "nodex" + strconv.Itoa(tempH.header) + "; "
		tempH = tempH.siguiente
	}
	cadena = cadena + "}\n"
	//POR ULTIMO SOLO RANKEAMOS VERTICAL
	cabecera = m.lst_v.first
	for cabecera != nil{
		cadena = cadena + "{ rank=same; nodey" + strconv.Itoa(cabecera.header) + ";"
		temp := cabecera.derecha
		for temp != nil{
			cadena = cadena + "node" + strconv.Itoa(temp.x) + "_" + strconv.Itoa(temp.y) + ";"
			temp = temp.derecha
		}
		cadena = cadena + "};\n"
		cabecera = cabecera.siguiente
	}
	return cadena
}

/*****************************************************************************/
//CLASE ARBOL AVL
type nodoAVL struct {
	indice   int
	meses *ListaM
	altura   int
	izq, der *nodoAVL
}
func newNodo(indice int, meses *ListaM) *nodoAVL {
	return &nodoAVL{indice, meses, 0, nil, nil}
}
type AVL struct {
	raiz *nodoAVL
}
func NewAVL() *AVL {
	return &AVL{nil}
}
func max(val1 int, val2 int) int {
	if val1 > val2 {
		return val1
	}
	return val2
}
func altura(temp *nodoAVL) int {
	if temp != nil {
		return temp.altura
	}
	return -1
}
func rotacionIzquierda(temp **nodoAVL) {
	aux := (*temp).izq
	(*temp).izq = aux.der
	aux.der = *temp
	(*temp).altura = max(altura((*temp).der), altura((*temp).izq)) + 1
	aux.altura = max(altura((*temp).izq), (*temp).altura) + 1
	*temp = aux
}
func rotacionDerecha(temp **nodoAVL) {
	aux := (*temp).der
	(*temp).der = aux.izq
	aux.izq = *temp
	(*temp).altura = max(altura((*temp).der), altura((*temp).izq)) + 1
	aux.altura = max(altura((*temp).der), (*temp).altura) + 1
	*temp = aux
}
func rotacionDobleIzquierda(temp **nodoAVL) {
	rotacionDerecha(&(*temp).izq)
	rotacionIzquierda(temp)
}
func rotacionDobleDerecha(temp **nodoAVL) {
	rotacionIzquierda(&(*temp).der)
	rotacionDerecha(temp)
}
func insert(indice int, meses *ListaM, root **nodoAVL) {
	if (*root) == nil {
		*root = newNodo(indice, meses)
		return
	}
	if indice < (*root).indice {
		insert(indice, meses, &(*root).izq)
		if (altura((*root).izq) - altura((*root).der)) == -2 {
			if indice < (*root).izq.indice{
				rotacionIzquierda(root)
			}else{
				rotacionDobleIzquierda(root)
			}
		}
	}else if indice > (*root).indice{
		insert(indice, meses, &(*root).der)
		if (altura((*root).der) - altura((*root).izq)) == 2{
			if indice > (*root).der.indice {
				rotacionDerecha(root)
			}else{
				rotacionDobleDerecha(root)
			}
		}
	}else{
		fmt.Println("Ya se inserto el indice")
	}

	(*root).altura = max(altura((*root).izq), altura((*root).der))+1
}
func (avl *AVL) Insertar(indice int, meses *ListaM) {
	insert(indice, meses, &avl.raiz)
}
func (avl *AVL) Print(){
	inOrden(avl.raiz)
}
func inOrden(temp *nodoAVL)  {
	if temp != nil {
		inOrden(temp.izq)
		fmt.Println("Index: ", temp.indice)
		inOrden(temp.der)
	}
}
func existe(temp *nodoAVL, indice int)bool{
	if temp != nil {
		if temp.indice == indice{
			return true
		} else {
			if (existe(temp.izq, indice) || existe(temp.der, indice)){
				return true
			}
		}
	}
	return false
}
func get(temp *nodoAVL, indice int)*nodoAVL{
	if temp != nil{
		if temp.indice == indice{
			return temp
		} else {
			if (get(temp.izq, indice) != nil){
				return get(temp.izq, indice)
			} else {
				return get(temp.der, indice)
			}
		}
	}
	return nil
}
func graph(temp *nodoAVL)string{
	etiqueta := ""
	if temp.izq == nil && temp.der == nil{
		etiqueta = "nodo" + strconv.Itoa(temp.indice) + " [ label =\"" + strconv.Itoa(temp.indice) +"\"];\n"
	} else {
		etiqueta = "nodo" + strconv.Itoa(temp.indice) + " [ label =\"" + strconv.Itoa(temp.indice) +"\"];\n"
	}
	if temp.izq != nil{
		etiqueta = etiqueta + graph(temp.izq) + "nodo" + strconv.Itoa(temp.indice) + "->nodo" + strconv.Itoa(temp.izq.indice) + "\n"
	}
	if temp.der != nil{
		etiqueta = etiqueta + graph(temp.der) + "nodo" + strconv.Itoa(temp.indice) + "->nodo" + strconv.Itoa(temp.der.indice) + "\n"
	}
	return etiqueta
}