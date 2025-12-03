---
linkTitle: generate
---

# `generate`

Komenda służy do wygenerowania pliku (plików XML) faktury na podstawie jednego z obsługiwanych formatów wejściowych. Użyj jej, jeśli nie chcesz wysyłać faktury a jedynie sprawdzić jej postać w formacie XML.

Przykładowe wywołanie:

```sh
./ksef generate -o wyjscie plik-wejsciowy
```

Powyższa komenda sparsuje `plik-wejsciowy` i zapisze wszystkie znajdujące się w nim faktury do katalogu `wyjscie`. Każda faktura zostanie zapisana w postaci pliku XML o nazwie `invoice-{nr}` gdzie `{nr}` oznacza kolejny numer tj. 0, 1, 2 itd.

Obsługiwane formaty wejściowe:

## CSV (z sekcjami)

Przykładowy plik wejściowy znajdziesz tutaj: [szablon.csv](https://github.com/toudi/ksef/blob/master/przykladowe-pliki-wejsciowe/csv/szablon.csv)

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
