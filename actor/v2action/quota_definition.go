package v2action

// TODO change to ccv2.QuotaDefinition
type QuotaDefinition struct {
	Name                    string
	InstanceMemoryLimit     int
	MemoryLimit             int
	TotalRoutes             int
	TotalServices           int
	NonBasicServicesAllowed bool
	AppInstanceLimit        int
	TotalReservedRoutePorts int
}

func (actor Actor) GetQuotaDefinition(quotaDefinitionGUID string) (QuotaDefinition, Warnings, error) {
	return QuotaDefinition{}, nil, nil
}
