package monthly

import (
	"fmt"
	"strings"

	"github.com/mozillazg/go-slugify"
)

func (i Invoice) ID() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%04d-", i.OrdNum))
	if i.Type == InvoiceTypeReceived {
		builder.WriteString(slugify.Slugify(i.Issuer.Name))
		builder.WriteString("-")
	}
	builder.WriteString(slugify.Slugify(i.RefNo))
	return builder.String()
}
