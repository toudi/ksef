---
title: PDF
---

Wydruk PDF

Program umożliwia stworzenie wydruku faktury w formacie PDF poprzez jeden z kilku dostępnych silników.

## `typst`

Na ten projekt natknąłem się przez kompletny przypadek, natomiast muszę szczerze przyznać, że niesamowicie mnie zachwycił. Umożliwia tworzenie dokumentów PDF ale zawiera masę przydatnych funkcjonalności - w tym czytanie plików xml, yaml itp itd.

W repozytorium z programem zamieszczam przykładowe implementacje wizualizacji (`examples/local-pdf-printout/typst`)

### działanie silnika

W konfiguracji (patrz poniżej) definiujesz katalog roboczy (`workdir`). Do tego katalogu zostaną skopiowane szablony. Zalecam stworzenie tego katalogu w ramdysku (np. `/tmp`), ponieważ kopiowane tam będą tymczasowo pliki faktur oraz UPO.

W trybie renderowania UPO program skopiuje upo pod nazwą `upo.xml`

W trybie renderowania faktury program skopiuje do katalogu z szablonem:

- fakturę pod nazwą invoice.xml
- dodatkowo - program utworzy plik o nazwie meta.yaml o następującej strukturze:

```yaml
page:
  header:
    left: Lewa strona nagłówka
    right: Prawa strona nagłówka

invoice:
  ksef-ref-no: 1111111111-20251218-AABBCCDDEEFF-AA
  qr-codes:
    invoice: https://ksef-test.mf.gov.pl/client-app/invoice/1111111111/17-12-2025/aabbccddeeffgghhiijjkkllmmnnooppqqrr_112233=
    #offline: https://ksef-test.mf.gov.pl/client-app/certificate/Nip/1111111111/1111111111/0011223344556677/ ....
```

## `cirfmf` (biblioteka Ministerstwa Finansów)

Ministerstwo Finansów udostępniło swoją bibliotekę do renderowania PDF. Można jej użyć zarówno do renderowania UPO jak i faktur. Domyślnie ta biblioteka używana jest jak aplikacja webowa, dlatego dodałem do niej pomocniczy skrypt który umożliwia używanie biblioteki w trybie wsadowym.

1. Sklonuj repozytorium z biblioteką

   ```
   git clone https://github.com/CIRFMF/ksef-pdf-generator
   ```

1. Zainstaluj potrzebne biblioteki oraz skompiluj bibliotekę

   ```
   cd ksef-pdf-generator
   npm install
   npm run build
   ```

1. Skopiuj otrzymany plik (`dist/ksef-fe-invoice-converter.js`) do katalogu w którym umieścisz pozostałe skrypty. Możesz użyć mojego przykładowego katalogu (tj. `examples/local-pdf-printout/cirfmf`)

1. Pamiętaj aby w konfiguracji wpisać ścieżkę do `node` (patrz poniżej)

## Konfiguracja silników

w pliku configuracyjnym (`config.yaml`) umieść następującą treść:

```yaml
pdf:
  - usage: "*"
    cirfmf:
      templates-dir: examples/local-pdf-printout/cirfmf
      workdir: /Users/toudi/printout/cirfmf
      node-bin: /opt/homebrew/bin/node
  - usage: ["upo", "invoice:issued"]
    typst:
      debug: true
      workdir: /Users/toudi/printout/typst
      templates-dir: examples/local-pdf-printout/typst
      invoice:
        template: invoice/invoice.typ
        header:
          left: Lewa strona nagłówka
          right: Prawa strona nagłówka
        footer:
          right: Prawa strona stopki
      upo:
        template: upo/upo.typ
```

Omówię teraz pokrótce w jaki sposób program zinterpretuje konfigurację.

Po pierwsze zwróć uwagę na pole `usage`. Jest ona hierarchiczna, tzn że program najpierw szuka szczegółowego zastosowania (np. `invoice:issued`) a dopiero później bierze ogólny (`*`). W zaproponowanym przykładzie dla faktur wystawionych (`issued`) oraz dla UPO zostanie użyty szablon oparty o silnik typst, w pozostałym przypadku - silnik ministerstwa finansów.
