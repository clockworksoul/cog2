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

package memory

import (
	"context"

	"github.com/getgort/gort/data/rest"
	"github.com/getgort/gort/dataaccess/errs"
)

// RoleCreate creates a new role.
func (da *InMemoryDataAccess) RoleCreate(ctx context.Context, name string) error {
	if name == "" {
		return errs.ErrEmptyRoleName
	}

	if nil != da.roles[name] {
		return errs.ErrRoleExists
	}

	da.roles[name] = &rest.Role{Name: name, Permissions: []rest.RolePermission{}}
	return nil
}

// RoleDelete
func (da *InMemoryDataAccess) RoleDelete(ctx context.Context, name string) error {
	if name == "" {
		return errs.ErrEmptyRoleName
	}

	if nil == da.roles[name] {
		return errs.ErrNoSuchRole
	}

	delete(da.roles, name)
	return nil
}

// RoleExists is used to determine whether a group exists in the data store.
func (da *InMemoryDataAccess) RoleExists(ctx context.Context, name string) (bool, error) {
	if name == "" {
		return false, errs.ErrEmptyRoleName
	}

	return da.roles[name] != nil, nil
}

// RoleGet gets a specific group.
func (da *InMemoryDataAccess) RoleGet(ctx context.Context, name string) (rest.Role, error) {
	role, ok := da.roles[name]

	if name == "" {
		return rest.Role{}, errs.ErrEmptyRoleName
	}

	if !ok {
		return rest.Role{}, errs.ErrNoSuchRole
	}

	return *role, nil
}

func (da *InMemoryDataAccess) RoleGrantPermission(ctx context.Context, rolename, bundle, permission string) error {
	role, ok := da.roles[rolename]

	if !ok {
		return errs.ErrNoSuchRole
	}

	role.Permissions = append(role.Permissions, rest.RolePermission{BundleName: bundle, Permission: permission})
	return nil
}

func (da *InMemoryDataAccess) RoleRevokePermission(ctx context.Context, rolename, bundle, permission string) error {
	role, ok := da.roles[rolename]

	if !ok {
		return errs.ErrNoSuchRole
	}

	perms := []rest.RolePermission{}
	for _, p := range role.Permissions {
		if p.BundleName == bundle && p.Permission == permission {
			continue
		}

		perms = append(perms, p)
	}

	role.Permissions = perms

	return nil
}
