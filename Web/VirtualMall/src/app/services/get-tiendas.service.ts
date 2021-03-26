import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http'
import { servidor } from "../apiURL/servidor"
import { Observable } from 'rxjs';

import { Tienda } from 'src/app/models/tienda/tienda';
import { Producto } from 'src/app/models/productos/producto'
import { Item } from '../models/item/item';

@Injectable({
  providedIn: 'root'
})
export class GetTiendasService {

  constructor(private http: HttpClient) {
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json'
      }),
    }
  }


  getListaTiendas():Observable<Tienda[]>{
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json'
      }),
    }
    return this.http.get<Tienda[]>(servidor + 'getTiendas', httpOptions)
  }

  getListaProductos(tiendita):Observable<Producto[]>{
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    }
    return this.http.post<Producto[]>(servidor + 'getProductos', tiendita, httpOptions)
  }

  getTiendaEspecifica(tiendita):Observable<any>{
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    }
    return this.http.post<any>(servidor + 'TiendaEspecifica', tiendita, httpOptions)
  }

  setItemCarrito(item):Observable<Item[]>{
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    }
    return this.http.post<Item[]>(servidor + 'setItemCarrito', item, httpOptions)
  }

  DeleteItemCarrito(item):Observable<any>{
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    }
    return this.http.post<any>(servidor + 'DeleteItemCarrito', item, httpOptions)
  }

  graficarPedidos(str):Observable<string>{
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    }
    return this.http.post<string>(servidor + 'graficarPedidos', str, httpOptions)
  }

  pedidosCargar(archivo):Observable<string>{
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
      }),
    }
    return this.http.post<string>(servidor + 'cargarPedidos', archivo, httpOptions)
  }

}


