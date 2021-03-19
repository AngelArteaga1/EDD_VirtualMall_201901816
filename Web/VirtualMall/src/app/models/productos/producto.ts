export class Producto {
    
    Nombre: string
    Codigo: number
    Descripcion: string
    Precio: Number
    Cantidad: number
    Imagen: string

    constructor(_Nombre: string, _Codigo: number, _Descripcion: string, _Precio: number, _Cantidad: number, _Imagen: string){
        this.Nombre = _Nombre
        this.Codigo = _Codigo
        this.Descripcion = _Descripcion
        this.Precio = _Precio
        this.Cantidad = _Cantidad
        this.Imagen = _Imagen
    }


}
