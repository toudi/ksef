kompilacja programu:

go build ksef/cmd/main.go

w celu rekompilacji schematów:

go run parse_schemas.go

wówczas program sparsuje schematy z katalogu "schemas" i wygeneruje odpowiednie struktury w katalogu "generators"

Jeśli zastanawiasz się po kiego grzyba jest ten generowany kod spieszę odpowiedzieć, że niestety ministerstwo używa typu sequence a on wymusza aby elementy w drzewie występowały w określonej kolejności (sic!) miałem więc do wyboru albo zaimplementować struktury w ten sposób, żeby ręcznie wklepać je do kodu w golang albo zaimplementować je w sposób ogólniejszy aby to użytkownik programu wypełniał te pola / atrybuty które wie, że potrzebuje a program na podstawie sparsowania schemy posortuje atrybuty według kolejności i XML przejdzie walidację.

## zapis tokenu

przesyłanie faktur poprzez sesję wsadową jest dość uciążliwe ponieważ wymaga ono każdorazowego podpisu paczki faktur, co w przypadku korzystania z bramki e-obywatel / profil zaufany powoduje generowanie kilku wiadomości SMS. Aby tego uniknąć KSeF przewidział sesję interaktywną. Po zalogowaniu się do aplikacji można tam wygenerować token. Następnie należy uzyć komendy `save-token` aby zapisać go do rejestru kluczy systemu operacyjnego.

```bash
  -nip string
    	numer NIP podatnika
  -t	użyj bramki testowej
  -token string
    	token wygenerowany na środowisku KSeF
```

Przykładowe wywołanie

```bash
./ksef save-token -t -nip 1111111111 -token AAABBBCCC....
```

Od tej pory, podczas wysyłki faktur, program rozpozna wystawcę faktur (tj. jego numer NIP) i pobierze z rejestru kluczy odpowiedni token

## generowanie faktur:

./ksef generate -d ';' -f faktury.csv -o katalog-docelowy [-t]

(parametr -t uzywa klucza publicznego bramki testowej do generowania metadanych)

## metadane (tylko tryb wsadowy)

```bash
  -p string
    	ścieżka do wygenerowanych plików
  -t	użyj bramki testowej
```

Jeśli nie chcesz używać tokenu i zamiast tego wolisz przesyłać faktury w sesji wsadowej, musisz najpierw wygenerować plik metadanych a następnie podpisać go.

Przykładowe wywołanie:

```bash
./ksef metadata -p katalog-z-plikami-faktur [-t]
```

w katalogu docelowym program stworzy pliki:

* metadata.xml [surowy plik metadanych który należy podpisać. Podpisanego pliku użyjesz w kolejnym kroku (`wysyłka faktur`)]
* metadata.zip [surowy plik archiwum, nie jest on wysyłany na serwer]
* metadata.zip.aes [plik archiwum zaszyfrowany odpowiednim kluczem ministerstwa, zależnym od wybranego trybu (testowy / produkcja) - to ten plik jest przesyłany]

### podpisywanie pliku metadanych

Aby podpisać plik metadanych użyj trybu "Osadzonego". Możesz użyć do tego celu karty kryptograficznej lub aplikacji m-obywatel: https://moj.gov.pl/nforms/signer/upload?xFormsAppName=SIGNER&xadesPdf=true. Po podpisaniu dokumentu bramka zwróci plik xml z doklejoną sekcją podpisu (`Signature`). Należy ten plik zapisać i przejść do kolejnego kroku (`wysyłka faktur`):

## wysylka faktur

```bash
  -i	użyj sesji interaktywnej
  -p string
    	ścieżka do katalogu z wygenerowanymi fakturami
  -sj
    	użyj formatu JSON do zapisu pliku statusu (domyślnie YAML)
  -t	użyj bramki testowej
```

Przewidziane zostały dwa tryby wysyłki faktur.

### tryb wsadowy
Aby skorzystać z trybu wsadowego upewnij się, że podisałeś plik metadanych (patrz sekcja `metadane`). Następnie wywołaj ksef w poniższy sposób:

```bash
./ksef upload -p podpisany-metadata.xml [-t]
```

### tryb interaktywny

Aby skorzystać z trybu interaktywnego należy uprzednio wygenerować token na stronie aplikacji KSeF oraz zapisać go do systemowego repozytorium kluczy (patrz sekcja `zapis tokenu`). Następnie wywołujemy ksef w następujący sposób:

```bash
./ksef upload -i -p katalog-z-plikami-xml [-t]
```

Niezależnie od wybranego trybu wysyłki, program utworzy plik `status.ref` który posłuży do sprawdzenia statusu i pobrania UPO (w przypadku pozytywnego przetworzenia faktur)

## pobieranie upo

```bash
Usage of status:
  -o string
    	ścieżka do zapisu UPO (domyślnie katalog pliku statusu + {nrRef}.pdf)
  -p string
    	ścieżka do pliku statusu
  -xml
    	zapis UPO jako plik XML
```

Przykładowe wywołanie:

```bash
./ksef status -p sciezka-do-pliku-status.ref
```

jesli status przetworzenia faktur bedzie poprawny, program pobierze upo i wygeneruje plik:

- {nrRef}.xml (jesli zostanie wybrany parametr -xml)
- {nrRef}.pdf (domyślnie)

UPO jest gotowe do druku.
