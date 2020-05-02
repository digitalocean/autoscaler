/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package digitalocean

import (
	"context"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/digitalocean/godo"
)

type capacity struct {
	cpus   int
	memory int
}

func (m *Manager) ensureCapacityMap(ctx context.Context) error {
	if len(m.capacityByDropletSlug) > 0 {
		return nil
	}

	var sizes []godo.Size
	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 100,
	}
	for {
		sizes, resp, err := m.sizeLister.List(ctx, opts)
		if err != nil {
			return err
		}

		sizes = append(sizes, sizes...)

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return err
		}

		opts.Page = page + 1
	}

	for _, size := range sizes {
		m.capacityByDropletSlug[size.Slug] = capacity{
			cpus:   size.Vcpus,
			memory: size.Memory,
		}
	}

	return nil
}
