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
	"sort"

	"github.com/getgort/gort/data/rest"
	"github.com/getgort/gort/dataaccess/errs"
)

// RoleCreate creates a new role.
func (da *InMemoryDataAccess) RoleCreate(ctx context.Context, rolename string) error {
	if rolename == "" {
		return errs.ErrEmptyRoleName
	}

	if nil != da.roles[rolename] {
		return errs.ErrRoleExists
	}

	da.roles[rolename] = &rest.Role{Name: rolename, Permissions: []rest.RolePermission{}}
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
func (da *InMemoryDataAccess) RoleGet(ctx context.Context, rolename string) (rest.Role, error) {
	role, ok := da.roles[rolename]

	if rolename == "" {
		return rest.Role{}, errs.ErrEmptyRoleName
	}

	if !ok {
		return rest.Role{}, errs.ErrNoSuchRole
	}

	return *role, nil
}

// RolePermissionExists returns true if the given role has been granted the
// specified permission. It returns an error if rolename is empty or if no
// such role exists.
func (da *InMemoryDataAccess) RolePermissionExists(ctx context.Context, rolename, bundlename, permission string) (bool, error) {
	perms, err := da.RolePermissionList(ctx, rolename)
	if err != nil {
		return false, err
	}

	for _, p := range perms {
		if p.BundleName == bundlename && p.Permission == permission {
			return true, nil
		}
	}

	return false, nil
}

// RoleList
func (da *InMemoryDataAccess) RoleList(ctx context.Context) ([]rest.Role, error) {
	list := make([]rest.Role, 0)

	for _, r := range da.roles {
		list = append(list, *r)
	}

	return list, nil
}

func (da *InMemoryDataAccess) RoleGroupAdd(ctx context.Context, rolename, groupname string) error {
	return da.GroupRoleAdd(ctx, groupname, rolename)
}

func (da *InMemoryDataAccess) RoleGroupDelete(ctx context.Context, rolename, groupname string) error {
	return da.GroupRoleDelete(ctx, groupname, rolename)
}

func (da *InMemoryDataAccess) RoleGroupExists(ctx context.Context, rolename, groupname string) (bool, error) {
	groups, err := da.RoleGroupList(ctx, rolename)
	if err != nil {
		return false, err
	}

	if exists, err := da.GroupExists(ctx, groupname); err != nil {
		return false, err
	} else if !exists {
		return false, errs.ErrNoSuchGroup
	}

	for _, g := range groups {
		if g.Name == groupname {
			return true, nil
		}
	}

	return false, nil
}

func (da *InMemoryDataAccess) RoleGroupList(ctx context.Context, rolename string) ([]rest.Group, error) {
	role, ok := da.roles[rolename]
	if !ok {
		return nil, errs.ErrNoSuchRole
	}
	return role.Groups, nil
}

func (da *InMemoryDataAccess) RolePermissionAdd(ctx context.Context, rolename, bundlename, permission string) error {
	role, ok := da.roles[rolename]
	if !ok {
		return errs.ErrNoSuchRole
	}

	role.Permissions = append(role.Permissions, rest.RolePermission{BundleName: bundlename, Permission: permission})

	return nil
}

func (da *InMemoryDataAccess) RolePermissionDelete(ctx context.Context, rolename, bundlename, permission string) error {
	role, ok := da.roles[rolename]

	if !ok {
		return errs.ErrNoSuchRole
	}

	perms := []rest.RolePermission{}
	for _, p := range role.Permissions {
		if p.BundleName == bundlename && p.Permission == permission {
			continue
		}

		perms = append(perms, p)
	}

	role.Permissions = perms

	return nil
}

// RolePermissionList returns returns an alphabetically-sorted list of
// fully-qualified (i.e., "bundle:permission") permissions granted to
// the role.
func (da *InMemoryDataAccess) RolePermissionList(ctx context.Context, rolename string) (rest.RolePermissionList, error) {
	role, err := da.RoleGet(ctx, rolename)
	if err != nil {
		return nil, err
	}

	perms := role.Permissions

	sort.Slice(perms, func(i, j int) bool { return perms[i].String() < perms[j].String() })

	return perms, nil
}
