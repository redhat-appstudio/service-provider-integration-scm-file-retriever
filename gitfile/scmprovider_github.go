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
	"regexp"

	"github.com/imroc/req"
	"go.uber.org/zap"
)

type GithubFile struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        int32  `json:"size"`
	Encoding    string `json:"encoding"`
	DownloadUrl string `json:"download_url"`
}

var GithubAPITemplate = "https://api.github.com/repos/%s/%s/contents/%s"
var GithubURLRegexp = regexp.MustCompile(`(?Um)^(?:https)(?:\:\/\/)github.com/(?P<repoUser>[^/]+)/(?P<repoName>[^/]+)(.git)?$`)
var GithubURLRegexpNames = GithubURLRegexp.SubexpNames()

// GitHubScmProvider implements Detector to detect GitHub URLs.
type GitHubScmProvider struct {
}

func (d *GitHubScmProvider) Detect(repoUrl, filepath, ref string) (bool, string, error) {
	if len(repoUrl) == 0 || !GithubURLRegexp.MatchString(repoUrl) {
		return false, "", nil
	}

	result := GithubURLRegexp.FindAllStringSubmatch(repoUrl, -1)
	m := map[string]string{}
	for i, n := range result[0] {
		m[GithubURLRegexpNames[i]] = n
	}
	param := req.Param{}
	if ref != "" {
		param["ref"] = ref
	}
	r, e := req.Get(fmt.Sprintf(GithubAPITemplate, m["repoUser"], m["repoName"], filepath), param)
	if e != nil {
		zap.L().Error("Failed to make GitHub API call", zap.Error(e))
		return true, "", e
	}
	statusCode := r.Response().StatusCode
	zap.L().Debug(fmt.Sprintf(
		"GitHub API call response code: %d", statusCode))
	if statusCode >= 400 {
		return true, "", fmt.Errorf("unexpected status code from GitHub API: %d. Response: %s", statusCode, r.String())
	}

	var file GithubFile
	err := r.ToJSON(&file)
	if err != nil {
		zap.L().Error("Failed to parse GitHub json response", zap.Error(err))
		return true, "", err
	}
	return true, file.DownloadUrl, nil
}
