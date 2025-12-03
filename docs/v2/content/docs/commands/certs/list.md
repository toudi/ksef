---
linkTitle: list
---

# `list`

{{< gateway-flags >}}

Komenda wyświetla listę certyfikatów zapisanych w bazie

Przykładowe wywołanie

```
./ksef certs list
```

Przykładowe wyjście

```
+-------+---------------------+------------------------+------------+---------------+
|  id   |     środowisko      |        funkcja         |    nip     | samopodpisany |
+-------+---------------------+------------------------+------------+---------------+
| aag1h | ksef-test.mf.gov.pl | KsefTokenEncryption    |            | false         |
| ac65d | ksef-test.mf.gov.pl | SymmetricKeyEncryption |            | false         |
| cfbc9 | ksef-test.mf.gov.pl | Authentication         | 1111111111 | false         |
| deabh | ksef-test.mf.gov.pl | Offline                | 1111111111 | false         |
| 08904 | ksef-test.mf.gov.pl | Authentication         | 1111111111 | true          |
| cd0b2 | ksef-test.mf.gov.pl | Authentication         | 3333333333 | false         |
| gf77g | ksef-test.mf.gov.pl | Offline                | 3333333333 | false         |
+-------+---------------------+------------------------+------------+---------------+
```

Wartości w kolumnie `funkcja`

`KsefTokenEncryption`
: Certyfikat MF używany do podpisywania tokenu KSeF. W zasadzie nie powinien być używany, ponieważ tokeny KSeF zostaną wycofane z użycia

`SymmetricKeyEncryption`
: Certyfikat MF używany do zaszyfrowania klucza szyfrującego faktury przy wysyłce

`Authentication`
: Certyfikat użytkownika (zależny od numeru NIP) używany do automatycznej autoryzacji oraz podpisywania plików wyzwania

`Offline`
: Certyfikat użytkownika (zależny od numeru NIP) używany do podpisywania faktur wystawionych w trybie offline
