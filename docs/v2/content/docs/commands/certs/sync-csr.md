---
linkTitle: sync-csr
---

# `sync-csr`

{{< gateway-flags >}}

Komenda wysyła oczekujące wnioski do KSeF i/lub pobiera listę gotowych certyfikatów

Przykładowe wywołanie

```
./ksef certs push-csr
```

{{< callout type="important" >}}
Proces pozyskiwania certyfikatów jest **ASYNCHRONICZNY**. W praktyce oznacza to, że komendy `sync-csr` musimy użyć co najmniej dwukrotnie. Pierwszy raz - aby wysłać wnioski do KSeF oraz drugi raz - aby pobrać gotowe certyfikaty. Opisuję to bardziej szczegółowo w rozdziale [Autoryzacja](/docs/authentication)
{{< /callout >}}

{{< callout type="warning" >}}
Warunkiem **koniecznym** prawidłowego działania komendy jest wcześniejsze zalogowanie się do systemu i uzyskanie tokenu sesyjnego. Opisuję to szczegółowo w artykule [Autoryzacja](/docs/authentication)
{{< /callout >}}
