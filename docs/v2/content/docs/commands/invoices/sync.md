---
linkTitle: sync
---

# `sync`

Synchronizuj listę faktur z KSeF.

Ta komenda wykonuje następujące czynności:

1. Wysyła wszystkie niewysłane faktury do KSeF
1. Czeka na zakończenie przetwarzania i (jeśli to możliwe) pobiera UPO
1. Pobiera faktury które zostały przygotowane **dla** podmiotu. Innymi słowy - faktury otrzymane, faktury gdzie podmiot występuje jako płatnik etc etc.
