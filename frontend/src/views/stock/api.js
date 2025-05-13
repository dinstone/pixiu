import { AddStock, AddTransaction, DeleteStock, DeleteTransaction, GetHolding, GetStockList, GetTransactions, UpdateStock, UpdateTransaction } from 'wailsjs/go/ipc/StockApi.js'

export default {
  getStocks: () => GetStockList(),
  addStock: data => AddStock(data),
  saveStock: data => UpdateStock(data),
  deleteStock: code => DeleteStock(code),
  getHolding: stockCode => GetHolding(stockCode),
  getTrades: holdingId => GetTransactions(holdingId),
  addTrade: data => AddTransaction(data),
  saveTrade: data => UpdateTransaction(data),
  deleteTrade: id => DeleteTransaction(id),
}
