---
title: PDF
---

Wydruk PDF

Program umożliwia stworzenie wydruku faktury w formacie PDF poprzez jeden z kilku dostępnych silników.

## `typst`

Na ten projekt natknąłem się przez kompletny przypadek, natomiast muszę szczerze przyznać, że niesamowicie mnie zachwycił. Umożliwia tworzenie dokumentów PDF ale zawiera masę przydatnych funkcjonalności - w tym czytanie plików xml, yaml itp itd.

W repozytorium z programem zamieszczam przykładowe implementacje wizualizacji (`examples/local-pdf-printout/typst`)

### konfiguracja silnika

| opcja              | znaczenie                                                                                                                                                                                                                                                                                                        |
| ------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `debug`            | Jeśli ustawisz opcję na `true`, to program podczas napotkanych błędów skopiuje standardowe wyjście programu `typst` do pliku z nazwą pliku konwertowanego i końcówką `-error.txt`, co pozwoli na szybszą diagnozę                                                                                                |
| `templates-dir`    | Katalog zawierający szablony typst dla renderowanych faktur oraz UPO                                                                                                                                                                                                                                             |
| `workdir`          | Jeśli wskażesz tu katalog, program skopiuje szablony z katalogu `templates` do tego katalogu. Podczas pracy, program będzie kopiować tu pliki XML UPO oraz faktur. Jeśli pozostawisz opcję niezdefiniowaną (lub ustawisz wartość na pusty string) - wtedy program będzie kopiować pliki do katalogu z szablonami |
| `upo.template`     | Nazwa szablonu do renderowania UPO                                                                                                                                                                                                                                                                               |
| `invoice.template` | Nazwa szablonu do renderowania faktury                                                                                                                                                                                                                                                                           |
| `invoice.header`   | Ustawienia nagłówka (dostępne podklucze: `left`, `center`, `right`)                                                                                                                                                                                                                                              |
| `invoice.footer`   | Ustawienia stopki (dostępne podklucze: `left`, `center`, `right`)                                                                                                                                                                                                                                                |
| `invoice.printout` | Dodatkowe zmienne które zostaną przekazane do wydruku                                                                                                                                                                                                                                                            |

#### przykładowe ustawienia

```yaml
pdf:
  - usage:
      - upo
    typst:
      debug: true
      templates-dir: examples/local-pdf-printout/typst
      upo:
        template: upo/upo.typ
  - usage:
      - invoice:issued
    typst:
      debug: true
      templates-dir: examples/local-pdf-printout/typst
      invoice:
        template: invoice/invoice.typ
        header:
          left: Lewa strona nagłówka
          center: środek nagłówka
          right: prawa strona nagłówka
        footer:
          right: WSI Pegasus
```

### działanie silnika

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

printout:
  wszystkie:
    dodatkowe: ["wartości", "które", "przekażesz"]
    a: ["które", "nie są", "przewidziane"]
    przez: ["format", "FA(3)"]

  bank-account:
    account-no: PL 1111 2222 0000 ( ... )
    name: Nazwa banku
```

### Spersonalizowane szablony

Opiszę dwa poziomy personalizacji. Pierwszy poziom polega na dodaniu warunków do ustawień silnika. Zacznijmy od przykładu bo on chyba najprościej pokaże do czego służy opcja personalizacji

```yaml
pdf:
  - usage: "*"
    cirfmf:
      templates-dir: examples/local-pdf-printout/cirfmf
      workdir: /Users/toudi/printout/cirfmf
      node-bin: /opt/homebrew/bin/node
  - usage:
      - upo
    typst:
      templates-dir: examples/local-pdf-printout/typst
      upo:
        template: upo/upo.typ
  - usage:
      - invoice:issued
    typst:
      templates-dir: examples/local-pdf-printout/typst
      invoice:
        template: invoice/invoice.typ
        footer:
          right: WSI Pegasus
        printout: {}
  - usage:
      - invoice:issued
    if: '{{ eq .buyer.nr_vat_ue "112233445566" }}'
    typst:
      debug: true
      templates-dir: /Users/toudi/projects/invoice-templates-git/typst
      invoice:
        template: invoice/invoice-dual-lang-np.typ
        printout: { no-qrcodes: true }
```

Podczas wyboru szablonu do wydruku faktury wystawionej (`invoice:issued`) program ma do wyboru trzy szablony (patrząc na tablicę powyżej - pierwszy (`usage: *`), trzeci oraz czwarty). Jeśli do renderowania zostanie wybrana faktura w której stroną kupującą jest podmiot identyfikowany przez `nr_vat_ue` o wartości `112233445566` - wówczas program wybierze ten szablon. Jeśli nie - zostanie użyty szablon numer 3, ponieważ jego zastosowanie (`usage`) jest bardziej szczegółowe niż ogólne (`*`)

Pola dostępne przy ewaluacji warunku:

| pole              | opis                                   |
| ----------------- | -------------------------------------- |
| `seller.nip`      | numer NIP wystawcy                     |
| `buyer.nip`       | Numer NIP kupującego                   |
| `buyer.nr_vat_ue` | Numer NIP kupującego z UE              |
| `buyer.nr_id`     | Numer NIP kupującego z poza obszaru UE |

Jeśli Twoja lista ustawień silnika stanie się zbyt długa, możesz rozważyć rozbicie ustawień per podmiot wystawiający. W tym celu:

1. Skorzystaj z komendy kopiowania ustawień:

   ```
   ./ksef subject-settings copy-pdf-config -n 1111111111
   ```

   Oczywiście powyższa komenda również reaguje na flagi odpowiedzialne za środowisko (`-t` / `--demo-gateway`). Spowoduje to skopiowanie ustawień z pliku konfiguracyjnego do ustawień podmiotu (`data/<środowisko>/<nip>/settings.yaml`).

1. Zmodyfikuj ustawienia w obrębie pojedynczego podmiotu

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
