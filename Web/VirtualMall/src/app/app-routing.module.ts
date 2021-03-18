import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes } from '@angular/router';
import { TiendasComponent } from "./components/tiendas/tiendas.component"
import { InicioComponent } from "./components/inicio/inicio.component"
import { ProductosComponent } from "./components/productos/productos.component"


const routes: Routes = [
  { path: '', component: InicioComponent },
  { path: 'tienda', component: TiendasComponent },
  { path: 'productos/:tienda/:departamento/:calificacion', component: ProductosComponent }
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
