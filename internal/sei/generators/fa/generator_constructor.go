package fa

func WithHooks(hooks GeneratorHooks) func(fa *FAGenerator) {
	return func(fa *FAGenerator) {
		fa.hooks = hooks
	}
}

func WithElementOrdering(ordering ElementOrdering) func(fa *FAGenerator) {
	return func(fa *FAGenerator) {
		fa.elementOrdering = ordering
	}
}

func WithCommonData(commonData map[string]string) func(fa *FAGenerator) {
	return func(fa *FAGenerator) {
		fa.commonData = commonData
	}
}
