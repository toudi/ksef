package certificates

import (
	"fmt"
	"ksef/internal/resources"
	"ksef/internal/runtime"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	certSerial = "serial"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "wyświetla listę dostępnych certyfikatów",
	RunE:  listCerts,
}

func init() {
	flagSet := listCommand.Flags()
	flagSet.String(certSerial, "", "filtruj wg. numeru seryjnego certyfikatu")
	CertificatesCommand.AddCommand(listCommand)
}

func listCerts(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()
	ctx := cmd.Context()

	certsDB, err := resources.CertsDB(vip)
	if err != nil {
		return err
	}
	defer resources.Close(ctx)

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "id"},
			{Align: simpletable.AlignCenter, Text: "środowisko"},
			{Align: simpletable.AlignCenter, Text: "funkcja"},
			{Align: simpletable.AlignCenter, Text: "nip"},
			{Align: simpletable.AlignCenter, Text: "profil"},
			{Align: simpletable.AlignCenter, Text: "samopodpisany"},
			{Align: simpletable.AlignCenter, Text: "podmiot"},
		},
	}

	certSerial := vip.GetString(certSerial)

	certs, err := certsDB.Certificates(ctx)
	if err != nil {
		return err
	}
	selectedGateway, err := runtime.Gateway(vip)
	if err != nil {
		return err
	}

	for cert := range certs {
		if cert.Gateway != selectedGateway {
			continue
		}

		if certSerial != "" && cert.SerialNumber != certSerial {
			continue
		}

		table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
			{
				Text: cert.UID,
			},
			{
				Text: string(cert.Gateway),
			},
			{
				Text: cert.UsageAsString(),
			},
			{
				Text: strings.Join(cert.NIP, ", "),
			},
			{
				Text: cert.ProfileName,
			},
			{
				Text: fmt.Sprintf("%t", cert.SelfSigned),
			},
			{
				Text: lo.FromPtrOr(cert.CN, "-"),
			},
		})
	}

	table.SetStyle(simpletable.StyleDefault)
	table.Println()
	return nil
}
