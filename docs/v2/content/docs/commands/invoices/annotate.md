---
linkTitle: annotate
weight: 6
---

# `annotate`

Zarządzanie adnotacjami faktur zakupowych – opisywanie pozycji faktur flagami i
komentarzami

---

## Cel i powód istnienia

Podczas pracy z fakturami i przekazywaniem ich do księgowości napotkałem na następujące sytuacje:

- Faktura której pozycja jest niezwykle enigmatyczna. Przykładowo - zakup stacji dokującej gdzie opisem jest jej numer modelu. Dla księgowych jest to spory problem, bo nie wiedzą co to w istocie jest
- Faktura której nie chcę wrzucać w koszty
- Faktura gdzie dozwolone jest odliczenie tylko 50% VAT (akurat w takim przypadku podeujrzewam, że księgowi wiedzą o wiele lepiej ode mnie jak rozpoznawać takie przypadki)

Dotychczas (tj. przed wprowadzeniem KSeF) do adnotacji wystarczało zapisanie małego komentarza na wydruku. Jednak kiedy faktury są elektroniczne a KSeF w ogóle nie przewiduje takiej funkcjonalności - postanowiłem dorobić taką funkcjonalność do programu.

Adnotacje są przechowywane w bazie faktur i mogą być lokalne (dotyczą
konkretnej faktury) albo globalne (ustawione dla danego sprzedawcy – stosowane
automatycznie do wszystkich faktur od tego kontrahenta).

---

## Związek z generowaniem JPK_V7M

Podczas generowania JPK program odczytuje rejestr faktur, a dla każdej pozycji
zakupu sprawdza, czy istnieje dla niej adnotacja. W zależności od flagi:

| Flaga               | Wpływ na JPK                                                                                                                                                 |
|---------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `exclude`           | Pozycja jest **pomijana** – nie trafia do JPK. Dotyczy np. zakupów prywatnych.                                                                              |
| `vat-50-percent`    | Do JPK trafia tylko **połowa** VAT-u. Dotyczy zakupów pojazdów, paliwa itp.                                        |
| `fixed-asset`       | Pozycja jest klasyfikowana jako **środek trwały**. Kwoty bazowe i VAT trafiają do pól **K_40/K_41** zamiast K_42/K_43.                                     |
| `comment`           | Nie wpływa na wartości JPK – komentarz jest widoczny przy wydruku **accountant-notes** (raport dla księgowego).                                             |

---

## Podstawowe użycie

Wszystkie komendy z grupy `annotate` oczekują ścieżki do pliku XML faktury
zakupowej (i - co mam nadzieję jest logiczne - **tylko** zakupowej. W przypadku próby wywołania komendy `annotate` na fakturze wystawionej program nie pozwoli na to i wyświetli komunikat błędu):

```bash
./ksef invoices annotate <komenda> <faktura.xml> [flagi]
```

W każdej komendzie możesz wskazać pozycje faktury za pomocą flagi
`--items` (lista numerów pozycji – od 1) lub znaku `*` (wszystkie pozycje).

---

## `comment` – dodawanie komentarza tekstowego

```bash
./ksef invoices annotate comment \
  --items 1 \
  --text "Myszka do laptopa" \
  data/.../otrzymane/0001-moja-firma-fv-xyz.xml
```

{{< callout type="info" >}}
Flaga `--text` jest wymagana.
{{< /callout >}}

---

## `exclude` – wyłączenie pozycji z JPK

Oznacza pozycje jako wyłączone z raportu JPK (zakup prywatny, niepodlegający
odliczeniu). 

```bash
./ksef invoices annotate exclude \
  --items 2 \
  data/.../otrzymane/0001-moja-firma-fv-xyz.xml
```

---

## `vat-50` – ograniczenie odliczenia VAT do 50%

Przydatne przy zakupie paliwa do pojazdów lub leasingu tychże, gdzie ustawa ogranicza odliczenie
VAT do 50%. 

```bash
./ksef invoices annotate vat-50 \
  --items 1 \
  data/.../otrzymane/0001-moja-firma-fv-xyz.xml
```

---

## `fixed-asset` – oznaczenie środka trwałego

Pozycja oznaczona jako środek trwały trafia do odrębnych pól deklaracji VAT
(K_40 – podstawa, K_41 – VAT zamiast domyślnych K_42/K_43).

```bash
./ksef invoices annotate fixed-asset \
  --items 1 \
  data/.../otrzymane/0001-moja-firma-fv-xyz.xml
```

---

## `list` – podgląd adnotacji

Wyświetla tabelę z pozycjami faktury oraz ich adnotacjami.

```bash
./ksef invoices annotate list data/.../otrzymane/0001-moja-firma-fv-xyz.xml
```

Przykładowe wyjście:

```
+--------+-----------------------------+------------+--------------------------------+
| wiersz |            nazwa            | stawka VAT |             uwagi              |
+--------+-----------------------------+------------+--------------------------------+
| 1      | USB-C Multifunction Adapter | 23         | Stacja dokująca USB do laptopa |
+--------+-----------------------------+------------+--------------------------------+
```

---

## `clear` – usunięcie wszystkich adnotacji

Usuwa wszystkie adnotacje z danej faktury:

```bash
./ksef invoices annotate clear data/.../otrzymane/0001-moja-firma-fv-xyz.xml
```

---

## Adnotacje globalne i lokalne

Domyślnie adnotacje są zapisywane lokalnie – w rejestrze miesięcznym (`registry.yaml`). Można
je jednak zapisać globalnie, w ustawieniach podmiotu, za pomocą flagi `--global`
(`-g`). Adnotacje globalne są wiązane z NIP-em sprzedawcy i automatycznie
stosowane do wszystkich faktur od tego kontrahenta. Funkcję tę przewidziałem na wypadek leasingów bądź np. tankowania u tego samego sprzedawcy aby nie trzeba było co miesiąc powielać identycznych flag.

Przykład – globalne wyłączenie pozycji „Kawa" od kontrahenta o NIP 1111111111:

```bash
./ksef invoices annotate exclude \
  --items 3 \
  --global \
  data/.../otrzymane/0001-moja-firma-fv-xyz.xml
```

Od tej pory każda faktura od tego sprzedawcy, która zawiera pozycję pasującą
do hasha (nazwa, indeks, GTIN lub PKWiU), automatycznie otrzyma adnotację
`exclude`.
