// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"log"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallkanister() *cobra.Command {
	var kanisterApp = &cobra.Command{
		Use:          "kanister",
		Short:        "Install kanister for application-level data management",
		Long:         "Install kanister, An extensible open-source framework for application-level data management on Kubernetes",
		Example:      "arkade install kanister",
		SilenceUsage: true,
	}

	kanisterApp.Flags().StringP("namespace", "n", "default", "The namespace to install kanister (default: default")
	kanisterApp.Flags().Bool("update-repo", true, "Update the helm repo")
	kanisterApp.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set image.tag=0.69.0)")

	kanisterApp.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		log.Println(kubeConfigPath)
		namespace, _ := kanisterApp.Flags().GetString("namespace")

		overrides := map[string]string{}

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		kanisterOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("kanister/kanister-operator").
			WithHelmURL("https://charts.kanister.io/").
			WithOverrides(overrides).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(kanisterOptions)
		if err != nil {
			return err
		}

		println(KanisterInstallMsg)
		return nil
	}

	return kanisterApp
}

const KanisterInfoMsg = `# Get started with kanister here:
# https://kanister.io/

# See kanister docs here
# https://docs.kanister.io/install.html

# Once deployed use kanctl to create an S3 compliant kanister profile:

kanctl create profile s3compliant --bucket <bucket> --access-key ${AWS_ACCESS_KEY_ID} \
                                    --secret-key ${AWS_SECRET_ACCESS_KEY}               \
                                    --region <region>                                   \
                                    --namespace kanister


`

const KanisterInstallMsg = `=======================================================================
= kanister has been installed.                                   =
=======================================================================` +
	"\n\n" + KanisterInfoMsg + "\n\n" + pkg.ThanksForUsing
