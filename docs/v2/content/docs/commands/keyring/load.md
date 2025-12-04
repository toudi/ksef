---
linkTitle: load
---

# `load`

{{< gateway-flags >}}

Komenda odczytuje keyring plikowy i przenosi klucze do keyringu systemowego (jest to odwrotność komendy `dump`)

```
./ksef -v -t -n 1111111111 keyring load --keyring.file.password-file password.txt --keyring.file.path keyring-plikowy
```

Zwróć uwagę, że wymagane jest podanie numeru NIP podmiotu.
