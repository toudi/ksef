#import "../../../common/xml-utils.typ": children, extract

#let item-price(item) = {
  let quantity = decimal(extract(item, "P_8B"))
  let net-price-node = children(item, "P_9A")
  let gross-price-node = children(item, "P_9B")

  let net-price = decimal(0)
  let gross-price = decimal(0)

  let vat-multiplier = decimal(0)
  let vat-field = extract(item, "P_12")
  let vat-divisor = decimal(0)

  if vat-field.contains(regex("[0-9]+")) {
    vat-multiplier = decimal(vat-field)
    vat-divisor = (100 + vat-multiplier) / decimal(100)
  }

  if net-price-node.len() > 0 {
    // we have to calculate net -> gross
    net-price = decimal(net-price-node.first().children.first())
    gross-price = net-price
    if vat-multiplier > 0 {
      gross-price = decimal(net-price * vat-divisor)
    }
  } else {
    // we have to calculate gross -> net
    gross-price = decimal(gross-price-node).first().children.first()
    net-price = decimal(gross-price / vat-divisor)
  }
  let gross-amount = gross-price * quantity
  let net-amount = net-price * quantity
  let vat-amount = (gross-amount) - (net-amount)
  return (
    net: calc.round(net-price, digits: 2),
    gross: calc.round(gross-price, digits: 2),
    amount: (
      net: calc.round(net-amount, digits: 2),
      gross: calc.round(gross-amount, digits: 2),
      vat: calc.round(vat-amount, digits: 2),
    ),
  )
}
