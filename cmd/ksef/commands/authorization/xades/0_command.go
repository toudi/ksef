package xades

import (
	"github.com/spf13/cobra"
)

var XadesCommand = &cobra.Command{
	Use:   "xades",
	Short: "zarządzanie opcjami autoryzacji opartej o podpis XAdES",
}
