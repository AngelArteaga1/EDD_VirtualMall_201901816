import { Component, OnInit } from '@angular/core';
import { GetTiendasService } from "../../services/get-tiendas.service"
import { ActivatedRoute, Params } from '@angular/router';
import { TiendaEspecifica } from 'src/app/models/tienda/tienda-especifica'
import { Producto } from 'src/app/models/productos/producto'


@Component({
  selector: 'app-productos',
  templateUrl: './productos.component.html',
  styleUrls: ['./productos.component.css']
})
export class ProductosComponent implements OnInit {

  store: TiendaEspecifica
  comment: string = "";
  productos: Producto[]
  logo: string =""
  //nombreTienda: string = ""
  descripcion: string = ""
  contacto: string = ""
  mostrarMensajeError = false

  constructor(private rutaActiva: ActivatedRoute, private TiendasService: GetTiendasService) { 
    this.store = {
      Departamento: this.rutaActiva.snapshot.params.departamento,
      Nombre: this.rutaActiva.snapshot.params.tienda,
      Calificacion: Number(this.rutaActiva.snapshot.params.calificacion)
    }
    this.rutaActiva.params.subscribe(
      (params: Params) => {
        this.store.Departamento = params.departamento;
        this.store.Nombre = params.tienda;
        this.store.Calificacion = Number(params.calificacion);
      }
    );

    this.TiendasService.getTiendaEspecifica(this.store).subscribe((dataList:any)=>{
      this.descripcion = dataList.Descripcion
      this.logo = dataList.Logo
      this.contacto = dataList.Contacto
    }, (err) => {
      this.mostrarMensajeError = true
      console.log("Adios bb")
    })

    this.TiendasService.getListaProductos(this.store).subscribe((dataList:Producto[])=>{
      console.log(dataList)
      this.productos = dataList;
      console.log("cantidad de la lista: ", this.productos.length)
      let contador = 0
      let i = 0
      for (let product of this.productos) {
        //PARA SABER EL RATING
        if (contador == 0){
          this.comment = this.comment + '<div class="row text-center">'
        }
        this.comment = this.comment + 
        '<div class="card col-3">'+
            '<img src="' + product.Imagen + '" class="card-img-top" alt="products">'+
            '<div class="card-body">'+
              '<h5 class="card-title">' + product.Nombre + '</h5>'+
              '<p class="card-text">' + product.Descripcion + '</p>'+
              '<p class="text-success">Precio: Q' + product.Precio + '</p>'+
              '<a role="button" class="btn btn-primary" href="http://localhost:4200/carrito/' + this.store.Nombre + 
              '/' + this.store.Departamento + '/' + this.store.Calificacion + 
              '/' + product.Nombre + '/' + product.Descripcion + '/_' + 
              '/' + product.Precio + '/' + product.Codigo + '/add">Comprar</a>'+
            '</div>'+
        '</div>'
        if (contador == 3 || i == this.productos.length - 1){
          this.comment = this.comment + '</div>\n<hr class="featurette-divider"><br>'
          contador = 0
        } else {
          contador++
        }
        i++
      }
      
    }, (err) => {
      this.mostrarMensajeError = true
      console.log("Adios bb")
    })
  }

  onClickMe(){
    console.log("HOLA bb")
  }

  ngOnInit(): void {
    
    //console.log(this.store)
  }

}
