# Generowanie faktur

W tym rozdziale zajmiemy się generowaniem faktur. Program obsługuje kilka formatów plików wejściowych.

```shell
./ksef generate
Usage of generate:
  -d string
    	łańcuch znaków rozdzielający pola (tylko dla CSV) (default ",")
  -e string
    	użyj pliku z konwersją znaków (tylko dla CSV)
  -f string
    	nazwa pliku do przetworzenia
  -g string
    	nazwa generatora (default "fa-2")
  -o string
    	nazwa katalogu wyjściowego
  -s string
    	Nazwa skoroszytu do przetworzenia (tylko dla XLSX)
```

## Uwagi techniczne

::: info
Kolejność pól nie ma znaczenia. Program podczas konwersji sortuje pola według schematu przewidzianego przez ministerstwo
:::

::: info
Większość z pól przewidzianych przez ministerstwo i tak jest opcjonalna więc w plikach źródłowych uzupełniaj tylko te których potrzebujesz
:::

::: info
Korzystaj z mnemoników. Zamiast mało przyjaznych nazw pól takich jak `P_7`, `P_12` możesz posłużyć się mnemonikami które program w locie przetłumaczy na oczekiwane przez ministerstwo:

| mnemonik         | pole | znaczenie               |
| ---------------- | ---- | ----------------------- |
| item             | P_7  | Nazwa towaru / usługi   |
| units            | P_8A | Jednostka miary         |
| quantity         | P_8B | Ilość                   |
| unit-price-net   | P_9A | Cena jednostkowa netto  |
| unit-price-gross | P_9B | Cena jednostkowa brutto |
| vat-rate         | P_12 | Stawka VAT              |

:::

## CSV z sekcjami

Szczerze mówiąc nie wiem czy uprawnione jest nazywanie tego formatu CSV jako że każda z sekcji ma inną liczbę komórek. Z drugiej strony format ten przemyślałem dla integracji które nie będą posiadać jakichkolwiek bibliotek i będą generować dane źródłowe "na piechotę" (np. programy księgowe dla DOS). Z tego samego powodu konwersja z CSV wspiera także tablicę konwersji polskich ogonków

Przykładowy plik CSV znajdziesz w repozytorium projektu: https://github.com/toudi/ksef/tree/master/przykladowe-pliki-wejsciowe/csv

CSV z sekcjami wygląda w ten sposób, że na początku w pliku umieszczamy wszelkie informachje powtarzalne dla każdej z faktur (czyli np. dane wystawcy czy nazwę systemu raportowaną do KSeF) a następnie każda z sekcji `faktura.fa` rozpoczyna deklarację nowej faktury. Tym samym, w jednym pliku CSV możemy zdeklarować wiele faktur a następnie za pomocą jednej komendy wszystkie je wrzucić do KSeF

Przykłady wywołania:

```shell
./ksef generate -f plik.csv -o katalog-wyjsciowy
./ksef generate -f plik.csv -e ogonki.txt -o katalog-wyjsciowy
./ksef generate -f plik.csv -d ';' -o katalog-wyjsciowy
```

::: info
Jeśli Twój system wejściowy zapisuje kwoty za pomocą liczb całkowitych, możesz je równiez w ten sposób wyeksportować. W tym celu oprócz wyeksportowania wartości w wybranym przez Ciebie polu utwórz kolejne z dopiskiem `.decimal-places` i wstaw tam mnożnik. Dla przykładu, zapis:<br />

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
:::

## XLSX / Excell 2007+ / Libreoffice

Ten format przewidziałem dla integracji gdzie źródłem danych są faktury wystawiane w arkuszu kalkulacyjnym.

## YAML

Tu dochodzimy do formatu gdzie integracja najprawdopodobniej umożliwia zastosowanie biblioteki generującej dane wyjściowe

::: info
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

:::
