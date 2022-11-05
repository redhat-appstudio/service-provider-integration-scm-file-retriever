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
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"net/http"
)

var (
	invalidSourceError = errors.New("invalid source string")
)

type DetectError struct {
	Error      error `json:"error"`
	StatusCode int   `json:"status"`
}

// ScmProvider defines the interface that in order to determine if URL belongs to SCM provider
type ScmProvider interface {
	// detect will check whether the provided repository URL matches a known SCM pattern,
	// and transform input params into valid file download URL.
	// Params are repository, path to the file inside the repository, Git reference (branch/tag/commitId),
	// HTTP client instance and authentication headers holder struct
	detect(ctx context.Context, repoUrl, filepath, ref string, client *req.Client, auth HeaderStruct) (bool, string, DetectError)
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

func detect(ctx context.Context, repoUrl, filepath, ref string, client *req.Client, auth HeaderStruct) (string, DetectError) {
	for _, d := range ScmProviders {
		ok, resultUrl, detectError := d.detect(ctx, repoUrl, filepath, ref, client, auth)
		if !ok {
			continue
		}
		if detectError.Error != nil {
			return "", detectError
		}
		return resultUrl, detectError
	}
	return "", DetectError{StatusCode: http.StatusInternalServerError, Error: fmt.Errorf("%w: %s for %s", invalidSourceError, repoUrl, filepath)}
}
