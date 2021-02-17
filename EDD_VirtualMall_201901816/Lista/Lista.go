package Lista

import "fmt"

//CLASE NODO
type Nodo struct{
	Siguiente *Nodo
	Anterior *Nodo
	Dato string
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
func (l *Lista) Add(valor string){
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
		fmt.Print("<--[", Aux.Dato, "]-->")
		Aux = Aux.Siguiente
	}
	fmt.Println()
	fmt.Println("Tamaño de la lista: ", l.len)
}
