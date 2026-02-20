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

  if vat-field.match(regex("^[0-9]+$")) != none {
    vat-multiplier = decimal(vat-field)
    vat-divisor = vat-multiplier / decimal(100)
  }

  let gross-amount = decimal(0)
  let net-amount = decimal(0)
  let vat-amount = decimal(0)

  if net-price-node.len() > 0 {
    // we have to calculate net -> gross
    net-price = calc.round(decimal(net-price-node.first().children.first()), digits: 2)
    net-amount = calc.round(net-price * quantity, digits: 2)
    vat-amount = calc.round(net-amount * vat-divisor, digits: 2)

    gross-amount = net-amount + vat-amount
    gross-price = net-price
    if vat-multiplier > 0 {
      gross-price = calc.round(decimal(net-price * (1 + vat-divisor)), digits: 2)
    }
  } else {
    // we have to calculate gross -> net
    gross-price = calc.round(decimal(gross-price-node.first().children.first()), digits: 2)
    gross-amount = calc.round(gross-price * quantity, digits: 2)
    net-amount = calc.round(gross-amount / (1 + vat-divisor), digits: 2)
    vat-amount = gross-amount - net-amount
    net-price = gross-price
    if vat-divisor > 0 {
      net-price = calc.round(decimal(gross-price / (1 + vat-divisor)), digits: 2)
    }
  }

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

#let add-percent-if-numeric-rate(value) = {
  if value.match(regex("^[0-9]+$")) != none {
    return value + " %"
  }

  return value
}
