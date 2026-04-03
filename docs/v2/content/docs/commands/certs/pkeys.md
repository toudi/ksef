---
linkTitle: pkeys
title: Zarządzanie kluczami prywatnymi
---

# `pkeys`

Komenda służy do zarządzania szyfrowaniem kluczy prywatnych w bazie certyfikatów. Dostępne podkomendy umożliwiają zaszyfrowanie lub odszyfrowanie wszystkich kluczy prywatnych przechowywanych w systemie.

## Szyfrowanie kluczy prywatnych

Aby zaszyfrować wszystkie klucze prywatne certyfikatów (zapisując hasło do keyringu):

```
./ksef certs pkeys encrypt
```

W wyniku wykonania komendy:
- Program przeiteruje przez wszystkie certyfikaty z pliku `certificates.yaml`
- Dla każdego certyfikatu program sprawdzi klucz prywatny powiązany z takim certyfikatem. Następnie:
  - Jeśli klucz jest już zaszyfrowany - zostanie on pominięty
  - W przeciwnym wypadku:
    - Klucz zostanie zaszyfrowany
    - Hasło do szyfrowania zostanie zapisane w systemowym keyringu

## Odszyfrowanie kluczy prywatnych

Aby odszyfrować wszystkie zaszyfrowane klucze prywatne (używając kluczy z keyringu) należy użyć komendy `decrypt`:

```
./ksef certs pkeys decrypt
```

W wyniku wykonania komendy:
- Program przeiteruje przez wszystkie certyfikaty z pliku `certificates.yaml`
- Dla każdego certyfikatu program sprawdzi klucz prywatny powiązany z takim certyfikatem. Następnie:
  - Jeśli klucz jest już odszyfrowany - zostanie on pominięty
  - W przeciwnym wypadku:
    - Klucz zostanie odczytany z keyringu
    - Klucz prywatny zostanie odszyfrowany
