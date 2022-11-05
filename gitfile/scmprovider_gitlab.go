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
	"context"
	"github.com/imroc/req/v3"
	"strings"
)

// GitLabScmProvider implements Detector to detect Gitlab URLs.
type GitLabScmProvider struct{}

func (d *GitLabScmProvider) detect(ctx context.Context, repoUrl, filepath, ref string, cl *req.Client, auth HeaderStruct) (bool, string, DetectError) {
	if len(repoUrl) == 0 {
		return false, "", DetectError{}
	}

	if strings.HasPrefix(repoUrl, "https://gitlab.com/") {
		return true, "", DetectError{}
	}

	return false, "", DetectError{}
}
