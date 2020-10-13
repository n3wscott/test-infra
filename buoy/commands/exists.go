/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"knative.dev/test-infra/pkg/gomod"
)

func addNextCmd(root *cobra.Command) {
	var domain string
	var release string
	var verbose bool
	var tag bool

	var cmd = &cobra.Command{
		Use:   "exists go.mod",
		Short: "Determine if the release branch exists for a given module.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			gomodFile := args[0]

			var out io.Writer
			if verbose {
				out = os.Stderr
			}

			meta, err := gomod.ReleaseStatus(gomodFile, release, domain, out)
			if err != nil {
				return err
			}

			if tag {
				fmt.Printf(meta.Release)
			}

			if !meta.ReleaseBranchExists {
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&domain, "domain", "d", "", "domain filter (i.e. knative.dev) [required]")
	_ = cmd.MarkFlagRequired("domain")
	cmd.Flags().StringVarP(&release, "release", "r", "", "release should be '<major>.<minor>' (i.e.: 1.23 or v1.23) [required]")
	_ = cmd.MarkFlagRequired("release")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print verbose output (stderr)")
	cmd.Flags().BoolVarP(&tag, "next", "t", false, "Print the next release tag (stdout)")

	root.AddCommand(cmd)
}