import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes } from '@angular/router';
import { TiendasComponent } from "./components/tiendas/tiendas.component"
import { InicioComponent } from "./components/inicio/inicio.component"
import { ProductosComponent } from "./components/productos/productos.component"
import { CarritoComponent } from "./components/carrito/carrito.component"
import { ReportesComponent } from "./components/reportes/reportes.component"


const routes: Routes = [
  { path: '', component: InicioComponent },
  { path: 'tienda', component: TiendasComponent },
  { path: 'productos/:tienda/:departamento/:calificacion', component: ProductosComponent },
  { path: 'carrito/:tienda/:departamento/:calificacion/:producto/:descripcion/:imagen/:precio/:codigo/:estado', component: CarritoComponent },
  { path: 'reportes', component: ReportesComponent }
];

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    RouterModule.forRoot(routes)
  ],
  exports: [
    RouterModule
  ]
})
export class AppRoutingModule { }
