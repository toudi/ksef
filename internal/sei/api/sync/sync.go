package sync

import (
	"fmt"
	"ksef/internal/logging"
	registryPkg "ksef/internal/registry"
	"ksef/internal/sei/api/client"
	"ksef/internal/sei/api/session/interactive"
	"path"
	"time"
)

type SyncInvoicesConfig struct {
	Income            bool
	Cost              bool
	Subject3          bool
	SubjectAuthorized bool
	StartDate         time.Time
	DestPath          string
	SubjectTIN        string
	IssuerToken       string
	Token             string
}

const queryInit string = "/api/online/Query/Invoice/Sync"

func SyncInvoices(
	apiClient *client.APIClient,
	params *SyncInvoicesConfig,
	registry *registryPkg.InvoiceRegistry,
) error {
	var err error

	if registry == nil {
		registry = registryPkg.NewRegistry()
		registry.QueryCriteria.Type = "range"
		registry.QueryCriteria.DateFrom = params.StartDate
		registry.Issuer = params.SubjectTIN
		if params.Cost {
			registry.QueryCriteria.SubjectType = "subject2"
		} else if params.Income {
			registry.QueryCriteria.SubjectType = "subject1"
		} else if params.Subject3 {
			registry.QueryCriteria.SubjectType = "subject3"
		} else if params.SubjectAuthorized {
			registry.QueryCriteria.SubjectType = "subjectAuthorized"
		}
	}

	interactiveSession := interactive.InteractiveSessionInit(apiClient)

	if params.Token == "" {
		if params.IssuerToken != "" {
			interactiveSession.SetIssuerToken(params.IssuerToken)
		}

		if err = interactiveSession.Login(registry.Issuer, true); err != nil {
			return fmt.Errorf("unable to login: %v", err)
		}
	}

	httpSession := interactiveSession.HTTPSession()

	if params.Token != "" {
		httpSession.SetHeader("SessionToken", params.Token)
	}

	var queryInitStruct struct {
		Criteria registryPkg.QueryCriteria `json:"queryCriteria"`
	}

	var queryInitResponse struct {
		DocumentCount int                   `json:"numberOfElements"`
		Offset        int                   `json:"pageOffset"`
		Invoices      []registryPkg.Invoice `json:"invoiceHeaderList"`
	}

	queryInitStruct.Criteria = registry.QueryCriteria
	queryInitStruct.Criteria.DateTo = time.Now().Truncate(time.Second)

	var processedInvoices int = 0
	var pageOffset int = 0
	var syncFinished bool = false

	for !syncFinished {
		response, err := httpSession.JSONRequest(
			client.JSONRequestParams{
				Method:   "POST",
				Endpoint: queryInit + fmt.Sprintf("?PageSize=50&PageOffset=%d", pageOffset),
				Payload:  queryInitStruct,
				Response: &queryInitResponse,
				Logger:   logging.DownloadHTTPLogger,
			},
		)
		if err != nil {
			fmt.Printf("response code: %d\n", response.StatusCode)
			return fmt.Errorf("unable to send queryInit: %v", err)
		}
		for _, invoice := range queryInitResponse.Invoices {
			if !registry.Contains(invoice.SEIReferenceNumber) {
				invoice.SubjectFrom.TIN = invoice.SubjectFrom.IssuedBy.TIN
				invoice.SubjectFrom.FullName = invoice.SubjectFrom.Issuer.FullName
				invoice.SubjectTo.TIN = invoice.SubjectTo.IssuedBy.TIN
				invoice.SubjectTo.FullName = invoice.SubjectTo.Issuer.FullName

				registry.Invoices = append(registry.Invoices, invoice)
			} else {
				logging.DownloadLogger.Debug(
					"invoice already in registry; no-op",
					"invoiceRefNo", invoice.SEIReferenceNumber,
				)
			}

			processedInvoices += 1
		}
		syncFinished = processedInvoices >= queryInitResponse.DocumentCount
		pageOffset += 1
	}

	registry.QueryCriteria = queryInitStruct.Criteria
	registry.Environment = apiClient.EnvironmentAlias

	return registry.Save(path.Join(params.DestPath, "registry.yaml"))
}
