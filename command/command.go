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

package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/getgort/gort/types"
)

var (
	// ErrInvalidBundleCommandPair is returned by FindCommandEntry when the
	// command entry string doesn't look like  "command" or "bundle:command".
	ErrInvalidBundleCommandPair = errors.New("invalid bundle:comand pair")
)

// Command represents a command typed in by a user. It is typically
// generated by the Parse function.
type Command struct {
	Bundle     string
	Command    string
	Options    map[string]CommandOption
	Parameters CommandParameters
}

func (c Command) OptionsValues() map[string]types.Value {
	m := map[string]types.Value{}

	for _, o := range c.Options {
		m[o.Name] = o.Value
	}

	return m
}

// CommandOption represents a command option or flag, and its string
// value (if any).
type CommandOption struct {
	Name  string
	Value types.Value
}

type CommandParameters []types.Value

func (c CommandParameters) String() string {
	if len(c) == 0 {
		return ""
	}

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%v", c[0]))

	for i := 0; i < len(c); i++ {
		b.WriteRune(' ')
		b.WriteString(fmt.Sprintf("%v", c[i]))
	}

	return b.String()
}

// Parse accepts a slice of token strings and constructs a Command value.
// Its behavior may be modified by passing one or more ParseOptions.
func Parse(tokens []string, options ...ParseOption) (Command, error) {
	infer := types.Inferrer{}.ComplexTypes(false).StrictStrings(false)

	po := &parseOptions{false, false, map[string]string{}, map[string]bool{}}
	for _, o := range options {
		o(po)
	}

	if len(tokens) == 0 {
		return Command{}, fmt.Errorf("empty tokens list")
	}

	bundleName, commandName, err := SplitCommand(tokens[0])
	if err != nil {
		return Command{}, fmt.Errorf("command parse failure: %w", err)
	}

	cmd := Command{
		Bundle:     bundleName,
		Command:    commandName,
		Options:    map[string]CommandOption{},
		Parameters: []types.Value{},
	}

	tokens = tokens[1:]

	var lastOption *CommandOption = nil

	for i, t := range tokens {
		// Double slash indicates the end of options
		if t == "--" {
			cmd.Parameters, err = infer.InferAll(tokens[i+1:])
			if err != nil {
				return cmd, err
			}
			break
		}

		// Format: --option
		if len(t) >= 2 && dashCount(t) == 2 {
			lastOption = buildOption(t[2:], po)
			cmd.Options[lastOption.Name] = *lastOption
			continue
		}

		// Format: -I or -Ik
		if len(t) >= 1 && dashCount(t) == 1 {
			if po.agnosticDashes {
				lastOption = buildOption(t[1:], po)
				cmd.Options[lastOption.Name] = *lastOption
			} else {
				for _, ch := range t[1:] {
					lastOption = buildOption(string(ch), po)
					cmd.Options[lastOption.Name] = *lastOption
				}
			}

			continue
		}

		// If we got here, the token isn't an option.

		// Was the previous token an option?
		if lastOption != nil {
			// Does it expect an argument?
			explicitHas, ok := po.hasArg[lastOption.Name]
			hasArgument := (ok && explicitHas) || (!ok && po.assumeOptionArguments)

			// Expect an option:
			if hasArgument {
				term, err := infer.Infer(t)
				if err != nil {
					return cmd, err
				}

				lastOption.Value = term
				cmd.Options[lastOption.Name] = *lastOption
				lastOption = nil
				continue
			}
		}

		// Not an option; not an argument. Must be command args.
		cmd.Parameters, err = infer.InferAll(tokens[i:])
		if err != nil {
			return cmd, err
		}
		break
	}

	return cmd, nil
}

type parseOptions struct {
	agnosticDashes        bool
	assumeOptionArguments bool
	aliases               map[string]string
	hasArg                map[string]bool
}

type ParseOption func(*parseOptions)

// ParseAgnosticDashes modifies how dashes are interpreted. If true, double and
// single dashes are treated the same. If false (default), then double-dashed
//options are "long" and single-dashed options are "short".
func ParseAgnosticDashes(agnostic bool) ParseOption {
	return func(po *parseOptions) {
		po.agnosticDashes = agnostic
	}
}

// ParseAssumeOptionArguments changes the assumption about whether an option
// not specified using ParseOptionHasArgument should be treated as having an
// option. If false (default), unknown options are assumed to have no argument.
func ParseAssumeOptionArguments(assume bool) ParseOption {
	return func(po *parseOptions) {
		po.assumeOptionArguments = assume
	}
}

// ParseOptionHasArgument allows specific options to be specified as expecting
// an option (or not). Options not specified are treated according to
// ParseAssumeOptionArguments.
func ParseOptionHasArgument(option string, hasArg bool) ParseOption {
	return func(po *parseOptions) {
		po.hasArg[option] = hasArg
	}
}

// ParseOptionAlias allows option aliases to be set, most often "short options"
// to "long options". All references to "alias" are treated as "name".
func ParseOptionAlias(alias, name string) ParseOption {
	return func(po *parseOptions) {
		po.aliases[alias] = name
	}
}

// SplitCommand accepts a string in the style of "bundle:command" or "command"
// and returns the bundle and command as a pair of strings. If there's no
// indicated bundle, the bundle string (the first string) will be empty. If
// there's more than one colon, an error will be returned.
func SplitCommand(name string) (bundle, command string, err error) {
	split := strings.Split(name, ":")

	switch len(split) {
	case 1:
		command = split[0]
	case 2:
		bundle = split[0]
		command = split[1]
	default:
		err = ErrInvalidBundleCommandPair
	}

	return
}

func buildOption(name string, po *parseOptions) *CommandOption {
	if n, ok := po.aliases[name]; ok {
		name = n
	}

	return &CommandOption{Name: name, Value: types.BoolValue{V: true}}
}

func dashCount(str string) int {
	count := 0

	for _, ch := range str {
		if ch != '-' {
			break
		}

		count++
	}

	return count
}

// TokenizeAndParse is a helper function that combines the Tokenize and Parse functions.
func TokenizeAndParse(str string, options ...ParseOption) (Command, error) {
	t, err := Tokenize(str)
	if err != nil {
		return Command{}, err
	}

	return Parse(t, options...)
}
