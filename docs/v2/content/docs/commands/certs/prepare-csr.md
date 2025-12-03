---
linkTitle: prepare-csr
---

# `prepare-csr`

{{< gateway-flags >}}

Komenda przygotowuje żądanie wystawienia certyfikatów

Przykładowe wywołanie

```
./ksef certs prepare-csr -a -o -n 1111111111
```

Obsługiwane flagi:

| flaga       | flaga (skrót) | znaczenie                                                 |
| ----------- | ------------- | --------------------------------------------------------- |
| `--auth`    | `-a`          | wygeneruj żądanie wystawienia certyfikatu autoryzacyjnego |
| `--offline` | `-o`          | wygeneruj żądanie wystawienia certyfikatu offline         |
| `--nip`     | `-n`          | numer NIP podmiotu                                        |

Przynajmniej jedna z flag (`-a` / `-o`) musi być wybrana. Przy czym czysty pragmatyzm sugerowałby wybranie obu tych flag - w szczególności tego do wystawiania faktur w trybie offline ;-)

{{< callout type="warning" >}}
Warunkiem **koniecznym** prawidłowego działania komendy jest wcześniejsze zalogowanie się do systemu i uzyskanie tokenu sesyjnego. KSeF sam wypełnia dane użytkownika na podstawie kontekstu sesji. Opisuję to szczegółowo w artykule [Autoryzacja](/docs/authentication)
{{< /callout >}}

W wyniku wywołania komenda przygotuje wnioski o wystawienie certyfikatów i zapisze je w pliku `certificates/certificates.yaml`
