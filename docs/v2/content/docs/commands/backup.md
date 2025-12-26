---
linkTitle: backup
---

# `backup`

Komenda tworzy kopię zapasową plików programu:

- config.yaml
- certificates
- data

Dostępne flagi:

| Flaga               | Opis                                                      |
| ------------------- | --------------------------------------------------------- |
| `-d` / `--add-date` | Dodaj bieżącą datę do nazwy pliku archiwum                |
| `-o` / `--output`   | Zmień nazwę pliku archiwum (domyślnie: `ksef-backup.zip`) |
| `--invoice.pdf`     | Archiwizuj pliki PDF faktur                               |
| `--upo`             | Archiwizuj pliki UPO (domyślnie: tak)                     |
| `--upo.pdf`         | Archiwizuj pliki PDF UPO                                  |
| `-n` / `--nip`      | Archiwizuj dane tylko dla wskazanego numeru NIP           |

## Szyfrowanie archiwum

Aby zwiększyć bezpieczeństwo archiwum (w szczególności jeśli planujesz wysłać je do publicznej chmury) możesz użyć hasła szyfrowania, które zostanie zapisane w keyringu.
Oczywiście aby powyższa opcja zadziałała musisz skonfigurować keyring. Polecam zapisać konfigurację keyringu do pliku konfiguracyjnego aby nie wydłużać linii komend.

### Ustawienie hasła

```
./ksef backup set-password
```

Hasło dla wybranego podmiotu

```
./ksef backup set-password -n 1111111111
```

### Usunięcie hasła

Globalne hasło

```
./ksef backup set-password -d
```

Hasło dla wybranego podmiotu

```
./ksef backup set-password -d -n 1111111111
```
