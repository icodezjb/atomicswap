// Copyright 2019 icodezjb
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cmd

import (
	"fmt"
	"runtime"
)

var (
	Commit    = "unknown-commit"
	BuildTime = "unknown-buildtime"

	Version = "0.2.0"
)

// VersionFunc holds the textual version string.
func VersionFunc() string {
	return fmt.Sprintf(": %s\ncommit: %s\nbuild time: %s\ngolang version: %s\n",
		Version, Commit, BuildTime, runtime.Version()+" "+runtime.GOOS+"/"+runtime.GOARCH)
}
