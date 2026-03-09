export interface Stock {
  symbol: string
  name: string
  price: number
  change: number
}

export class StockTable {
  code: string = ''
  price: number = 0

  constructor() {}
}
