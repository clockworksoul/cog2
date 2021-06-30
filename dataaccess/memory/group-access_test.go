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
	"testing"

	"github.com/getgort/gort/data/rest"
	"github.com/getgort/gort/dataaccess/errs"
	"github.com/stretchr/testify/assert"
)

func testGroupAccess(t *testing.T) {
	t.Run("testGroupAddUser", testGroupAddUser)
	t.Run("testGroupCreate", testGroupCreate)
	t.Run("testGroupDelete", testGroupDelete)
	t.Run("testGroupExists", testGroupExists)
	t.Run("testGroupGet", testGroupGet)
	t.Run("testGroupGrantRole", testGroupGrantRole)
	t.Run("testGroupList", testGroupList)
	t.Run("testGroupListRoles", testGroupListRoles)
	t.Run("testGroupRemoveUser", testGroupRemoveUser)
}

func testGroupAddUser(t *testing.T) {
	err := da.GroupAddUser(ctx, "foo", "bar")
	assert.Error(t, err, errs.ErrNoSuchGroup)

	da.GroupCreate(ctx, rest.Group{Name: "foo"})
	defer da.GroupDelete(ctx, "foo")

	err = da.GroupAddUser(ctx, "foo", "bar")
	assert.Error(t, err, errs.ErrNoSuchUser)

	da.UserCreate(ctx, rest.User{Username: "bar", Email: "bar"})
	defer da.UserDelete(ctx, "bar")

	err = da.GroupAddUser(ctx, "foo", "bar")
	assert.NoError(t, err)

	group, _ := da.GroupGet(ctx, "foo")

	if len(group.Users) != 1 {
		t.Error("Users list empty")
		t.FailNow()
	}

	if len(group.Users) > 0 && group.Users[0].Username != "bar" {
		t.Error("Wrong user!")
		t.FailNow()
	}
}

func testGroupCreate(t *testing.T) {
	var err error
	var group rest.Group

	// Expect an error
	err = da.GroupCreate(ctx, group)
	assert.Error(t, err, errs.ErrEmptyGroupName)

	// Expect no error
	err = da.GroupCreate(ctx, rest.Group{Name: "test-create"})
	defer da.GroupDelete(ctx, "test-create")
	assert.NoError(t, err)

	// Expect an error
	err = da.GroupCreate(ctx, rest.Group{Name: "test-create"})
	assert.Error(t, err, errs.ErrGroupExists)
}

func testGroupDelete(t *testing.T) {
	// Delete blank group
	err := da.GroupDelete(ctx, "")
	assert.Error(t, err, errs.ErrEmptyGroupName)

	// Delete group that doesn't exist
	err = da.GroupDelete(ctx, "no-such-group")
	assert.Error(t, err, errs.ErrNoSuchGroup)

	da.GroupCreate(ctx, rest.Group{Name: "test-delete"}) // This has its own test
	defer da.GroupDelete(ctx, "test-delete")

	err = da.GroupDelete(ctx, "test-delete")
	assert.NoError(t, err)

	exists, _ := da.GroupExists(ctx, "test-delete")
	if exists {
		t.Error("Shouldn't exist anymore!")
		t.FailNow()
	}
}

func testGroupExists(t *testing.T) {
	var exists bool

	exists, _ = da.GroupExists(ctx, "test-exists")
	if exists {
		t.Error("Group should not exist now")
		t.FailNow()
	}

	// Now we add a group to find.
	da.GroupCreate(ctx, rest.Group{Name: "test-exists"})
	defer da.GroupDelete(ctx, "test-exists")

	exists, _ = da.GroupExists(ctx, "test-exists")
	if !exists {
		t.Error("Group should exist now")
		t.FailNow()
	}
}

func testGroupGet(t *testing.T) {
	var err error
	var group rest.Group

	// Expect an error
	_, err = da.GroupGet(ctx, "")
	assert.Error(t, err, errs.ErrEmptyGroupName)

	// Expect an error
	_, err = da.GroupGet(ctx, "test-get")
	assert.Error(t, err, errs.ErrNoSuchGroup)

	da.GroupCreate(ctx, rest.Group{Name: "test-get"})
	defer da.GroupDelete(ctx, "test-get")

	// da.Group ctx, should exist now
	exists, _ := da.GroupExists(ctx, "test-get")
	if !exists {
		t.Error("Group should exist now")
		t.FailNow()
	}

	// Expect no error
	group, err = da.GroupGet(ctx, "test-get")
	assert.NoError(t, err)
	if group.Name != "test-get" {
		t.Errorf("Group name mismatch: %q is not \"test-get\"", group.Name)
		t.FailNow()
	}
}

