import { CodigoX } from "./codigo-x"

export class Pedido {
    Fecha: string
    Tienda: string
    Departamento: string
    Calificacion: number
    Productos: CodigoX[]

    constructor(_Fecha: string, _Tienda: string, _Departamento: string, _Calificacion: number, _Productos: CodigoX[]){
        this.Fecha = _Fecha
        this.Tienda = _Tienda
        this.Departamento = _Departamento
        this.Calificacion = _Calificacion
        this.Productos = _Productos
    }

}
