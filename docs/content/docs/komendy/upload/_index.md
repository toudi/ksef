---
title: Upload (przesyłka faktur)
---

KSeF przewiduje dwa tryby wysyłki faktur - wsadowy (batch) oraz interaktywny. W przypadku sesji interaktywnej wszystkie faktury pakujemy do archiwum, generujemy plik metadanych który następnie podpisujemy a na samym końcu wysyłamy zaszyfrowane archiwum do ministerstwa. Minusem (lub plusem - w zależności jak na to patrzeć) jest to, że w przypadku błędu walidacji jednej faktury odrzucona zostaje cała paczka.

W sesji interaktywnej posługujemy się tokenem (choć czytałem też że zamiast tego ministerstwo planuje wprowadzić indywidualne certyfikaty). W każdym razie wysyłamy faktury jedna po drugiej i odrzucenie którejkolwiek z faktur nie powoduje odrzucenia całej paczki. Jak wspomniałem wyżej może być to paradoksalnie niepożądana sytuacja - jeśli masz do wyeksportowania 1000 faktur to teoretycznie mógłbyś chcieć wysłać je wszystkie w jednej sesji. Z drugiej jednak strony jeśli masz 1000 faktur to śmiem wątpić że używasz mojego programu :-)
