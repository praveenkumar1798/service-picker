package prompt

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
	"github.com/wednesday-solutions/picky/pickyhelpers"
)

// PromptSetupInfra is the prompt for the setup infra option of Home prompt.
func PromptSetupInfra() {
	var p PromptInput
	p.Label = "Do you want to setup infrastructure for your project"
	p.GoBack = PromptHome
	response := p.PromptYesOrNoSelect()
	if response {
		cloudProvider := PromptCloudProvider()
		var stacks []string
		for {
			stacks = PromptSelectExistingStacks()
			if len(stacks) > 0 {
				break
			}
		}
		environment := PromptEnvironment()
		err := CreateInfra(stacks, cloudProvider, environment)
		errorhandler.CheckNilErr(err)
		PromptDeployAfterInfra(stacks, environment)
	}
	PromptHome()
}

// PromptCloudProvider is a prompt for selecting a cloud provider.
func PromptCloudProvider() string {
	var p PromptInput
	p.Label = "Choose a cloud provider"
	p.Items = []string{constants.AWS}
	p.GoBack = PromptHome
	return p.PromptSelect()
}

// PromptEnvironment is a prompt for selecting an environment.
func PromptEnvironment() string {
	var p PromptInput
	p.Label = "Choose an environment"
	p.Items = []string{constants.Development, constants.QA, constants.Production}
	p.GoBack = PromptHome
	return p.PromptSelect()
}

// CreateInfra execute all the functionalities of infra setup.
func CreateInfra(directories []string, cloudProvider string, environment string) error {
	switch cloudProvider {
	case constants.AWS:
		status := pickyhelpers.IsInfraFilesExist()
		var (
			stack, database string
			stackInfo       map[string]interface{}
			err             error
		)
		done := make(chan bool)
		go pickyhelpers.ProgressBar(20, "Generating", done)

		if !status {
			err = pickyhelpers.CreateInfraSetup()
			errorhandler.CheckNilErr(err)
		}
		for _, dirName := range directories {
			service := utils.FindService(dirName)
			stack, database = utils.FindStackAndDatabase(dirName)
			stackInfo = pickyhelpers.GetStackInfo(stack, database, environment)

			err = pickyhelpers.CreateInfraStacks(service, stack, database, dirName, environment)
			if err != nil {
				if err.Error() != errorhandler.ErrExist.Error() {
					errorhandler.CheckNilErr(err)
				}
			}
			if service == constants.Backend {
				err = pickyhelpers.UpdateEnvByEnvironment(dirName, environment)
				errorhandler.CheckNilErr(err)
			}
		}
		err = pickyhelpers.CreateSstConfigFile(stackInfo, directories)
		errorhandler.CheckNilErr(err)
		<-done
		fmt.Printf("\n%s %s", "Generating", errorhandler.CompleteMessage)
	default:
		fmt.Printf("\nWork in Progress. Please stay tuned..!\n")
	}
	return nil
}
