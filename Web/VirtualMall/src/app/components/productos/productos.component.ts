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
  comment: string;
  productos: Producto[]
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
    console.log("Dale papi: ", this.store.Calificacion);
    console.log(this.store)
    console.log("Dale papi + 1: ", this.store.Calificacion + 1)
    this.TiendasService.getListaProductos(this.store).subscribe((dataList:any)=>{
      console.log(dataList)
      this.productos = dataList;
      console.log("cantidad de la lista: ", this.productos.length)
      let contador = 0
      let i = 0
      for (let product of this.productos) {
        if (contador == 0){
          this.comment = this.comment + '<div class="row text-center">'
        }
        this.comment = this.comment + 
        '<div class="card col-3">'+
                '<img src="' + product.Imagen + '" class="card-img-top" alt="...">'+
                '<div class="card-body">'+
                '<h5 class="card-title">' + product.Nombre + '</h5>'+
                '<p class="card-text">' + product.Descripcion + '</p>'+
                '<a href="#" class="btn btn-primary">Comprar</a>'+
                '</div>'
            '</div>'
        if (contador == 3 || i == this.productos.length - 1){
          this.comment = this.comment + '</div>\n<br><hr class="featurette-divider"><br>'
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

  ngOnInit(): void {
    
    //console.log(this.store)
  }

}
