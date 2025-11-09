package certificates

import (
	"fmt"
	"ksef/internal/certsdb"
	"ksef/internal/config"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "wyświetla listę dostępnych certyfikatów",
	RunE:  listCerts,
}

func init() {
	CertificatesCommand.AddCommand(listCommand)
}

func listCerts(cmd *cobra.Command, _ []string) error {
	env := config.GetGateway(viper.GetViper())
	certDB, err := certsdb.OpenOrCreate(env)
	if err != nil {
		return err
	}

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "id"},
			{Align: simpletable.AlignCenter, Text: "środowisko"},
			{Align: simpletable.AlignCenter, Text: "funkcja"},
			{Align: simpletable.AlignCenter, Text: "nip"},
			{Align: simpletable.AlignCenter, Text: "samopodpisany"},
		},
	}

	for _, cert := range certDB.Certs() {
		table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
			{
				Text: cert.UID,
			},
			{
				Text: string(cert.Environment),
			},
			{
				Text: cert.UsageAsString(),
			},
			{
				Text: cert.NIP,
			},
			{
				Text: fmt.Sprintf("%t", cert.SelfSigned),
			},
		})
	}

	table.SetStyle(simpletable.StyleDefault)
	table.Println()
	return nil
}
