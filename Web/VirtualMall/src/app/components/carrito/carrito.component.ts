import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Item } from 'src/app/models/item/item'
import { GetTiendasService } from "../../services/get-tiendas.service"
import { Pedido } from "src/app/models/pedidos/pedido"
import { CodigoX } from "src/app/models/pedidos/codigo-x"
import { ArchivoPedidos } from "src/app/models/pedidos/archivo-pedidos"


@Component({
  selector: 'app-carrito',
  templateUrl: './carrito.component.html',
  styleUrls: ['./carrito.component.css']
})
export class CarritoComponent implements OnInit {

  item :Item
  comment: string = ""
  items: Item[]
  precioTotal: string = "Precio Total: Q0.00"
  estado: string = ""

  constructor(private rutaActiva: ActivatedRoute, private TiendasService: GetTiendasService) { 
    this.item = {
      Tienda: this.rutaActiva.snapshot.params.tienda,
      Departamento: this.rutaActiva.snapshot.params.departamento,
      Calificacion: Number(this.rutaActiva.snapshot.params.calificacion),
      Producto: this.rutaActiva.snapshot.params.producto,
      Descripcion: this.rutaActiva.snapshot.params.descripcion,
      Imagen: this.rutaActiva.snapshot.params.imagen,
      Precio: Number(this.rutaActiva.snapshot.params.precio),
      Codigo: Number(this.rutaActiva.snapshot.params.codigo)
    }
    this.estado = this.rutaActiva.snapshot.params.estado
    this.rutaActiva.params.subscribe(
      (params: Params) => {
        this.item.Tienda = params.tienda;
        this.item.Departamento = params.departamento;
        this.item.Calificacion = Number(params.calificacion);
        this.item.Producto = params.producto;
        this.item.Descripcion = params.descripcion;
        this.item.Imagen = params.imagen;
        this.item.Precio = Number(params.precio);
        this.item.Codigo = Number(params.codigo);
        this.estado = params.estado;
      }
    );

    console.log(this.item)
    
      if (this.estado == "delete"){
        this.TiendasService.DeleteItemCarrito(this.item).subscribe((dataList:any)=>{
          console.log("ELIMINADO CORRECTAMENTE")
        }, (err) => {
          console.log("Adios bb")
        })
        this.item.Producto = "_"
        this.item.Tienda = "_"
        this.item.Departamento = "_"
        this.item.Producto = "_"
        this.item.Descripcion = "_"
        this.item.Calificacion = -1
        this.item.Codigo = -1
        this.item.Precio = -1
      }

    //ENVIAR LOS DATOS OBTENIDOS Y CARGAR LA PAGINA PAPA
    this.TiendasService.setItemCarrito(this.item).subscribe((dataList:any)=>{
      console.log("FUNCIONO CORRECTAMENTE")
      this.items = dataList
      //recorremos todo papi
      let total = 0
      for (let elemento of this.items) {
        this.comment = this.comment + '<div class="card bg-light mb-3">'+
        '<h5 class="card-header">' + elemento.Producto + '</h5>'+
        '<div class="card-body row">'+
          '<img class="card-img-top col-2" src="' + elemento.Imagen + '" alt="Card image cap">'+
          '<div class="col-10">'+
            '<h5 class="card-title">' + elemento.Descripcion + '</h5>'+
            '<p class="card-text">Precio: Q' + elemento.Precio + '</p>'+
            '<a href="http://localhost:4200/carrito/' + elemento.Tienda + 
            '/' + elemento.Departamento + '/' + elemento.Calificacion + 
            '/' + elemento.Producto + '/' + elemento.Descripcion + '/_' + 
            '/' + elemento.Precio + '/' + elemento.Codigo + '/delete" class="btn btn-danger">Eliminar</a>'+
          '</div>'+
        '</div>'+
      '</div>'
      total = total + elemento.Precio
      }
      this.precioTotal = "Precio Total: Q" + total
    }, (err) => {
      console.log("Adios bb")
    })

  }

  CargarPedidos(){
    //obtenemos la fecha de hoy
    let today = new Date().toLocaleDateString()
    var spl = today.split("/", 3);
    let date = spl[0] + '-' + spl[1] + '-' + spl[2] 
    console.log(date)
    let pedidos: Pedido[] = new Array(this.items.length) 
    let i = 0
    for (let elemento of this.items) {
      let codigos: CodigoX[] = new Array(1) 
      let codigo = new CodigoX(elemento.Codigo)
      codigos[0] = codigo
      let pedido = new Pedido(date, elemento.Tienda, elemento.Departamento, elemento.Calificacion, codigos)
      pedidos[i] = pedido
      i++
    }
    let archivo = new ArchivoPedidos(pedidos)

    this.TiendasService.pedidosCargar(archivo).subscribe((dataList:any)=>{
      console.log("PEDIDO REALIZADO CORRECTAMENTE")
    }, (err) => {
      console.log("EL PEDIDO NO SE REALIZO CORRECTAMENTE")
    })
  }

  ngOnInit(): void {
  }

}
