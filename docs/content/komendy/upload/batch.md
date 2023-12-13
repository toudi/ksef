# Wysyłka wsadowa (batch)

Wysyłka w trybie wsadowym sprowadza się do następujących kroków:

1. Spakowanie źródłowych plików XML do archiwum ZIP
2. Wygenerowania klucza szyfrującego do AES
3. Zaszyfrowanie klucza z pkt. 2 za pomocą klucza publicznego dostarczonego przez Ministerstwo Finansów
4. Zaszyfrowanie archiwum wygeneowanego w pkt. 1 kluczem wygenerowanym w pkt. 2
5. Wygenerowanie pliku metadanych który zawiera informację o sumie kontrolnej archiwum oraz spakowanego archiwum
6. Podpisanie pliku z pkt. 5 za pomocą osadzonego certyfikatu
7. Wysłanie podpisanego pliku do KSeF

Program zajmuje się punktami 1 - 5

Aby skorzystać z wysyłki wsadowej, wydaj następujące komendy:

## metadata

```text
./ksef metadata
Usage of metadata:
  -p string
        ścieżka do wygenerowanych plików
  -t    użyj bramki testowej
```

Przykładowe wywołanie dla bramki testowej

```shell
./ksef metadata -p katalog-z-plikami-xml -t
```

Program wygeneruje następujące pliki:

| plik             | znaczenie                                                                                                                                    |
| ---------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| metadata.xml     | surowy plik metadanych który należy podpisać. Podpisanego pliku użyjesz w kolejnym kroku (`upload`)                                          |
| metadata.zip     | surowy plik archiwum, nie jest on wysyłany na serwer                                                                                         |
| metadata.zip.aes | plik archiwum zaszyfrowany odpowiednim kluczem ministerstwa, zależnym od wybranego trybu (testowy / produkcja) - to ten plik jest przesyłany |

## podpisanie pliku

Możesz posłużyć się bramką rządową do podpisywania plików [ https://moj.gov.pl/nforms/signer/upload?xFormsAppName=SIGNER&xadesPdf=true ]. Po podpisaniu bramka zwróci Ci podpisany plik metadanych. Zapisz go pod inną nazwą (na wypadek ewentualnych błędów) np. `metadata-signed.xml`

# upload

```text
./ksef upload
Usage of upload:
  -p string
        ścieżka do katalogu z wygenerowanymi fakturami
  -t    użyj bramki testowej
```

Przykładowe wywołanie dla bramki testowej

```shell
./ksef upload -p katalog-z-plikami-xml -t
```

W wyniku działania komendy program stworzy plik statusu który posłuży do pobierania wizualizacji i/lib UPO
