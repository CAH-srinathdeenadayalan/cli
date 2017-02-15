package v2_test

import (
	"errors"

	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/command/commandfakes"
	"code.cloudfoundry.org/cli/command/v2"
	"code.cloudfoundry.org/cli/command/v2/shared"
	"code.cloudfoundry.org/cli/command/v2/v2fakes"
	"code.cloudfoundry.org/cli/util/configv3"
	"code.cloudfoundry.org/cli/util/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = FDescribe("org Command", func() {
	var (
		cmd        v2.OrgCommand
		testUI     *ui.UI
		fakeConfig *commandfakes.FakeConfig
		// 	fakeSharedActor *commandfakes.FakeSharedActor
		fakeActor *v2fakes.FakeOrgActor
		// 	binaryName      string
		executeErr error
	)

	BeforeEach(func() {
		testUI = ui.NewTestUI(nil, NewBuffer(), NewBuffer())
		fakeConfig = new(commandfakes.FakeConfig)
		// 	fakeSharedActor = new(commandfakes.FakeSharedActor)
		fakeActor = new(v2fakes.FakeOrgActor)

		cmd = v2.OrgCommand{
			UI:     testUI,
			Config: fakeConfig,
			// SharedActor: fakeSharedActor,
			Actor: fakeActor,
		}

		// 	cmd.RequiredArgs.AppName = "some-app"

		// 	binaryName = "faceman"
		// 	fakeConfig.BinaryNameReturns(binaryName)

		// 	// TODO: remove when experimental flag is removed
		// 	fakeConfig.ExperimentalReturns(true)
		cmd.RequiredArgs.Organization = "some-org"
	})

	JustBeforeEach(func() {
		executeErr = cmd.Execute(nil)
	})

	// // TODO: remove when experimental flag is removed
	// It("Displays the experimental warning message", func() {
	// 	Expect(testUI.Out).To(Say(command.ExperimentalWarning))
	// })

	// Context("when the user is logged in, and org and space are targeted", func() {
	// 	BeforeEach(func() {
	// 		fakeConfig.HasTargetedOrganizationReturns(true)
	// 		fakeConfig.TargetedOrganizationReturns(configv3.Organization{Name: "some-org"})
	// 		fakeConfig.HasTargetedSpaceReturns(true)
	// 		fakeConfig.TargetedSpaceReturns(configv3.Space{
	// 			GUID: "some-space-guid",
	// 			Name: "some-space"})
	// 		fakeConfig.CurrentUserReturns(
	// 			configv3.User{Name: "some-user"},
	// 			nil)
	// 	})

	// 	Context("when getting the current user returns an error", func() {
	// 		var expectedErr error

	// 		BeforeEach(func() {
	// 			expectedErr = errors.New("getting current user error")
	// 			fakeConfig.CurrentUserReturns(
	// 				configv3.User{},
	// 				expectedErr)
	// 		})

	// 		It("returns the error", func() {
	// 			Expect(executeErr).To(MatchError(expectedErr))
	// 		})
	// 	})

	// 	It("displays flavor text", func() {
	// 		Expect(testUI.Out).To(Say("Showing health and status for app some-app in org some-org / space some-space as some-user..."))
	// 	})

	Context("when the --guid flag is provided", func() {
		BeforeEach(func() {
			cmd.GUID = true
		})

		Context("when no errors occur", func() {
			BeforeEach(func() {
				fakeActor.GetOrganizationByNameReturns(
					v2action.Organization{GUID: "some-org-guid"},
					v2action.Warnings{"warning-1", "warning-2"},
					nil)
			})

			It("displays the org guid and outputs all warnings", func() {
				Expect(executeErr).ToNot(HaveOccurred())

				Expect(testUI.Out).To(Say("some-org-guid"))
				Expect(testUI.Err).To(Say("warning-1"))
				Expect(testUI.Err).To(Say("warning-2"))

				Expect(fakeActor.GetOrganizationByNameCallCount()).To(Equal(1))
				orgName := fakeActor.GetOrganizationByNameArgsForCall(0)
				Expect(orgName).To(Equal("some-org"))
			})
		})

		Context("when getting the org returns an error", func() {
			Context("when the error is translatable", func() {
				BeforeEach(func() {
					fakeActor.GetOrganizationByNameReturns(
						v2action.Organization{},
						v2action.Warnings{"warning-1", "warning-2"},
						v2action.OrganizationNotFoundError{Name: "some-org"})
				})

				It("returns a translatable error and outputs all warnings", func() {
					Expect(executeErr).To(MatchError(shared.OrganizationNotFoundError{Name: "some-org"}))

					Expect(testUI.Err).To(Say("warning-1"))
					Expect(testUI.Err).To(Say("warning-2"))
				})
			})

			Context("when the error is not translatable", func() {
				var expectedErr error

				BeforeEach(func() {
					expectedErr = errors.New("get org error")
					fakeActor.GetOrganizationByNameReturns(
						v2action.Organization{},
						v2action.Warnings{"warning-1", "warning-2"},
						expectedErr)
				})

				It("returns the error and all warnings", func() {
					Expect(executeErr).To(MatchError(expectedErr))

					Expect(testUI.Err).To(Say("warning-1"))
					Expect(testUI.Err).To(Say("warning-2"))
				})
			})
		})
	})

	Context("when the --guid flag is not provided", func() {
		Context("when no errors occur", func() {
			BeforeEach(func() {
				fakeConfig.CurrentUserReturns(
					configv3.User{
						Name: "some-user",
					},
					nil)

				fakeActor.GetOrganizationByNameReturns(
					v2action.Organization{
						Name:                "some-org",
						GUID:                "some-org-guid",
						QuotaDefinitionGUID: "some-quota-definition-guid",
					},
					v2action.Warnings{"warning-1", "warning-2"},
					nil)

				fakeActor.GetOrganizationDomainNamesReturns(
					[]string{
						"a-shared.com",
						"b-private.com",
						"c-shared.com",
						"d-private.com",
					},
					v2action.Warnings{"warning-3", "warning-4"},
					nil)

				fakeActor.GetQuotaDefinitionReturns(
					v2action.QuotaDefinition{
						Name:                    "default",
						InstanceMemoryLimit:     456,
						MemoryLimit:             123,
						TotalRoutes:             789,
						TotalServices:           987,
						NonBasicServicesAllowed: true,
						AppInstanceLimit:        654,
						TotalReservedRoutePorts: 321,
					},
					v2action.Warnings{"warning-5", "warning-6"},
					nil)

				fakeActor.GetOrganizationSpacesReturns(
					[]v2action.Space{
						{Name: "space2"},
						{Name: "space1"},
					},
					v2action.Warnings{"warning-7", "warning-8"},
					nil)

				fakeActor.GetOrganizationSpaceQuotaDefinitionsReturns(
					[]v2action.SpaceQuotaDefinition{
						{Name: "space-quota2"},
						{Name: "space-quota1"},
					},
					v2action.Warnings{"warning-9", "warning-10"},
					nil)
			})

			It("displays warnings and a table with org domains, quotas, spaces and space quotas", func() {
				Expect(executeErr).To(BeNil())

				Eventually(testUI.Out).Should(Say("Getting info for org %s as some-user\\.\\.\\.", cmd.RequiredArgs.Organization))
				Expect(testUI.Err).To(Say("warning-1"))
				Expect(testUI.Err).To(Say("warning-2"))
				Expect(testUI.Err).To(Say("warning-3"))
				Expect(testUI.Err).To(Say("warning-4"))
				Expect(testUI.Err).To(Say("warning-5"))
				Expect(testUI.Err).To(Say("warning-6"))
				Expect(testUI.Err).To(Say("warning-7"))
				Expect(testUI.Err).To(Say("warning-8"))
				Expect(testUI.Err).To(Say("warning-9"))
				Expect(testUI.Err).To(Say("warning-10"))

				Eventually(testUI.Out).Should(Say("%s:", cmd.RequiredArgs.Organization))

				Eventually(testUI.Out).Should(Say("domains:\\s+a-shared.com, b-private.com, c-shared.com, d-private.com"))

				Eventually(testUI.Out).Should(Say("quota:\\s+default \\(123M memory limit, 456M instance memory limit, 789 routes, 987 services, paid services allowed, 654 app instance limit, 321 route ports\\)"))

				Eventually(testUI.Out).Should(Say("spaces:\\s+space1, space2"))

				Eventually(testUI.Out).Should(Say("space quotas:\\s+space-quota1, space-quota2"))

				Eventually(testUI.Out).Should(Say("OK"))

				Expect(fakeConfig.CurrentUserCallCount()).To(Equal(1))

				Expect(fakeActor.GetOrganizationByNameCallCount()).To(Equal(1))
				orgName := fakeActor.GetOrganizationByNameArgsForCall(0)
				Expect(orgName).To(Equal("some-org"))

				Expect(fakeActor.GetOrganizationDomainNamesCallCount()).To(Equal(1))
				orgGUID := fakeActor.GetOrganizationDomainNamesArgsForCall(0)
				Expect(orgGUID).To(Equal("some-org-guid"))

				Expect(fakeActor.GetQuotaDefinitionCallCount()).To(Equal(1))
				quotaDefinitionGUID := fakeActor.GetQuotaDefinitionArgsForCall(0)
				Expect(quotaDefinitionGUID).To(Equal("some-quota-definition-guid"))

				Expect(fakeActor.GetOrganizationSpacesCallCount()).To(Equal(1))
				orgGUID = fakeActor.GetOrganizationSpacesArgsForCall(0)
				Expect(orgGUID).To(Equal("some-org-guid"))

				Expect(fakeActor.GetOrganizationSpaceQuotaDefinitionsCallCount()).To(Equal(1))
				orgGUID = fakeActor.GetOrganizationSpaceQuotaDefinitionsArgsForCall(0)
				Expect(orgGUID).To(Equal("some-org-guid"))
			})

			Context("when route ports and app instances are unlimited", func() {
				BeforeEach(func() {
					fakeActor.GetQuotaDefinitionReturns(
						v2action.QuotaDefinition{
							Name:                    "default",
							InstanceMemoryLimit:     456,
							MemoryLimit:             123,
							TotalRoutes:             789,
							TotalServices:           987,
							NonBasicServicesAllowed: true,
							AppInstanceLimit:        -1,
							TotalReservedRoutePorts: -1,
						},
						v2action.Warnings{"warning-5", "warning-6"},
						nil)
				})

				It("displays unlimited", func() {
					Expect(executeErr).To(BeNil())

					Eventually(testUI.Out).Should(Say("quota:\\s+default \\(123M memory limit, 456M instance memory limit, 789 routes, 987 services, paid services allowed, unlimited app instance limit, unlimited route ports\\)"))
				})
			})
		})

		Context("when getting the current user returns an error", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = errors.New("getting current user error")
				fakeConfig.CurrentUserReturns(
					configv3.User{},
					expectedErr)
			})

			It("returns the error", func() {
				Expect(executeErr).To(MatchError(expectedErr))
			})
		})

		Context("when getting the org returns an error", func() {
			Context("when the error is translatable", func() {
				BeforeEach(func() {
					fakeActor.GetOrganizationByNameReturns(
						v2action.Organization{},
						v2action.Warnings{"warning-1", "warning-2"},
						v2action.OrganizationNotFoundError{Name: "some-org"})
				})

				It("returns a translatable error and outputs all warnings", func() {
					Expect(executeErr).To(MatchError(shared.OrganizationNotFoundError{Name: "some-org"}))

					Expect(testUI.Err).To(Say("warning-1"))
					Expect(testUI.Err).To(Say("warning-2"))
				})
			})

			Context("when the error is not translatable", func() {
				var expectedErr error

				BeforeEach(func() {
					expectedErr = errors.New("get org error")
					fakeActor.GetOrganizationByNameReturns(
						v2action.Organization{},
						v2action.Warnings{"warning-1", "warning-2"},
						expectedErr)
				})

				It("returns the error and all warnings", func() {
					Expect(executeErr).To(MatchError(expectedErr))

					Expect(testUI.Err).To(Say("warning-1"))
					Expect(testUI.Err).To(Say("warning-2"))
				})
			})
		})

		Context("when getting the org domain names returns an error", func() {
			Context("when the error is translatable", func() {
				BeforeEach(func() {
					fakeActor.GetOrganizationDomainNamesReturns(
						nil,
						v2action.Warnings{"warning-1", "warning-2"},
						v2action.OrganizationNotFoundError{Name: "some-org"})
				})

				It("returns a translatable error and outputs all warnings", func() {
					Expect(executeErr).To(MatchError(shared.OrganizationNotFoundError{Name: "some-org"}))

					Expect(testUI.Err).To(Say("warning-1"))
					Expect(testUI.Err).To(Say("warning-2"))
				})
			})

			Context("when the error is not translatable", func() {
				var expectedErr error

				BeforeEach(func() {
					expectedErr = errors.New("get org domains error")
					fakeActor.GetOrganizationDomainNamesReturns(
						nil,
						v2action.Warnings{"warning-1", "warning-2"},
						expectedErr)
				})

				It("returns the error and all warnings", func() {
					Expect(executeErr).To(MatchError(expectedErr))

					Expect(testUI.Err).To(Say("warning-1"))
					Expect(testUI.Err).To(Say("warning-2"))
				})
			})
		})

		Context("when getting the quota definition returns an error", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = errors.New("get quota definition error")
				fakeActor.GetQuotaDefinitionReturns(
					v2action.QuotaDefinition{},
					v2action.Warnings{"warning-1", "warning-2"},
					expectedErr)
			})

			It("returns the error and all warnings", func() {
				Expect(executeErr).To(MatchError(expectedErr))

				Expect(testUI.Err).To(Say("warning-1"))
				Expect(testUI.Err).To(Say("warning-2"))
			})
		})

		Context("when getting the org spaces returns an error", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = errors.New("get org spaces error")
				fakeActor.GetOrganizationSpacesReturns(
					nil,
					v2action.Warnings{"warning-1", "warning-2"},
					expectedErr)
			})

			It("returns the error and all warnings", func() {
				Expect(executeErr).To(MatchError(expectedErr))

				Expect(testUI.Err).To(Say("warning-1"))
				Expect(testUI.Err).To(Say("warning-2"))
			})
		})

		Context("when getting the org space quota definitions returns an error", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = errors.New("get org space quota definitions error")
				fakeActor.GetOrganizationSpaceQuotaDefinitionsReturns(
					nil,
					v2action.Warnings{"warning-1", "warning-2"},
					expectedErr)
			})

			It("returns the error and all warnings", func() {
				Expect(executeErr).To(MatchError(expectedErr))

				Expect(testUI.Err).To(Say("warning-1"))
				Expect(testUI.Err).To(Say("warning-2"))
			})
		})
	})
})
