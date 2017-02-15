package v2

import (
	"fmt"
	"sort"
	"strings"

	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/flag"
	"code.cloudfoundry.org/cli/command/v2/shared"
)

//go:generate counterfeiter . OrgActor

type OrgActor interface {
	GetOrganizationByName(orgName string) (v2action.Organization, v2action.Warnings, error)
	// GetOrganizationPrivateDomains(orgGUID string) ([]v2action.Domain, v2action.Warnings, error)
	// GetSharedDomains() ([]v2action.Domain, v2action.Warnings, error)
	GetOrganizationDomainNames(orgGUID string) ([]string, v2action.Warnings, error)
	GetQuotaDefinition(quotaDefinitionGUID string) (v2action.QuotaDefinition, v2action.Warnings, error)
	GetOrganizationSpaces(orgGUID string) ([]v2action.Space, v2action.Warnings, error)
	GetOrganizationSpaceQuotaDefinitions(orgGUID string) ([]v2action.SpaceQuotaDefinition, v2action.Warnings, error)
}

type OrgCommand struct {
	RequiredArgs    flag.Organization `positional-args:"yes"`
	GUID            bool              `long:"guid" description:"Retrieve and display the given org's guid.  All other output for the org is suppressed."`
	usage           interface{}       `usage:"CF_NAME org ORG [--guid]"`
	relatedCommands interface{}       `related_commands:"org-users, orgs"`

	UI     command.UI
	Config command.Config
	Actor  OrgActor
}

func (cmd *OrgCommand) Setup(config command.Config, ui command.UI) error {
	return nil
}

func (cmd OrgCommand) Execute(args []string) error {
	if cmd.GUID {
		return cmd.displayOrgGUID()
	} else {
		return cmd.displayOrgSummary()
	}
}

func (cmd OrgCommand) displayOrgGUID() error {
	org, warnings, err := cmd.Actor.GetOrganizationByName(cmd.RequiredArgs.Organization)
	cmd.UI.DisplayWarnings(warnings)
	if err != nil {
		return shared.HandleError(err)
	}

	cmd.UI.DisplayText(org.GUID)

	return nil
}

func (cmd OrgCommand) displayOrgSummary() error {
	user, err := cmd.Config.CurrentUser()
	if err != nil {
		return shared.HandleError(err)
	}

	cmd.UI.DisplayTextWithFlavor(
		"Getting info for org {{.OrgName}} as {{.Username}}...",
		map[string]interface{}{
			"OrgName":  cmd.RequiredArgs.Organization,
			"Username": user.Name,
		})
	cmd.UI.DisplayNewline()

	org, warnings, err := cmd.Actor.GetOrganizationByName(cmd.RequiredArgs.Organization)
	cmd.UI.DisplayWarnings(warnings)
	if err != nil {
		return shared.HandleError(err)
	}

	domainNames, warnings, err := cmd.Actor.GetOrganizationDomainNames(org.GUID)
	cmd.UI.DisplayWarnings(warnings)
	if err != nil {
		return shared.HandleError(err)
	}

	quotaDefinition, warnings, err := cmd.Actor.GetQuotaDefinition(org.QuotaDefinitionGUID)
	cmd.UI.DisplayWarnings(warnings)
	if err != nil {
		return shared.HandleError(err)
	}

	spaces, warnings, err := cmd.Actor.GetOrganizationSpaces(org.GUID)
	cmd.UI.DisplayWarnings(warnings)
	if err != nil {
		return shared.HandleError(err)
	}

	var spaceNames []string
	for _, space := range spaces {
		spaceNames = append(spaceNames, space.Name)
	}

	sort.Strings(spaceNames)

	spaceQuotaDefinitions, warnings, err := cmd.Actor.GetOrganizationSpaceQuotaDefinitions(org.GUID)
	cmd.UI.DisplayWarnings(warnings)
	if err != nil {
		return shared.HandleError(err)
	}

	var spaceQuotaDefinitionNames []string
	for _, spaceQuotaDefinition := range spaceQuotaDefinitions {
		spaceQuotaDefinitionNames = append(spaceQuotaDefinitionNames, spaceQuotaDefinition.Name)
	}

	sort.Strings(spaceQuotaDefinitionNames)

	cmd.UI.DisplayText(fmt.Sprintf("%s:", org.Name))

	table := [][]string{
		{cmd.UI.TranslateText("domains:"), strings.Join(domainNames, ", ")},
		{cmd.UI.TranslateText("quota:"), cmd.formatQuotaDefinition(quotaDefinition)},
		{cmd.UI.TranslateText("spaces:"), strings.Join(spaceNames, ", ")},
		{cmd.UI.TranslateText("space quotas:"), strings.Join(spaceQuotaDefinitionNames, ", ")},
	}
	cmd.UI.DisplayTable("", table, 3)

	cmd.UI.DisplayOK()

	return nil
}

func (cmd OrgCommand) formatQuotaDefinition(quotaDefinition v2action.QuotaDefinition) string {
	parts := []string{}

	parts = append(parts, cmd.UI.TranslateText("{{.MemoryLimit}}M memory limit", map[string]interface{}{
		"MemoryLimit": quotaDefinition.MemoryLimit,
	}))

	parts = append(parts, cmd.UI.TranslateText("{{.InstanceMemoryLimit}}M instance memory limit", map[string]interface{}{
		"InstanceMemoryLimit": quotaDefinition.InstanceMemoryLimit,
	}))

	parts = append(parts, cmd.UI.TranslateText("{{.RoutesLimit}} routes", map[string]interface{}{
		"RoutesLimit": quotaDefinition.TotalRoutes,
	}))

	parts = append(parts, cmd.UI.TranslateText("{{.ServicesLimit}} services", map[string]interface{}{
		"ServicesLimit": quotaDefinition.TotalServices,
	}))

	var paidAllowed string
	if quotaDefinition.NonBasicServicesAllowed {
		paidAllowed = cmd.UI.TranslateText("allowed")
	} else {
		paidAllowed = cmd.UI.TranslateText("disallowed")
	}
	parts = append(parts, cmd.UI.TranslateText("paid services {{.NonBasicServicesAllowed}}", map[string]interface{}{
		"NonBasicServicesAllowed": paidAllowed,
	}))

	var appInstanceLimit string
	if quotaDefinition.AppInstanceLimit == -1 {
		appInstanceLimit = cmd.UI.TranslateText("unlimited")
	} else {
		appInstanceLimit = fmt.Sprintf("%d", quotaDefinition.AppInstanceLimit)
	}
	parts = append(parts, cmd.UI.TranslateText("{{.AppInstanceLimit}} app instance limit", map[string]interface{}{
		"AppInstanceLimit": appInstanceLimit,
	}))

	var routePorts string
	if quotaDefinition.TotalReservedRoutePorts == -1 {
		routePorts = cmd.UI.TranslateText("unlimited")
	} else {
		routePorts = fmt.Sprintf("%d", quotaDefinition.TotalReservedRoutePorts)
	}
	parts = append(parts, cmd.UI.TranslateText("{{.ReservedRoutePorts}} route ports", map[string]interface{}{
		"ReservedRoutePorts": routePorts,
	}))

	return fmt.Sprintf("%s (%s)", quotaDefinition.Name, strings.Join(parts, ", "))
}
