import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http'
import { servidor } from "../apiURL/servidor"
import { Observable } from 'rxjs';

import { Tienda } from 'src/app/models/tienda/tienda';

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

  getListaProductos(tiendita):Observable<any>{
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json'
      }),
    }
    return this.http.post<any>(servidor + 'getProductos', tiendita, httpOptions)
  }

}
