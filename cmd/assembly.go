package commands

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/SAP/cloud-mta-build-tool/internal/archive"
	"github.com/SAP/cloud-mta-build-tool/internal/artifacts"
	"github.com/SAP/cloud-mta-build-tool/internal/logs"
	"strconv"
)

const (
	defaultPlatform string = "cf"
)

var assembleCmdSrc string
var assembleCmdTrg string

func init() {
	assemblyCommand.Flags().StringVarP(&assembleCmdSrc,
		"source", "s", "", "the path to the MTA project; the current path is default")
	assemblyCommand.Flags().StringVarP(&assembleCmdTrg,
		"target", "t", "", "the path to the MBT results folder; the current path is default")
}

// Generate mtar from build artifacts
var assemblyCommand = &cobra.Command{
	Use:       "assemble",
	Short:     "Assembles MTA Archive",
	Long:      "Assembles MTA Archive",
	ValidArgs: []string{"Deployment descriptor location"},
	Args:      cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := assembly(assembleCmdSrc, assembleCmdTrg, defaultPlatform, os.Getwd)
		logError(err)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func assembly(source, target, platform string, getWd func() (string, error)) error {
	logs.Logger.Info("assembling the MTA project...")
	// copy from source to target
	err := artifacts.CopyMtaContent(source, target, dir.Dep, getWd)
	if err != nil {
		return errors.Wrap(err, "assembly of the MTA project failed when copying the MTA content")
	}
	// Generate meta artifacts
	err = artifacts.ExecuteGenMeta(source, target, dir.Dep, platform, false, getWd)
	if err != nil {
		return errors.Wrap(err, "assembly of the MTA project failed when generating the meta information")
	}
	// generate mtar
	err = artifacts.ExecuteGenMtar(source, target, strconv.FormatBool(target != ""), dir.Dep, "", getWd)
	if err != nil {
		return errors.Wrap(err, "assembly of the MTA project failed when generating the MTA archive")
	}
	// cleanup
	err = artifacts.ExecuteCleanup(source, target, dir.Dep, getWd)
	if err != nil {
		return errors.Wrap(err, "assembly of the MTA project failed when executing cleanup")
	}
	return nil
}
