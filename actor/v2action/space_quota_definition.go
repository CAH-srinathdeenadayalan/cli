package v2action

// TODO move to ccv2.SpaceQuotaDefinition
type SpaceQuotaDefinition struct {
	Name string
}

func (actor Actor) GetOrganizationSpaceQuotaDefinitions(orgGUID string) ([]SpaceQuotaDefinition, Warnings, error) {

	return []SpaceQuotaDefinition{}, nil, nil
}
