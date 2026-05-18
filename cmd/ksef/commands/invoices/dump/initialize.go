package dump

import (
	"fmt"
	"ksef/internal/invoicesdb/annotations"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	dumpMonth          time.Time
	dumpRegistry       *monthlyregistry.Registry
	dumpAnnotationsMgr *annotations.Annotations
	dumpAnnotationCfg  AnnotationConfig
)

// initializeDump validates month arguments, parses the target month,
// opens the monthly registry, and creates the annotations manager.
// Results are stored in package-level variables for use by RunE handlers.
func initializeDump(cmd *cobra.Command, args []string) error {
	if len(args) > 0 && len(args) < 2 {
		return fmt.Errorf("użycie: dump [rok] [miesiąc]")
	}

	vip := viper.GetViper()

	var month time.Time = time.Now().Local().AddDate(0, -1, 0)
	if len(args) > 0 {
		yearNum, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		monthNum, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		month = time.Date(yearNum, time.Month(monthNum), 1, 0, 0, 0, 0, time.Local)
	}

	var err error
	dumpRegistry, err = monthlyregistry.OpenForMonth(vip, month)
	if err != nil {
		return err
	}

	dumpMonth = month

	dumpAnnotationsMgr, err = annotations.Manager(
		vip,
		annotations.WithMonthlyRegistry(dumpRegistry),
	)
	if err != nil {
		return err
	}

	dumpAnnotationCfg = GetConfig(vip)

	return nil
}
