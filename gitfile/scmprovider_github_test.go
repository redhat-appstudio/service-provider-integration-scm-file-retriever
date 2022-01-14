// Copyright (c) 2022 Red Hat, Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitfile

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileHead(t *testing.T) {
	r1, err := Detect("https://github.com/redhat-appstudio/service-provider-integration-operator", "Makefile", "HEAD")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	assert.Equal(t, r1, "https://raw.githubusercontent.com/redhat-appstudio/service-provider-integration-operator/HEAD/Makefile")
}

func TestGetFileHead2(t *testing.T) {
	r1, err := Detect("https://github.com/redhat-appstudio/service-provider-integration-operator.git", "Makefile", "HEAD")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	assert.Equal(t, r1, "https://raw.githubusercontent.com/redhat-appstudio/service-provider-integration-operator/HEAD/Makefile")
}

func TestGetFileOnBranch(t *testing.T) {
	r1, err := Detect("https://github.com/redhat-appstudio/service-provider-integration-operator.git", "Makefile", "v0.1.0")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	assert.Equal(t, r1, "https://raw.githubusercontent.com/redhat-appstudio/service-provider-integration-operator/v0.1.0/Makefile")
}

func TestGetFileOnCommitId(t *testing.T) {
	r1, err := Detect("https://github.com/redhat-appstudio/service-provider-integration-operator.git", "Makefile", "efaf08a367921ae130c524db4a531b7696b7d967")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	assert.Equal(t, r1, "https://raw.githubusercontent.com/redhat-appstudio/service-provider-integration-operator/efaf08a367921ae130c524db4a531b7696b7d967/Makefile")
}

func TestGetUnexistingFile(t *testing.T) {
	_, err := Detect("https://github.com/redhat-appstudio/service-provider-integration-operator.git", "Makefile-Non-Exist", "")

	if err == nil {
		t.Error("error expected")
	}
	assert.Equal(t, "unexpected status code from GitHub API: 404. Response: {\"message\":\"Not Found\",\"documentation_url\":\"https://docs.github.com/rest/reference/repos#get-repository-content\"}", fmt.Sprint(err))
}
