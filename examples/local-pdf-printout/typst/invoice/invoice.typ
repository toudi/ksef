#let meta = yaml("meta.yaml")
#let invoice = xml("invoice.xml").first()

#import "../common/xml-utils.typ": *

#import "./components/qr-codes.typ": qr-codes
#import "./components/header.typ": invoice-header

#let page-header(meta) = {
  if "page" in meta and "header" in meta.at("page") {
    let header = meta.at("page").at("header")
    let left = [#h(1fr)]
    let middle = [#h(1fr)]
    let right = [#h(1fr)]
    if "left" in header and header.left.len() > 0 {
      left = [#{ header.left }]
    }
    if "center" in header and header.center.len() > 0 {
      middle = [#{ header.center }]
    }
    if "right" in header and header.right.len() > 0 {
      right = [#{ header.right }]
    }
    grid(
      grid(
        columns: 3,
        left, middle, right,
      ),
      [#v(0.5em)],
      [#line(length: 100%)],
    )
  }
}

#set text(font: "CMU Sans Serif")
#set page(
  margin: (x: 1cm, bottom: 1cm, top: 2cm),
  header: page-header(meta),
)

#grid(
  columns: (3fr, 2fr),
  align: horizon,
  [#{ qr-codes(meta) }], [#align(right, { invoice-header(invoice) })],
)

#import "./components/participants.typ": participants

#{ participants(invoice) }

#import "./components/items.typ": items, vat-summary

#v(2em)

#{ items(invoice) }

#v(1em)

#align(end, { vat-summary(invoice) })
