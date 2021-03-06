/*
 * Copyright 2021 The Gort Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/getgort/gort/client"
)

// $ cogctl group grant --help
// Usage: cogctl group grant [OPTIONS] GROUP ROLES...
//
//   Grant one or more roles to an existing group.
//
// Options:
//   --help  Show this message and exit.

const (
	groupGrantUse   = "grant"
	groupGrantShort = "Grant a role to an existing group"
	groupGrantLong  = "Grant a role to an existing group."
	groupGrantUsage = `Usage:
  gort group grant [flags] group_name role_name

Flags:
  -h, --help   Show this message and exit

Global Flags:
  -P, --profile string   The Gort profile within the config file to use
`
)

// GetGroupGrantCmd is a command
func GetGroupGrantCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   groupGrantUse,
		Short: groupGrantShort,
		Long:  groupGrantLong,
		RunE:  groupGrantCmd,
		Args:  cobra.ExactArgs(2),
	}

	cmd.SetUsageTemplate(groupGrantUsage)

	return cmd
}

func groupGrantCmd(cmd *cobra.Command, args []string) error {
	groupname := args[0]
	rolename := args[1]

	gortClient, err := client.Connect(FlagGortProfile)
	if err != nil {
		return err
	}

	err = gortClient.GroupRoleAdd(groupname, rolename)
	if err != nil {
		return err
	}

	fmt.Printf("role added to %s: %s\n", groupname, rolename)

	return nil
}
