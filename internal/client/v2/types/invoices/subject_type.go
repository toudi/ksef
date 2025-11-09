package invoices

type SubjectType string

const (
	SubjectTypeInvalid    SubjectType = ""
	SubjectTypeIssuer     SubjectType = "Subject1"
	SubjectTypeRecipient  SubjectType = "Subject2"
	SubjectTypePayer      SubjectType = "Subject3"
	SubjectTypeAuthorized SubjectType = "SubjectAuthorized"
)
