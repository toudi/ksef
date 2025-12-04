---
linkTitle: dump
---

# `dump`

{{< gateway-flags >}}

Komenda zrzuca klucze zapisane w pliku systemowym do zaszyfrowanego keyringu plikowego

```
./ksef -v -t -n 1111111111 keyring dump --keyring.file.password-file password.txt --keyring.file.path keyring-plikowy
```

Zwróć uwagę, że wymagane jest podanie numeru NIP podmiotu. Komenda zrzuci tylko te wpisy które należą do programu (tj. dotyczą KSeF)
