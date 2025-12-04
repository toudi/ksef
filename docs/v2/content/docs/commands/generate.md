---
linkTitle: generate
---

# `generate`

Komenda służy do wygenerowania pliku (plików XML) faktury na podstawie jednego z obsługiwanych formatów wejściowych. Użyj jej, jeśli nie chcesz wysyłać faktury a jedynie sprawdzić jej postać w formacie XML.

Przykładowe wywołanie:

```sh
./ksef generate -o wyjscie plik-wejsciowy
```

Powyższa komenda sparsuje `plik-wejsciowy` i zapisze wszystkie znajdujące się w nim faktury do katalogu `wyjscie`. Każda faktura zostanie zapisana w postaci pliku XML o nazwie `invoice-{nr}.xml` gdzie `{nr}` oznacza kolejny numer tj. 0, 1, 2 itd.

## Uwagi techniczne

{{< callout type="info" >}}
Kolejność pól nie ma znaczenia. Program podczas konwersji sortuje pola według schematu przewidzianego przez ministerstwo
{{< /callout >}}

{{< callout type="info" >}}
Większość z pól przewidzianych przez ministerstwo i tak jest opcjonalna więc w plikach źródłowych uzupełniaj tylko te których potrzebujesz
{{< /callout >}}

{{< callout >}}
Korzystaj z mnemoników. Zamiast mało przyjaznych nazw pól takich jak `P_7`, `P_12` możesz posłużyć się mnemonikami które program w locie przetłumaczy na oczekiwane przez ministerstwo:

| mnemonik           | pole | znaczenie               |
| ------------------ | ---- | ----------------------- |
| `item`             | P_7  | Nazwa towaru / usługi   |
| `units`            | P_8A | Jednostka miary         |
| `quantity`         | P_8B | Ilość                   |
| `unit-price-net`   | P_9A | Cena jednostkowa netto  |
| `unit-price-gross` | P_9B | Cena jednostkowa brutto |
| `vat-rate`         | P_12 | Stawka VAT              |

{{< /callout >}}

{{< callout type="warning" >}}
Uważaj na stosowanie stawki VAT "NP" (zazwyczaj stosowanej w przypadku eksportu usług). KSeF przewiduje dwie sytuacje do jej stosowania.

- na podstawie art. 100 ust 1 pkt 4 ustawy (po prostu określ stawkę jako "np II")
- z wyłączeniem art. 100 ust 1 pkt 4 ustawy (określ stawkę jako "np I")

Przewidziałem też określenie stawki jako "np" oraz dodatkowe pole "except"

Przykłady wierszy:

```csv
"Sekcja";"Faktura.Fa.FaWiersze.FaWiersz";
"P_7";"P_8A";"P_8B";"P_9B";"P_12";"P_9B.decimal-places";
"Nazwa towaru";"szt";150;"123";23;2;;
"NP na podstawie ustawy";szt;45;"1.23";"np II";;
"NP z wyłączeniem ustawy";szt;90;"1.23";"np I";;
```

```yaml
- FaWiersz:
    item: NP z wyłączeniem ustawy
    units: m2
    quantity: "1.02"
    unit-price-net:
      value: "1.23"
      decimal-places: 4
    vat-rate: np
    vat-rate.except: 1
```

{{< /callout >}}

# Obsługiwane formaty wejściowe

## CSV (z sekcjami)

Szczerze mówiąc nie wiem czy uprawnione jest nazywanie tego formatu CSV jako że każda z sekcji ma inną liczbę komórek. Z drugiej strony format ten przemyślałem dla integracji które nie będą posiadać jakichkolwiek bibliotek i będą generować dane źródłowe "na piechotę" (np. programy księgowe dla DOS). Z tego samego powodu konwersja z CSV wspiera także tablicę konwersji polskich ogonków

Przykładowy plik wejściowy znajdziesz tutaj: [szablon.csv](https://github.com/toudi/ksef/blob/master/przykladowe-pliki-wejsciowe/csv/szablon.csv)

CSV z sekcjami wygląda w ten sposób, że na początku w pliku umieszczamy wszelkie informachje powtarzalne dla każdej z faktur (czyli np. dane wystawcy czy nazwę systemu raportowaną do KSeF) a następnie każda z sekcji `faktura.fa` rozpoczyna deklarację nowej faktury. Tym samym, w jednym pliku CSV możemy zdeklarować wiele faktur a następnie za pomocą jednej komendy wszystkie je wrzucić do KSeF

{{< callout >}}
Jeśli Twój system wejściowy zapisuje kwoty za pomocą liczb całkowitych (tj. w groszach), możesz je równiez w ten sposób wyeksportować. W tym celu oprócz wyeksportowania wartości w wybranym przez Ciebie polu utwórz kolejne z dopiskiem `.decimal-places` i wstaw tam mnożnik. Dla przykładu, zapis:<br />

```
"P_9B";"P_9B.decimal-places";
"123";2;
```

Oznaczać będzie liczbę `1.23` (123 / 10 do potęgi 2)

podczas gdy zapis

```
"P_9B";"P_9B.decimal-places";
"123";4;
```

Oznaczać będzie liczbę `0.0123`
{{< /callout >}}

Przydatne parametry podczas konwersji

| flaga             | flaga (skrót) | znaczenie                                 |
| ----------------- | ------------- | ----------------------------------------- |
| `--csv.delimiter` | `-d`          | separator pól. Domyślnie: `,`             |
| `--csv.encoding`  | `-e`          | plik z konwersją strony kodowej do UTF-8. |

Przykładowe pliki konwersji:

- [cp852-dos.txt](https://github.com/toudi/ksef/blob/master/przykladowe-pliki-wejsciowe/csv/cp852-dos.txt)
- [win1250.txt](https://github.com/toudi/ksef/blob/master/przykladowe-pliki-wejsciowe/csv/win1250.txt)

## XLSX (arkusz kalkulacyjny typu Microsoft Open XML)

Przydatne parametry podczas konwersji

Przykładowy plik wejściowy znajdziesz tutaj: [przyklad.xlsx](https://github.com/toudi/ksef/blob/dbed687a3ead4e6346dabca4c008df0b8e63d5cb/przykladowe-pliki-wejsciowe/przyklad.xlsx)

| flaga          | flaga (skrót) | znaczenie                                                                                         |
| -------------- | ------------- | ------------------------------------------------------------------------------------------------- |
| `--xlsx.sheet` | `-s`          | nazwa skoroszytu. Jeśli nie zostanie podana, program użyje pierwszego skoroszytu jako bazy faktur |

## yaml

w tym formacie możesz zawrzeć zarówno pojedyńczą fakturę jak i kilka. Oto przykładowe pliki wejściowe:

- [single.yaml](https://github.com/toudi/ksef/blob/master/przykladowe-pliki-wejsciowe/yaml/single.yaml)
- [bulk.yaml](https://github.com/toudi/ksef/blob/master/przykladowe-pliki-wejsciowe/yaml/bulk.yaml)

{{< callout >}}
YAML umożliwia zapisywanie liczb zmiennoprzecinkowych. Tym niemniej, bezpieczniejszym sposobem może być albo wyeksportowanie kwoty jako string albo w formie bazowej. Sprowadzając rzecz do konkretów, kwotę `1.23` możesz zapisać w następujące sposoby:

```yaml
unit-price-net: 1.23
```

```yaml
unit-price-net: "1.23"
```

```yaml
unit-price-net:
  value: 123
  decimal-places: 2
```

{{< /callout >}}
