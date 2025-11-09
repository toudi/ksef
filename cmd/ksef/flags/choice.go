package flags

import "fmt"

// https://github.com/spf13/pflag/issues/236#issuecomment-3098646174
// choiceValue implements the [pflag.Value] interface.
type choiceValue struct {
	value    string
	validate func(string) error
}

// Set sets the value of the choice.
func (f *choiceValue) Set(s string) error {
	err := f.validate(s)
	if err != nil {
		return err
	}

	f.value = s
	return nil
}

// Type returns the type of the choice, which must be "string" for [pflag.FlagSet.GetString].
func (f *choiceValue) Type() string { return "string" }

// String returns the current value of the choice.
func (f *choiceValue) String() string { return f.value }

// StringChoice returns a [choiceValue] that validates the value against a set
// of choices. Only the last value will be used if multiple values are set.
func StringChoice(choices []string) *choiceValue {
	return &choiceValue{
		validate: func(s string) error {
			for _, choice := range choices {
				if s == choice {
					return nil
				}
			}
			return fmt.Errorf("must be one of %v", choices)
		},
	}
}
