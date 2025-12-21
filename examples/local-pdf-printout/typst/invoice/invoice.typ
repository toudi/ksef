#let meta = yaml("meta.yaml")
#let invoice = xml("invoice.xml").first()

#import "../common/xml-utils.typ": *

#import "./components/qr-codes.typ": qr-codes
#import "./components/header.typ": invoice-header

#let page-header(meta) = {
  if "header" in meta {
    let left = [#h(1fr)]
    let middle = [#h(1fr)]
    let right = [#h(1fr)]
    if "left" in meta.at("header") {
      left = [#{ meta.header.left }]
    }
    if "center" in meta.at("header") {
      middle = [#{ meta.header.center }]
    }
    if "right" in meta.at("header") {
      right = [#{ meta.header.right }]
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
