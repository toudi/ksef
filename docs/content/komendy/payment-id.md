# Generowanie identyfikatora płatności

Jeśli masz kilka faktur i chciałbyś aby kontrahent zapłacił Ci za nie za pomocą jednego przelewu to ministerstwo przewidziało do tego zbiorczy identyfikator płatności.

```text
./ksef payment-id
Usage of payment-id:
  -json
        Użyj formatu JSON do zapisania wyjścia
  -o string
        Plik do zapisania wyjścia
  -p string
        ścieżka do pliku rejestru
  -token string
        Token sesji interaktywnej lub nazwa zmiennej środowiskowej która go zawiera
  -yaml
        Użyj formatu YAML do zapisania wyjścia
```

Komenda `payment-id` przewiduje dwa tryby działania:

## Generowanie identyfikatora płatności

```shell
./ksef payment-id -p lokalizacja-pliku-registry.yaml 'numer faktury 1' 'numer faktury 2' ...
```

::: info
Numerem faktury może być zarówno identyfikator wewnętrzny (czyli po prostu Twój numer faktury) jak i numer KSeF
:::

Komenda zapisze do plku rejestru nową pozycję, np:

```yaml
payment-ids:
  - ksefPaymentRefNo: 1111111111-XXXXXXXX-YYYYYYYYYYYY-ZZ
    ksefInvoiceRefNumbers:
      - 1111111111-XXXXXXXX-YYYYYYYYYYYY-ZZ
      - 1111111111-XXXXXXXX-YYYYYYYYYYYY-ZZ
```

Możesz więc bezproblemowo odczytać plik rejestru i dobrać się do powyższych identyfikatorów

## Wyświetlanie listy faktur podlegających pod identyfikator

Ta opcja jest przydatna wtedy kiedy potrzebujesz wygenerować identyfikatory faktur dla podanego identyfikatora płatności. Przykład wywołania:

```shell
./ksef payment-id -p registry.yaml 1111111111-XXXXXXXX-YYYYYYYYYYYY-ZZ
invoiceRefNo,ksefInvoiceRefNo
abcabcabc,1111111111-XXXXXXXX-YYYYYYYYYYYY-ZZ
abcabcabc,1111111111-XXXXXXXX-YYYYYYYYYYYY-ZZ
```

(format CSV jest domyślnym wyjściem)

Dodatkowo możesz użyć flag `-yaml` oraz `-json` w zależności od formatu który jest Ci potrzebny.

Możesz przekierować wyjście komendy do pliku używając przełącznika `-o`