func testGroupGrantRole(t *testing.T) {
	var err error

	groupName := "group-group-grant-role"
	roleName := "role-group-grant-role"
	bundleName := "bundle-group-grant-role"
	permissionName := "perm-group-grant-role"

	da.GroupCreate(ctx, rest.Group{Name: groupName})
	defer da.GroupDelete(ctx, groupName)

	err = da.RoleCreate(ctx, roleName)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	err = da.RoleGrantPermission(ctx, roleName, bundleName, permissionName)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	err = da.GroupGrantRole(ctx, groupName, roleName)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	expectedRoles := []rest.Role{
		{
			Name:        roleName,
			Permissions: []rest.RolePermission{{BundleName: bundleName, Permission: permissionName}},
		},
	}

	roles, err := da.GroupListRoles(ctx, groupName)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Equal(t, expectedRoles, roles)

	err = da.GroupRevokeRole(ctx, groupName, roleName)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	expectedRoles = []rest.Role{}

	roles, err = da.GroupListRoles(ctx, groupName)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Equal(t, expectedRoles, roles)
}

func testGroupList(t *testing.T) {
	da.GroupCreate(ctx, rest.Group{Name: "test-list-0"})
	defer da.GroupDelete(ctx, "test-list-0")
	da.GroupCreate(ctx, rest.Group{Name: "test-list-1"})
	defer da.GroupDelete(ctx, "test-list-1")
	da.GroupCreate(ctx, rest.Group{Name: "test-list-2"})
	defer da.GroupDelete(ctx, "test-list-2")
	da.GroupCreate(ctx, rest.Group{Name: "test-list-3"})
	defer da.GroupDelete(ctx, "test-list-3")

	groups, err := da.GroupList(ctx)
	assert.NoError(t, err)

	if len(groups) != 4 {
		t.Errorf("Expected len(groups) = 4; got %d", len(groups))
		t.FailNow()
	}

	for _, u := range groups {
		if u.Name == "" {
			t.Error("Expected non-empty name")
			t.FailNow()
		}
	}
}

func testGroupListRoles(t *testing.T) {
	da.GroupCreate(ctx, rest.Group{Name: "group-test-group-list-roles"})
	defer da.GroupDelete(ctx, "group-test-group-list-roles")

	da.RoleCreate(ctx, "role-test-group-list-roles-1")
	defer da.RoleDelete(ctx, "role-test-group-list-roles-1")

	da.RoleCreate(ctx, "role-test-group-list-roles-0")
	defer da.RoleDelete(ctx, "role-test-group-list-roles-0")

	da.RoleCreate(ctx, "role-test-group-list-roles-2")
	defer da.RoleDelete(ctx, "role-test-group-list-roles-2")

	roles, err := da.GroupListRoles(ctx, "group-test-group-list-roles")
	if !assert.NoError(t, err) && !assert.Empty(t, roles) {
		t.FailNow()
	}

	err = da.GroupGrantRole(ctx, "group-test-group-list-roles", "role-test-group-list-roles-1")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	err = da.GroupGrantRole(ctx, "group-test-group-list-roles", "role-test-group-list-roles-0")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	err = da.GroupGrantRole(ctx, "group-test-group-list-roles", "role-test-group-list-roles-2")
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	// Note: alphabetically sorted!
	expected := []rest.Role{
		{Name: "role-test-group-list-roles-0", Permissions: []rest.RolePermission{}},
		{Name: "role-test-group-list-roles-1", Permissions: []rest.RolePermission{}},
		{Name: "role-test-group-list-roles-2", Permissions: []rest.RolePermission{}},
	}

	actual, err := da.GroupListRoles(ctx, "group-test-group-list-roles")
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Equal(t, expected, actual)
}

func testGroupRemoveUser(t *testing.T) {
	da.GroupCreate(ctx, rest.Group{Name: "foo"})
	defer da.GroupDelete(ctx, "foo")

	da.UserCreate(ctx, rest.User{Username: "bat"})
	defer da.UserDelete(ctx, "bat")

	err := da.GroupAddUser(ctx, "foo", "bat")
	assert.NoError(t, err)

	group, err := da.GroupGet(ctx, "foo")
	assert.NoError(t, err)

	if len(group.Users) != 1 {
		t.Error("Users list empty")
		t.FailNow()
	}

	if len(group.Users) > 0 && group.Users[0].Username != "bat" {
		t.Error("Wrong user!")
		t.FailNow()
	}

	err = da.GroupRemoveUser(ctx, "foo", "bat")
	assert.NoError(t, err)

	group, err = da.GroupGet(ctx, "foo")
	assert.NoError(t, err)

	if len(group.Users) != 0 {
		t.Error("User not removed")
		t.FailNow()
	}
}
