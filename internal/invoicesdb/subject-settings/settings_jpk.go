package subjectsettings

type JPKFormMeta struct {
	IRSCode    int               `yaml:"irs-code,omitempty"`
	SystemName string            `yaml:"system-name,omitempty"`
	Subject    *SubjectData      `yaml:"subject,omitempty"`
	Defaults   map[string]string `yaml:"defaults,omitempty"`
}

type SubjectData struct {
	SubjectType string            `yaml:"type"`
	Data        map[string]string `yaml:"data"`
}

type JPKSettings struct {
	FormMeta JPKFormMeta `yaml:"form,omitempty"`
}
