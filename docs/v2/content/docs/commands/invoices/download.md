---
linkTitle: download
weight: 1
---

# `download`

Komenda służy do pobierania faktur z systemu KSeF i zapisywania ich w katalogu danych (`data`)

{{< gateway-flags >}}

# Najważniejsze flagi

| Flagi | Opis | Domyślna wartość |
|-------|------|------------------|
| `-n, --nip` | Numer NIP podmiotu | - |
| `-w, --workers` | Ilość workerów (domyślnie 0; wartość > 0 oznacza tryb równoległy) | `0` |
| `--invoice.pdf` | Generuj PDF dla pobranych faktur | `false` |
| `--start-date` | Data początkowa zakresu | pierwszy dzień miesiąca |
| `--end-date` | Data końcowa zakresu | puste (brak limitu) |
| `--income` | Pobieranie faktur przychodowych | - |
| `--cost` | Pobieranie faktur kosztowych | - |
| `--payer` | Pobieranie faktur płatnika | - |
| `--authorized` | Pobieranie faktur strony upoważnionej | - |
| `--use-export-mode` | Używaj eksportu faktur do pobierania | `false` |
| `--use-smart-mode` | Inteligentny tryb pobierania | `false` |


## Podstawowe użycie

```
./ksef invoices download -n 1112223344
```

Jeśli nie podasz żadnego z parametrów (oprócz numeru NIP), komenda pobierze faktury tylko w formacie XML i zapisze je do katalogu danych. Program pobierze faktury od daty ostatniej synchronizacji lub od początku miesiąca (jeśli nigdy wcześniej nie pobierałeś faktur). Dodatkowo, program spróbuje pobrać faktury dla wszystkich rodzajów strony podmiotu (tj. Subject2 (Zakupowe) + Subject3 (Odbiorca?) + SubjectAuthorized (Strona upoważniona (???))). Aby ograniczyć podmiot do tylko jednego rodzaju skorzystaj z flag `--cost` / `--payer` / `--authorized`

## Tryb współbieżny / równoległy

Ta opcja jest przewidziana dla użytkowników którzy obsługują więcej niż jeden NIP. Użyj parametru `-w` aby wskazać liczbę wątków które będą współbieżnie pobierać faktury.

Przykładowe wywołanie:
```
./ksef invoices download -w 4
# lub
./ksef invoices download --workers 4
```

{{< callout type="error" >}}
Komenda **NIE** jest przewidziana do równoległego pobierania faktur dla jednego NIPu - przy limitach MF byłoby to kompletnie bezsensowne.
{{< /callout >}}

{{< callout type="info" >}}
Podawanie numeru NIP **nie** jest wymagane - program odczytuje bazę certyfikatów (`certificates.yaml`) i na jej podstawie jest w stanie sam określić wszystkie numery NIP.

Dodatkowo: program automatycznie dostosowuje liczbę workerów, jeśli jest ich więcej niż NIPów do przetworzenia

Przykład: Jeśli masz 3 NIP-y i podasz `-w 5`, system użyje tylko 3 workerów.
{{< /callout >}}

## Generowanie PDF

Aby pobrać faktury wraz z plikami PDF w formacie dokumentu:

```
./ksef invoices download --invoice.pdf
```

## Określanie zakresu dat

Możesz określić niestandardowy zakres dat dla pobierania:

```
./ksef invoices download \
  --start-date 2025-01-01 \
  --end-date 2025-01-31
```

{{< callout type="info" >}}
**Domyślna data początkowa**: pierwszy dzień bieżącego miesiąca

Możesz pobrać faktury z dowolnego zakresu dat, np. z całego roku lub z konkretnego tygodnia.
{{< /callout >}}

## Tryb eksportu i inteligentnego pobierania

Ministerstwo Finansów wprowadziło dość drakońskie limity jeśli idzie o pobieranie faktur. Jest to (jeśli dobrze pamiętam) 8 pobrań / minutę i 64 / godzinę. Program oczywiście obsługuje tzw. rate limiting i czeka odpowiednią liczbę sekund, tym niemniej jeśli spodziewasz się sporej liczby faktur, warto wspomóc się trybem inteligentnym lub eksportu.

W trybie eksportu program używa pojedynczego requestu aby skonstruować żądanie eksportu a następnie pobiera paczkę ZIP zawierającą wyeksportowane faktury. Teoretycznie mogłoby się wydawać, że warto byłoby użyć tego trybu domyślnie. Niestety - tryb eksportu jest asynchroniczny co w praktyce oznacza, że możemy stracić więcej czasu na oczekiwanie na archiwum niż zużylibyśmy pobierając faktury jedna po drugiej. 

Użycie trybu eksportu:

```bash
./ksef invoices download --use-export-mode
```

Zamiast tego (tzn jeśli nie wiesz ile masz faktur do pobrania) rozważ tryb inteligentny.

W trybie inteligentnym:
- Program pobiera listę oczekujących faktur (1 request)
- Jeśli liczbę faktur da radę pobrać mieszcząc się dwukrotnie w limicie / minutę - używane jest pobieranie "proste" (czyli po prostu, jedna faktura po drugiej)
- W przeciwnym wypadku - używane jest pobieranie w trybie eksportu.

```
./ksef invoices download --use-smart-mode
```
