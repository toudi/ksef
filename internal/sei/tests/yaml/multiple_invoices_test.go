package yaml

import (
	"bytes"
	"ksef/internal/sei"
	"ksef/internal/sei/tests/recorder"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

const inputYAML = `
common:
  Faktura:
    Podmiot1:
      DaneIdentyfikacyjne:
        NIP: 1112223344
        Nazwa: Wytwórnia programów bezwartościowych
      Adres:
        KodKraju: PL
        AdresL1: ul. Kwietniowa 3/14
        AdresL2: 01-234 Mielno
    Fa:
      Platnosc:
        FormaPlatnosci: 6
        RachunekBankowy:
          - NrRB: 1122334455
            NazwaBanku: mBank
          - NrRB: 2233445566
            NazwaBanku: ING
invoices:
  - Faktura:
      Podmiot2:
        DaneIdentyfikacyjne:
          NIP: 2222222222
          Nazwa: Czesław Fikcyjny sp. j.
        Adres:
          KodKraju: PL
          AdresL1: ul. Czereśniowa 3/4
          AdresL2: 00-111 Chmielowo
      Fa:
        KodWaluty: PLN
        P_1: "2026-01-01"
        P_1M: Mielno
        P_2: FV 00/11/23
        P_6: "2026-01-02"
        RodzajFaktury: VAT
        Adnotacje:
          P_16: 2
          P_17: 2
          P_18: 2
          P_18A: 2
          Zwolnienie.P_19N: 1
          NoweSrodkiTransportu.P_22N: 1
          P_23: 2
          PMarzy.P_PMarzyN: 1
        Platnosc:
          Zaplacono: 1
          DataZaplaty: "2020-01-01"
          FormaPlatnosci: 6
        FaWiersze:
          - FaWiersz:
              item: Wprowadzanie do obrotu środków pochodzących z działalności przestępczej
              units: szt
              quantity: 150
              unit-price-net: 123
              vat-rate: 23
  - Faktura:
      Podmiot2:
        DaneIdentyfikacyjne:
          NIP: 8888888888
          Nazwa: Capone & Moran sp. z o.o.
        Adres:
          KodKraju: PL
          AdresL1: ul. Wiśniowa 3/4
          AdresL2: 11-222 Chlebowo
      Fa:
        KodWaluty: PLN
        P_1: "2026-02-02"
        P_1M: Suwałki
        P_2: FV 11/22/34
        P_6: "2026-02-03"
        RodzajFaktury: VAT
        Adnotacje:
          P_16: 2
          P_17: 2
          P_18: 2
          P_18A: 2
          Zwolnienie.P_19N: 1
          NoweSrodkiTransportu.P_22N: 1
          P_23: 2
          PMarzy.P_PMarzyN: 1
        FaWiersze:
          - FaWiersz:
              item: Zestaw wytrychów
              units: szt
              quantity: 1
              unit-price-net: 123
              vat-rate: 23
          - FaWiersz:
              item: Łom
              units: szt
              quantity: 10
              unit-price-net: "1.23"
              vat-rate: 23
`

const expectedInvoice1 = `
<?xml version="1.0" encoding="utf-8"?>
<Faktura xmlns="http://crd.gov.pl/wzor/2025/06/25/13775/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <Naglowek>
    <KodFormularza kodSystemowy="FA (3)" wersjaSchemy="1-0E">FA</KodFormularza>
    <WariantFormularza>3</WariantFormularza>
    <DataWytworzeniaFa>2026-04-14T16:00:00Z</DataWytworzeniaFa>
    <SystemInfo>WSI Pegasus</SystemInfo>
  </Naglowek>
  <Podmiot1>
    <DaneIdentyfikacyjne>
      <NIP>1112223344</NIP>
      <Nazwa>Wytwórnia programów bezwartościowych</Nazwa>
    </DaneIdentyfikacyjne>
    <Adres>
      <KodKraju>PL</KodKraju>
      <AdresL1>ul. Kwietniowa 3/14</AdresL1>
      <AdresL2>01-234 Mielno</AdresL2>
    </Adres>
  </Podmiot1>
  <Podmiot2>
    <DaneIdentyfikacyjne>
      <NIP>2222222222</NIP>
      <Nazwa>Czesław Fikcyjny sp. j.</Nazwa>
    </DaneIdentyfikacyjne>
    <Adres>
      <KodKraju>PL</KodKraju>
      <AdresL1>ul. Czereśniowa 3/4</AdresL1>
      <AdresL2>00-111 Chmielowo</AdresL2>
    </Adres>
    <JST>2</JST>
    <GV>2</GV>
  </Podmiot2>
  <Fa>
    <KodWaluty>PLN</KodWaluty>
    <P_1>2026-01-01</P_1>
    <P_1M>Mielno</P_1M>
    <P_2>FV 00/11/23</P_2>
    <P_6>2026-01-02</P_6>
    <P_13_1>18450.00</P_13_1>
    <P_14_1>4243.50</P_14_1>
    <P_15>22693.5</P_15>
    <Adnotacje>
      <P_16>2</P_16>
      <P_17>2</P_17>
      <P_18>2</P_18>
      <P_18A>2</P_18A>
      <Zwolnienie>
        <P_19N>1</P_19N>
      </Zwolnienie>
      <NoweSrodkiTransportu>
        <P_22N>1</P_22N>
      </NoweSrodkiTransportu>
      <P_23>2</P_23>
      <PMarzy>
        <P_PMarzyN>1</P_PMarzyN>
      </PMarzy>
    </Adnotacje>
    <RodzajFaktury>VAT</RodzajFaktury>
    <FaWiersz>
      <NrWierszaFa>1</NrWierszaFa>
      <P_7>Wprowadzanie do obrotu środków pochodzących z działalności przestępczej</P_7>
      <P_8A>szt</P_8A>
      <P_8B>150</P_8B>
      <P_9A>123</P_9A>
      <P_11>18450</P_11>
      <P_12>23</P_12>
    </FaWiersz>
    <Platnosc>
      <Zaplacono>1</Zaplacono>
      <DataZaplaty>2020-01-01</DataZaplaty>
      <FormaPlatnosci>6</FormaPlatnosci>
      <RachunekBankowy>
        <NrRB>1122334455</NrRB>
        <NazwaBanku>mBank</NazwaBanku>
      </RachunekBankowy>
      <RachunekBankowy>
        <NrRB>2233445566</NrRB>
        <NazwaBanku>ING</NazwaBanku>
      </RachunekBankowy>
    </Platnosc>
  </Fa>
</Faktura>
`

const expectedInvoice2 = `
<?xml version="1.0" encoding="utf-8"?>
<Faktura xmlns="http://crd.gov.pl/wzor/2025/06/25/13775/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <Naglowek>
    <KodFormularza kodSystemowy="FA (3)" wersjaSchemy="1-0E">FA</KodFormularza>
    <WariantFormularza>3</WariantFormularza>
    <DataWytworzeniaFa>2026-04-14T16:00:01Z</DataWytworzeniaFa>
    <SystemInfo>WSI Pegasus</SystemInfo>
  </Naglowek>
  <Podmiot1>
    <DaneIdentyfikacyjne>
      <NIP>1112223344</NIP>
      <Nazwa>Wytwórnia programów bezwartościowych</Nazwa>
    </DaneIdentyfikacyjne>
    <Adres>
      <KodKraju>PL</KodKraju>
      <AdresL1>ul. Kwietniowa 3/14</AdresL1>
      <AdresL2>01-234 Mielno</AdresL2>
    </Adres>
  </Podmiot1>
  <Podmiot2>
    <DaneIdentyfikacyjne>
      <NIP>8888888888</NIP>
      <Nazwa><![CDATA[Capone & Moran sp. z o.o.]]></Nazwa>
    </DaneIdentyfikacyjne>
    <Adres>
      <KodKraju>PL</KodKraju>
      <AdresL1>ul. Wiśniowa 3/4</AdresL1>
      <AdresL2>11-222 Chlebowo</AdresL2>
    </Adres>
    <JST>2</JST>
    <GV>2</GV>
  </Podmiot2>
  <Fa>
    <KodWaluty>PLN</KodWaluty>
    <P_1>2026-02-02</P_1>
    <P_1M>Suwałki</P_1M>
    <P_2>FV 11/22/34</P_2>
    <P_6>2026-02-03</P_6>
    <P_13_1>135.30</P_13_1>
    <P_14_1>31.12</P_14_1>
    <P_15>166.42</P_15>
    <Adnotacje>
      <P_16>2</P_16>
      <P_17>2</P_17>
      <P_18>2</P_18>
      <P_18A>2</P_18A>
      <Zwolnienie>
        <P_19N>1</P_19N>
      </Zwolnienie>
      <NoweSrodkiTransportu>
        <P_22N>1</P_22N>
      </NoweSrodkiTransportu>
      <P_23>2</P_23>
      <PMarzy>
        <P_PMarzyN>1</P_PMarzyN>
      </PMarzy>
    </Adnotacje>
    <RodzajFaktury>VAT</RodzajFaktury>
    <FaWiersz>
      <NrWierszaFa>1</NrWierszaFa>
      <P_7>Zestaw wytrychów</P_7>
      <P_8A>szt</P_8A>
      <P_8B>1</P_8B>
      <P_9A>123</P_9A>
      <P_11>123</P_11>
      <P_12>23</P_12>
    </FaWiersz>
    <FaWiersz>
      <NrWierszaFa>2</NrWierszaFa>
      <P_7>Łom</P_7>
      <P_8A>szt</P_8A>
      <P_8B>10</P_8B>
      <P_9A>1.23</P_9A>
      <P_11>12.3</P_11>
      <P_12>23</P_12>
    </FaWiersz>
    <Platnosc>
      <FormaPlatnosci>6</FormaPlatnosci>
      <RachunekBankowy>
        <NrRB>1122334455</NrRB>
        <NazwaBanku>mBank</NazwaBanku>
      </RachunekBankowy>
      <RachunekBankowy>
        <NrRB>2233445566</NrRB>
        <NazwaBanku>ING</NazwaBanku>
      </RachunekBankowy>
    </Platnosc>
  </Fa>
</Faktura>
`

func TestParsingMultipleInvoicesFromYAML(t *testing.T) {
	vip := viper.New()
	vip.Set("generator", "fa-3_1.0")
	recorder := recorder.NewRecorder()
	sei, err := sei.SEI_Init(vip, sei.WithInvoiceReadyFunc(recorder.Ready))
	require.NoError(t, err)

	reader := bytes.NewReader([]byte(inputYAML))

	err = sei.ProcessReader(reader, "dummy.yaml")
	require.NoError(t, err)
	require.Len(t, recorder.Invoices, 2)
	require.Equal(t, strings.TrimLeft(expectedInvoice1, " \n"), recorder.XMLInvoices[0])
	require.Equal(t, strings.TrimLeft(expectedInvoice2, " \n"), recorder.XMLInvoices[1])
}
