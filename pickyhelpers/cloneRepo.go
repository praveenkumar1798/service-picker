package pickyhelpers

import (
	"github.com/wednesday-solutions/picky/utils"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func CloneRepo(stack, dirName, path string) error {

	// Download the selected stack.
	err := utils.RunCommandWithoutLogs("", "git", "clone", constants.Repos()[stack], dirName)
	errorhandler.CheckNilErr(err)

	// Delete cd.yml file from the cloned repo.
	cdFilePatch := path + "/" + dirName + constants.CDFilePathURL
	status, _ := fileutils.IsExists(cdFilePatch)
	if status {
		err = fileutils.RemoveFile(cdFilePatch)
		errorhandler.CheckNilErr(err)
	}
	return nil
}
