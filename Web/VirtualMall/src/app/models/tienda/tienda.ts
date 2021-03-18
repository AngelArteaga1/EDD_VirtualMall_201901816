export class Tienda {
    Nombre: string
    Descripcion: string
    Contacto: string
    Calificacion: number
    Logo: string
    Departamento: string
    Indice: string

    constructor(_Nombre: string, _Descripcion: string, _Contacto: string, _Calificacion: number, _Logo: string, _Departamento: string, _Indice: string){
        this.Nombre = _Nombre
        this.Descripcion = _Descripcion
        this.Contacto = _Contacto
        this.Calificacion = _Calificacion
        this.Logo = _Logo
        this.Departamento = _Departamento
        this.Indice = _Indice
    }
}
