---
linkTitle: set-profile
---

# `set-profile`

Aby ułatwić pracę z wieloma certyfikatami wystawionymi w obrębie tego samego podmiotu (np. certyfikaty wystawiane per pracownik) wprowadziłem pojęcie "profilu" certyfikatu. Dlaczego operuję na grupie certyfikatów? Ponieważ w przypadku wystawiania faktur w trybie offline musimy użyć osobnego certyfikatu, który i tak należy do grupy

Można go podać przy imporcie lub dodać do już istniejących certyfikatów. Oto przykład:

## Dodawanie profilu do zaimportowanych certyfikatów

1. Wybór certyfikatów

   ```
   > ./ksef certs list
   time=2025-12-24T15:35:13.739+01:00 level=INFO msg="start programu"
   time=2025-12-24T15:35:13.739+01:00 level=INFO msg="wybrane środowisko" env=ksef.mf.gov.pl
   +-------+---------------------+------------------------+------------+--------+---------------+
   |  id   |     środowisko      |        funkcja         |    nip     | profil | samopodpisany |
   +-------+---------------------+------------------------+------------+--------+---------------+
   | aag1h | ksef-test.mf.gov.pl | KsefTokenEncryption    |            |        | false         |
   | ac65d | ksef-test.mf.gov.pl | SymmetricKeyEncryption |            |        | false         |
   | cfbc9 | ksef-test.mf.gov.pl | Authentication         | 1111111111 |        | false         | <===
   | deabh | ksef-test.mf.gov.pl | Offline                | 1111111111 |        | false         | <===
   ```

2. Chcę oflagować dwa certyfikaty tym samym profilem. Wywołuję komendę:

   ```
   > ./ksef certs set-profile --cert-ids deabh,cfbc9 --profile toudi
   time=2025-12-24T15:36:42.354+01:00 level=INFO msg="start programu"
   time=2025-12-24T15:36:42.354+01:00 level=INFO msg="wybrane środowisko" env=ksef.mf.gov.pl
   ```

3. Potwierdzenie:

   ```
   > ./ksef certs list
   time=2025-12-24T15:36:48.192+01:00 level=INFO msg="start programu"
   time=2025-12-24T15:36:48.192+01:00 level=INFO msg="wybrane środowisko" env=ksef.mf.gov.pl
   +-------+---------------------+------------------------+------------+--------+---------------+
   |  id   |     środowisko      |        funkcja         |    nip     | profil | samopodpisany |
   +-------+---------------------+------------------------+------------+--------+---------------+
   | aag1h | ksef-test.mf.gov.pl | KsefTokenEncryption    |            |        | false         |
   | ac65d | ksef-test.mf.gov.pl | SymmetricKeyEncryption |            |        | false         |
   | cfbc9 | ksef-test.mf.gov.pl | Authentication         | 1111111111 | toudi  | false         |
   | deabh | ksef-test.mf.gov.pl | Offline                | 1111111111 | toudi  | false         |
   ```

## Ustawianie nazwy profilu do nowo importowanych certyfikatów

```
./ksef certs import [ ... ] --profile toudi
```

## Użycie profilu

Kiedykolwiek masz do czynienia z autoryzacją (lub importem faktur w trybie offline) możesz wskazać preferowany certyfikat:

```
--cert-profile
```

Przykładowo:

```
./ksef auth sign --cert-profile toudi
./ksef invoices sync -n 1111111111 --cert-profile toudi [ ... ]
```
