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

package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/getgort/gort/data/rest"
)

// GroupDelete comments to be written...
func (c *GortClient) GroupDelete(groupname string) error {
	url := fmt.Sprintf("%s/v2/groups/%s", c.profile.URL.String(), groupname)

	resp, err := c.doRequest("DELETE", url, []byte{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return getResponseError(resp)
	}

	return nil
}

// GroupExists simply returns true if a group exists with the specified
// groupname; false otherwise.
func (c *GortClient) GroupExists(groupname string) (bool, error) {
	url := fmt.Sprintf("%s/v2/groups/%s", c.profile.URL.String(), groupname)
	resp, err := c.doRequest("GET", url, []byte{})
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, getResponseError(resp)
	}
}

// GroupGet comments to be written...
func (c *GortClient) GroupGet(groupname string) (rest.Group, error) {
	url := fmt.Sprintf("%s/v2/groups/%s", c.profile.URL.String(), groupname)
	resp, err := c.doRequest("GET", url, []byte{})
	if err != nil {
		return rest.Group{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return rest.Group{}, getResponseError(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rest.Group{}, err
	}

	group := rest.Group{}
	err = json.Unmarshal(body, &group)
	if err != nil {
		return rest.Group{}, err
	}

	return group, nil
}

// GroupList comments to be written...
func (c *GortClient) GroupList() ([]rest.Group, error) {
	url := fmt.Sprintf("%s/v2/groups", c.profile.URL.String())
	resp, err := c.doRequest("GET", url, []byte{})
	if err != nil {
		return []rest.Group{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []rest.Group{}, getResponseError(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []rest.Group{}, err
	}

	groups := []rest.Group{}
	err = json.Unmarshal(body, &groups)
	if err != nil {
		return []rest.Group{}, err
	}

	return groups, nil
}

// GroupMemberAdd comments to be written...
func (c *GortClient) GroupMemberAdd(groupname string, username string) error {
	url := fmt.Sprintf("%s/v2/groups/%s/members/%s", c.profile.URL.String(), groupname, username)
	resp, err := c.doRequest("PUT", url, []byte{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return getResponseError(resp)
	}

	return nil
}

// GroupMemberDelete comments to be written...
func (c *GortClient) GroupMemberDelete(groupname string, username string) error {
	url := fmt.Sprintf("%s/v2/groups/%s/members/%s", c.profile.URL.String(), groupname, username)
	resp, err := c.doRequest("DELETE", url, []byte{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return getResponseError(resp)
	}

	return nil
}

// GroupMemberList comments to be written...
func (c *GortClient) GroupMemberList(groupname string) ([]rest.User, error) {
	url := fmt.Sprintf("%s/v2/groups/%s/members", c.profile.URL.String(), groupname)
	resp, err := c.doRequest("GET", url, []byte{})
	if err != nil {
		return []rest.User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []rest.User{}, getResponseError(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []rest.User{}, err
	}

	users := []rest.User{}
	err = json.Unmarshal(body, &users)
	if err != nil {
		return []rest.User{}, err
	}

	return users, nil
}

// GroupSave comments to be written...
func (c *GortClient) GroupSave(group rest.Group) error {
	url := fmt.Sprintf("%s/v2/groups/%s", c.profile.URL.String(), group.Name)

	bytes, err := json.Marshal(group)
	if err != nil {
		return err
	}

	resp, err := c.doRequest("PUT", url, bytes)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return getResponseError(resp)
	}

	return nil
}

// GroupRoleAdd adds a role to a group.
func (c *GortClient) GroupRoleAdd(groupname string, rolename string) error {
	url := fmt.Sprintf("%s/v2/groups/%s/roles/%s", c.profile.URL.String(), groupname, rolename)
	resp, err := c.doRequest("PUT", url, []byte{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return getResponseError(resp)
	}

	return nil
}

// GroupMemberDelete deletes a role from a group.
func (c *GortClient) GroupRoleDelete(groupname string, rolename string) error {
	url := fmt.Sprintf("%s/v2/groups/%s/roles/%s", c.profile.URL.String(), groupname, rolename)
	resp, err := c.doRequest("DELETE", url, []byte{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return getResponseError(resp)
	}

	return nil
}

// GroupRoleList retrieves all roles added to a group.
func (c *GortClient) GroupRoleList(groupname string) ([]rest.Role, error) {
	url := fmt.Sprintf("%s/v2/groups/%s/roles", c.profile.URL.String(), groupname)
	resp, err := c.doRequest("GET", url, []byte{})
	if err != nil {
		return []rest.Role{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []rest.Role{}, getResponseError(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []rest.Role{}, err
	}

	roles := []rest.Role{}
	err = json.Unmarshal(body, &roles)
	if err != nil {
		return []rest.Role{}, err
	}

	return roles, nil
}
