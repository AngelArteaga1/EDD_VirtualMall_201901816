export class Item {

    Tienda: string
    Departamento: string
    Calificacion: number
    Producto: string
    Descripcion: string
    Imagen: string
    Precio: number
    Codigo: number

    constructor(_Tienda: string, _Departamento: string, _Calificacion: number, _Producto: string, _Descripcion: string, _Imagen: string, _Precio: number, _Codigo: number){
        this.Tienda = _Tienda
        this.Departamento = _Departamento
        this.Calificacion = _Calificacion
        this.Producto = _Producto
        this.Descripcion = _Descripcion
        this.Imagen = _Imagen
        this.Precio = _Precio
        this.Codigo = _Codigo
    }

}
