#import "../../common/xml-utils.typ": children, extract
#import "../../common/colors.typ": light-gray, table-border

#let extract-auth-context(elem) = {
  let auth_type_child = elem
    .filter(e => "tag" in e and e.tag == "IdKontekstu")
    .first()
    .children
    .filter(e => "tag" in e)
    .first()

  return (
    ctx_type: auth_type_child.tag,
    value: auth_type_child.children.first(),
  )
}

#let preamble(upo) = {
  let opisPotwierdzenia = extract(upo, "OpisPotwierdzenia")
  let uwierzytelnienie = extract(upo, "Uwierzytelnienie").filter(e => "tag" in e)
  let authContext = extract-auth-context(uwierzytelnienie)

  table(
    columns: 2,
    fill: (x, _) => if x == 0 { light-gray },
    stroke: 0.5pt + table-border,

    [Numer referencyjny sesji], [#{ extract(upo, "NumerReferencyjnySesji") }],

    [Strona dokumentu UPO ], [#{ extract(opisPotwierdzenia, "Strona") }],
    [Całkowita liczba stron dokumentu UPO ], [#{ extract(opisPotwierdzenia, "LiczbaStron") }],
    [Zakres dokumentów od ], [#{ extract(opisPotwierdzenia, "ZakresDokumentowOd") }],
    [Zakres dokumentów do ], [#{ extract(opisPotwierdzenia, "ZakresDokumentowDo") }],
    [Całkowita liczba dokumentów ], [#{ extract(opisPotwierdzenia, "CalkowitaLiczbaDokumentow") }],
    [Typ kontekstu ], [#{ upper(authContext.at("ctx_type")) }],
    [Identyfikator kontekstu uwierzytelnienia ], [#{ authContext.at("value") }],
    [Skrót dokumentu uwierzytelniającego ], [#{ extract(uwierzytelnienie, "SkrotDokumentuUwierzytelniajacego") }],
    [Nazwa pliku XSD struktury logicznej dotycząca przesłanego dokumentu ],
    [#{ extract(upo, "NazwaStrukturyLogicznej") }],

    [Kod formularza przedłoonego dokumentu elektronicznego ], [#{ extract(upo, "KodFormularza") }],
  )
}
