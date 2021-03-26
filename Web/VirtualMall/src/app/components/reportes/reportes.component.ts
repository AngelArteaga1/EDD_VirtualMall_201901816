import { Component, OnInit } from '@angular/core';
import {FormControl} from '@angular/forms';
import { GetTiendasService } from "../../services/get-tiendas.service"
import { Cad } from 'src/app/models/cad/cad'

@Component({
  selector: 'app-reportes',
  templateUrl: './reportes.component.html',
  styleUrls: ['./reportes.component.css']
})
export class ReportesComponent implements OnInit {

  date = new FormControl('')
  dep = new FormControl('')
  str: Cad
  constructor(private TiendasService: GetTiendasService) { }


  Graficar(){
    let Cadena = ""
    var spl = this.date.value.split("-", 3); 
    Cadena = spl[2] + '-' + spl[1] + '-' + spl[0] + '-' + this.dep.value
    this.str = {Cadena}
    console.log(Cadena) 

    this.TiendasService.graficarPedidos(this.str).subscribe((dataList:any)=>{
      console.log("GRAFICAS REALIZADAS EXITOSAMENTE")
    }, (err) => {
      console.log("NO SE REALIZARON LAS GRAFICAS")
    })

  }

  ngOnInit(): void {
  }

}
