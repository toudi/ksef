package certificates

import (
	"bufio"
	"errors"
	"iter"
	"ksef/cmd/ksef/commands/client"
	"ksef/core/certsdb"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var internalIdsCommand = &cobra.Command{
	Use:   "internal-ids",
	Short: "zarządza identyfikatorami wewnętrznymi",
	RunE:  cobra.NoArgs,
}

var addInternalIdsCommand = &cobra.Command{
	Use:     "add",
	Short:   "Dodaje jeden lub więcej identyfikatorów wewnętrznych do istniejącego certyfikatu",
	RunE:    addInternalIdsToCertificate,
	PreRunE: prepareCertsDB,
}

var removeInternalIdsCommand = &cobra.Command{
	Use:     "remove",
	Short:   "Usuwa jeden lub więcej identyfikatorów wewnętrznych z istniejącego certyfikatu",
	RunE:    removeInternalIdsFromCertificate,
	PreRunE: prepareCertsDB,
}

const (
	flagNameCertificateId = "cert-id"
	flagNameIds           = "ids"
)

var certsDB *certsdb.CertsDB

func init() {
	internalIdsFlags := internalIdsCommand.PersistentFlags()
	internalIdsFlags.StringSlice(flagNameIds, nil, "Lista identyfikatorów do dodania (lub ścieżka do pliku która go zawiera)")
	internalIdsFlags.String(flagNameCertificateId, "", "Identyfikator certyfikatu (wewnętrzny, numer seryjny lub nazwa profilu)")
	internalIdsCommand.MarkFlagRequired(flagNameCertificateId)
	internalIdsCommand.MarkFlagRequired(flagNameIds)
	internalIdsCommand.AddCommand(addInternalIdsCommand, removeInternalIdsCommand)
	CertificatesCommand.AddCommand(internalIdsCommand)
}

func prepareCertsDB(cmd *cobra.Command, _ []string) (err error) {
	vip := viper.GetViper()

	certsDB, err = client.CertsDB(vip)
	if err != nil {
		return err
	}

	return nil
}

func internalIdsIterator(flagSet *pflag.FlagSet) iter.Seq[string] {
	idsSlice, _ := flagSet.GetStringSlice(flagNameIds)
	return func(yield func(string) bool) {
		fileReader, err := os.Open(idsSlice[0])
		if err != nil && errors.Is(err, os.ErrNotExist) {
			// this is just a regular input, not file-based
			for _, id := range idsSlice {
				if !yield(id) {
					break
				}
			}
		} else {
			// this is a file-based input therefore let's read it line by line
			defer fileReader.Close()
			reader := bufio.NewScanner(fileReader)

			for reader.Scan() {
				line := strings.TrimSpace(reader.Text())
				if !yield(line) {
					break
				}
			}
		}
	}
}

func addInternalIdsToCertificate(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := cmd.Context()
	certId, _ := flags.GetString(flagNameCertificateId)
	idsIterator := internalIdsIterator(flags)

	if err := certsDB.AddInternalIDsToCertificate(ctx, certId, idsIterator); err != nil {
		return err
	}

	return certsDB.Close(ctx)
}

func removeInternalIdsFromCertificate(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()
	ctx := cmd.Context()
	certId, _ := flags.GetString(flagNameCertificateId)
	idsIterator := internalIdsIterator(flags)

	if err := certsDB.RemoveInternalIdsFromCertificate(ctx, certId, idsIterator); err != nil {
		return err
	}

	return certsDB.Close(ctx)
}
