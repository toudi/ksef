# Wysyłka interaktywna

W trybie wysyłki interaktywnej posługujemy się tokenem (choć czytałem też że zamiast tego ministerstwo planuje wprowadzić indywidualne certyfikaty). W każdym razie wysyłamy faktury jedna po drugiej i odrzucenie którejkolwiek z faktur nie powoduje odrzucenia całej paczki. Może być to paradoksalnie niepożądana sytuacja - jeśli masz do wyeksportowania 1000 faktur to teoretycznie mógłbyś chcieć wysłać je wszystkie w jednej sesji. Z drugiej jednak strony jeśli masz 1000 faktur to śmiem wątpić że używasz mojego programu :-)

## Token sesji

Upewnij się, że zapisałeś token [save-token](/content/komendy/save-token)

### Jeśli nie możesz użyć komendy `save-token`

::: danger
Zignoruj tą sekcję jeśli udało Ci się zapisać token do systemowego pęku kluczy
:::

Jeśli zapisanie tokenu do systemowego pęku kluczy nie jest to możliwe (np. dysponujesz serwerem w trybie headless który nie ma sesji dbus) przewidziałem przekazanie tokenu przez zmienną środowiskową i/lub otwartym tekstem

Oto przykład (w przykładzie surowy token zapisany jest w pliku token.plaintext)

```shell
more token.plaintext
12345
```

Szyfrowanie pliku:

```shell
gpg -c token.plaintext
```

Weryfikacja, że plik został zaszyfrowany:

```shell
hexdump -C token.plaintext.gpg
00000000  8c 0d 04 09 03 08 2d 40  22 6e b6 65 c0 90 ff d2  |......-@"n.e....|
00000010  4a 01 91 00 f3 52 53 04  1c e5 49 44 1d ab 6c dc  |J....RS...ID..l.|
00000020  db a2 09 5b 49 44 53 47  27 f8 83 d0 da 10 eb 67  |...[IDSG'......g|
00000030  58 88 e8 1e 49 fc bc 29  d0 84 8f fc 7e c3 37 8e  |X...I..)....~.7.|
00000040  47 87 c4 f8 c4 66 79 39  ea 3b c0 bd ff 3e ac 89  |G....fy9.;...>..|
00000050  d1 7d cb 79 20 02 e2 21  c4 ee 5b                 |.}.y ..!..[|
0000005b
```

Użycie komendy `upload` z przekazaniem tokenu w zmiennej środowiskowej

```shell
TOKEN=`gpg -d token.plaintext.gpg` ./ksef upload -t -i -p przyklady -token TOKEN
```

Użycie komendy `upload` z przekazaniem tokenu jawnie (z oczywistych względów jest to niebezpieczne i raczej nie powinno być stosowane)

```shell
./ksef upload -t -i -p przyklady -token AAABBBCCC...
```

## Wysyłka faktur

::: info
Program odczyta numer NIP wystawcy faktur ze źródłowych plików XML a następnie spróbuje odczytać token dla tego NIP'u z pęku kluczy. Dlatego istotne jest aby przed wywołaniem komendy `upload` wywołać najpierw komendę `save-token`
:::

```text
./ksef upload
Usage of upload:
  -i    użyj sesji interaktywnej
  -p string
        ścieżka do katalogu z wygenerowanymi fakturami
  -t    użyj bramki testowej
```

Przykład wywołania:

```shell
./ksef upload -i -p pliki-zrodlowe-xml -t
```
