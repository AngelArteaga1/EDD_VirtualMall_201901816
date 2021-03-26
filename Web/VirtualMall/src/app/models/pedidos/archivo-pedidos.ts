import { Pedido } from "./pedido";

export class ArchivoPedidos {
    Pedidos: Pedido[]
    constructor(_Pedidos: Pedido[]){
        this.Pedidos = _Pedidos
    }
}
