import { Component, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';
import { Tienda } from 'src/app/models/tienda/tienda';
import { GetTiendasService } from "../../services/get-tiendas.service"

@Component({
  selector: 'app-tiendas',
  templateUrl: './tiendas.component.html',
  styleUrls: ['./tiendas.component.css']
})
export class TiendasComponent implements OnInit {

  //mensajeError: string
  mostrarMensajeError = false
  mostrarMensajeExito = false
  listaTiendas: Tienda[]
  comment: string = ""

  constructor(private TiendasService: GetTiendasService) { 
    
    this.TiendasService.getListaTiendas().subscribe((dataList:Tienda[])=>{
      console.log(dataList)
      this.listaTiendas = dataList;
      console.log("cantidad de la lista: ", this.listaTiendas.length)
      let contador = 0
      let i = 0
      for (let tienda of this.listaTiendas) {
        if (contador == 0){
          this.comment = this.comment + '<div class="row text-center">'
        }
        this.comment = this.comment + 
        '<div class="col-lg-4">'+
        '<div>' +
        //'<div class="row p-3 mb-2 bg-secondary text-white border border-light">'+
        '<div class="row">'+
        '<p class="col-6">Departamento: ' + tienda.Departamento +'</p>'+
        '<p class="col-6">Calificacion: ' + tienda.Calificacion +'</p>'+
        '</div>' +
        '<hr class="featurette-divider">' +
        '<img class="rounded-circle" src="' + tienda.Logo + '" alt="Generic placeholder image" width="140" height="140">'+
        '<h2>' + tienda.Nombre + '</h2>'+
        '<p>' + tienda.Descripcion + '</p>'+
        '<p><a class="btn btn-secondary" href="/productos/' + tienda.Nombre + '/' + tienda.Departamento + '/' + tienda.Calificacion + '" role="button">Ver Tienda &raquo;</a></p>'+
        '</div>' +
        '</div><!-- /.col-lg-4 -->'
        if (contador == 2 || i == this.listaTiendas.length - 1){
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
  }

}
