// Copyright © 2018 Banzai Cloud
//
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

package oracle

import (
	"fmt"

	"github.com/banzaicloud/cloudinfo/pkg/cloudinfo"
)

var (
	ntwPerfMap = map[string][]string{
		cloudinfo.NTW_LOW:    {"0.6 Gbps"},
		cloudinfo.NTW_MEDIUM: {"1 Gbps", "1.2 Gbps", "2 Gbps", "2.4 Gbps"},
		cloudinfo.NTW_HIGH:   {"4.1 Gbps", "4.8 Gbps", "8.2 Gbps"},
		cloudinfo.NTW_EXTRA:  {"16.4 Gbps", "24.6 Gbps"},
	}
)

// OCINetworkMapper module object for handling Oracle specific VM to Networking capabilities mapping
type OCINetworkMapper struct {
}

// newNetworkMapper initializes the network performance mapper struct
func newNetworkMapper() *OCINetworkMapper {
	return &OCINetworkMapper{}
}

// MapNetworkPerf maps the network performance of the instance to the category supported by telescopes
func (nm *OCINetworkMapper) MapNetworkPerf(ntwPerf string) (string, error) {
	for perfCat, strVals := range ntwPerfMap {
		if cloudinfo.Contains(strVals, ntwPerf) {
			return perfCat, nil
		}
	}
	return "", fmt.Errorf("could not determine network performance for: [%s]", ntwPerf)
}
