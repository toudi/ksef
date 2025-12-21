#import "./item-price.typ": item-price

#let calc-total-amounts(items) = {
  let total-amounts = (
    net: decimal(0),
    gross: decimal(0),
    vat: decimal(0),
  )

  for item in items {
    let amounts = item-price(item)
    total-amounts.at("net") += amounts.amount.net
    total-amounts.at("gross") += amounts.amount.gross
    total-amounts.at("vat") += amounts.amount.vat
  }

  return total-amounts
}
