#let upo = xml("upo.xml").first()

#import "../common/xml-utils.typ": children, extract
#import "../common/colors.typ": light-gray, table-border


#set text(font: "CMU Sans Serif")
#set par(spacing: 0.4em)
#set page(
  paper: "a4",
  flipped: true,
  margin: (x: 1cm, bottom: 1cm),
  header: [
    #grid(
      columns: (1fr, auto),
      [#align(horizon, [*Krajowy system #text(fill: red, "e")-Faktur*])],
      grid(
        align: right,
        row-gutter: 0.4em,
        columns: 1,
        [Nazwa pełna podmiotu, któremu doręczono dokument elektroniczny: *Ministerstwo Finansów*],
        [Informacja o dokumencie: *Dokument został zarejestrowany w systemie teleinformatycznym Ministerstwa Finansów*],
      ),
    )
    #line(length: 100%)
  ],
)

#set heading(numbering: "1.")

#text(size: 1.2em, [*Urzędowe poświadczenie odbioru dokumentu elektronicznego KSeF*])
#v(2em)

#import "./components/preamble.typ": preamble

#{ preamble(upo) }

#v(1fr)

#let invoices = children(upo, "Dokument")


#import "./components/invoices.typ": invoices-table

#{ invoices-table(invoices) }
