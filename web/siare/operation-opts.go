package siare

type SSApiCallOpts struct {
	NotifyOperationResult  bool
	EnqueueNotification    bool
	ClearCache             bool
	CacheNamespacesToClear []string
	ResourceName           string
}

func NewSSApiCallOpts(resourceName string) *SSApiCallOpts {
	return &SSApiCallOpts{
		NotifyOperationResult:  true,
		EnqueueNotification:    true,
		ClearCache:             true,
		CacheNamespacesToClear: []string{resourceName},
		ResourceName:           resourceName,
	}
}

func InvoicePrintOperationDefaultOpts() *SSApiCallOpts {
	return NewSSApiCallOpts("invoice-print")
}

func InvoiceIssueOperationDefaultOpts() *SSApiCallOpts {
	return NewSSApiCallOpts("invoice-issue")
}

func InvoiceCancelOperationDefaultOpts() *SSApiCallOpts {
	return NewSSApiCallOpts("invoice-cancel")
}

func MetricsOperationDefaultOpts() *SSApiCallOpts {
	return NewSSApiCallOpts("metrics")
}
