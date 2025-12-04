# keyring

Tokeny KSeF przechowywane są w keyringu. Domyślnie jest to keyring systemowy (tj. wbudowany w system operacyjny taki jak MacOS, Windows czy Linux), ale w wyjątkowej sytuacji (tj. w przypadku systemu który nie jest w stanie wspierać tego rozwiązania) można użyć keyringu opartego o zaszyfrowany plik. Opiszę tutaj jak to zrobić

Flagi dotyczące keyringu

| Flaga                             | Opis                                                                                    |
| --------------------------------- | --------------------------------------------------------------------------------------- |
| `--keyring.engine`                | silnik keyringu (domyślna wartość to `system`). Dopuszczalne wartości: `system`, `file` |
| `--keyring.file.ask-password`     | pytaj o hasło do keyringu na wejściu standardowym (stdin)                               |
| `--keyring.file.buffered`         | buforuj keyring w pamięci (plik zostanie zapisany raz przy zamykaniu programu)          |
| `--keyring.file.password-env-var` | nazwa zmiennej środowiskowej która zawiera hasło do keyringu                            |
| `--keyring.file.password-file`    | ścieżka do pliku z hasłem keyringu                                                      |
| `--keyring.file.path`             | ścieżka do keyringu opartego o plik                                                     |

Jak nietrudno się domyślić najbardziej newralgicznymi flagami są te które odpowiadają za przekazanie hasła do programu. Poniżej zamieszczam kilka przykładów wywołań

{{< callout type="important" >}}
Hasło musi mieć długość 16 lub 32 znaków i nie jest to moja fanaberia tylko wynika to ze specyfikacji zastosowanego szyfru AES
{{< /callout >}}

## Przykład 1: Użycie hasła zapisanego w pliku

```
./ksef --keyring.engine file --keyring.file.path keyring-plikowy --keyring.file.password-file password.txt --keyring.file.buffered=true [ ... ]
```

- Program użyje keyringu zapisanego w zaszyfrowanym pliku `keyring-plikowy`.
- Jeśli plik nie istnieje - zostanie utworzony.
- Hasło zostanie odczytane z pliku `password.txt` (białe znaki są ignorowane)

{{< callout type="important" >}}
Plik z zapisanym hasłem musi mieć uprawnienia do odczytu i zapisu tylko i wyłącznie przez właściciela (0600). W przeciwnym wypadku program odmówi użycia takiego pliku i zakończy się z kodem błędu
{{< /callout >}}

## Przykład 2: Przekazanie hasła przez użytkownika

Nie sądzę żeby to było praktyczne, ale ..

```
./ksef --keyring.engine file --keyring.file.path keyring-plikowy --keyring.file.ask-password --keyring.file.buffered=true [ ... ]
```

- Program użyje keyringu zapisanego w zaszyfrowanym pliku `keyring-plikowy`.
- Jeśli plik nie istnieje - zostanie utworzony.
- Użytkownik zostanie zapytany o podanie hasła

## Przykład 3: Przekazanie nazwy zmiennej środowiskowej zawierającej hasło

```
ZMIENNA_Z_HASLEM=aaaaaaaaaaaaaaaa ./ksef --keyring-engine file --keyring.file.path keyring-plikowy --keyring.file.password-env-var ZMIENNA_Z_HASLEM --keyring.file.buffered=true [ ... ]
```
