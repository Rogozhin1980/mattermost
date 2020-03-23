// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package commands

import (
	"errors"
	"io/ioutil"

	"github.com/mattermost/mattermost-server/v5/audit"
	"github.com/spf13/cobra"
)

var LicenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Licensing commands",
}

var UploadLicenseCmd = &cobra.Command{
	Use:     "upload [license]",
	Short:   "Upload a license.",
	Long:    "Upload a license. Replaces current license.",
	Example: "  license upload /path/to/license/mylicensefile.mattermost-license",
	RunE:    uploadLicenseCmdF,
}

func init() {
	LicenseCmd.AddCommand(UploadLicenseCmd)
	RootCmd.AddCommand(LicenseCmd)
}

func uploadLicenseCmdF(command *cobra.Command, args []string) (cmdError error) {
	a, err := InitDBCommandContextCobra(command)
	if err != nil {
		return err
	}
	defer a.Shutdown()

	if len(args) != 1 {
		return errors.New("Enter one license file to upload")
	}

	auditRec := a.MakeAuditRecord("uploadLicense", audit.Fail)
	defer func() { a.LogAuditRec(auditRec, cmdError) }()
	auditRec.AddMeta("file", args[0])

	var fileBytes []byte
	if fileBytes, err = ioutil.ReadFile(args[0]); err != nil {
		return err
	}

	if _, err := a.SaveLicense(fileBytes); err != nil {
		return err
	}

	CommandPrettyPrintln("Uploaded license file")
	auditRec.Success()

	return nil
}
