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

package rules

import (
	"fmt"
	"regexp"

	"github.com/getgort/gort/types"
)

func Parse(rt RuleTokens) (Rule, error) {
	infer := types.Inferrer{}.ComplexTypes(true).StrictStrings(true)

	r := Rule{
		Command:     rt.Command,
		Conditions:  []Expression{},
		Permissions: []Permission{},
	}

	lastCondition := Undefined

	for _, p := range rt.Permissions {
		if p == "and" {
			lastCondition = And
			continue
		}

		if p == "or" {
			lastCondition = Or
			continue
		}

		r.Permissions = append(r.Permissions, Permission{
			Name:      p,
			Condition: lastCondition})
	}

	lastCondition = Undefined

	for _, c := range rt.Conditions {
		if c == "and" {
			lastCondition = And
			continue
		}

		if c == "or" {
			lastCondition = Or
			continue
		}

		a, b, o, m, err := ParseExpression(c)
		if err != nil {
			return r, fmt.Errorf("can't parse condition: %w", err)
		}

		va, err := infer.Infer(a)
		if err != nil {
			return r, fmt.Errorf("can't infer value: %w", err)
		}

		vb, err := infer.Infer(b)
		if err != nil {
			return r, fmt.Errorf("can't infer value: %w", err)
		}

		r.Conditions = append(r.Conditions, Expression{
			A:         va,
			B:         vb,
			Operator:  o,
			Modifier:  m,
			Condition: lastCondition})
	}

	return r, nil
}

var (
	reOperatorParts = regexp.MustCompile(`^(?:(all|any)\s+)?(.*)\s+([!<>=]{1,2}|in)\s+(.*)$`)
)

func ParseExpression(expr string) (a, b string, o Operator, m CollectionOperationModifier, err error) {
	subs := reOperatorParts.FindStringSubmatch(expr)

	if len(subs) != 5 {
		err = fmt.Errorf("expression doesn't conform to form A OP B")
		return
	}

	modifier := subs[1]
	op := subs[3]
	a, b = subs[2], subs[4]

	switch op {
	case "==":
		o = Equals
	case "!=":
		o = NotEquals
	case "<":
		o = LessThan
	case "<=":
		o = LessThanOrEqualTo
	case ">":
		o = GreaterThan
	case ">=":
		o = GreaterThanOrEqualTo
	case "in":
		o = In
	default:
		err = fmt.Errorf("unsupported operator: %s", op)
	}

	switch modifier {
	case "all":
		m = CollAll
	case "any":
		m = CollAny
	default:
		m = CollOne
	}

	return
}

// TokenizeAndParse is a helper function that wraps the Tokenize and Parse
// functions. It accepts a raw Gort rule of the form "COMMAND [when CONDITION
// (and|or)]? [allow|must have PERMISSION (and|or)]", and returns a RuleTokens
// value. A parsing error will produce a non-nil error. The RuleTokens' Command
// value should always be non-empty; Conditions and Permissions can both be
// empty (but non-nil). Empty Conditions always match the command. Empty
// Permissions indicating the use of the "allow" keyword and always pass.
func TokenizeAndParse(s string) (Rule, error) {
	rt, err := Tokenize(s)
	if err != nil {
		return Rule{}, err
	}

	return Parse(rt)
}
