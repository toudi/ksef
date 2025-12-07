package sendqueue

const (
	registryName = "send-queue.yaml"
)

type Invoice struct {
	Month    int
	FileName string
}

type SendQueue struct {
	invoices []*Invoice
}
