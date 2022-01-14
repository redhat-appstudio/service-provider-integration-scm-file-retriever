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

import "fmt"

// ScmProvider defines the interface that in order to determine if URL belongs to SCM provider
type ScmProvider interface {
	// Detect will detect whether the string matches a known SCM file pattern
	// and transform it to valid https url.
	Detect(repoUrl, filepath, ref string) (bool, string, error)
}

// ScmProviders is the list of detectors that are tried on an SCM URL.
// This is also the order they're tried (index 0 is first).
var ScmProviders []ScmProvider

func init() {
	ScmProviders = []ScmProvider{
		new(GitLabScmProvider),
		new(GitHubScmProvider),
	}
}

func Detect(repoUrl, filepath, ref string) (string, error) {
	for _, d := range ScmProviders {
		ok, resultUrl, err := d.Detect(repoUrl, filepath, ref)
		if err != nil {
			return "", err
		}
		if !ok {
			continue
		}
		return resultUrl, nil
	}
	return "", fmt.Errorf("invalid source string: %s for %s", repoUrl, filepath)
}
