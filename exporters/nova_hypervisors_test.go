package exporters

import (
	"context"
	"net/http"
	"strings"

	"github.com/jarcoal/httpmock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

// Hypervisor list from a cloud mixing libvirt, nova-incus and ironic computes.
// The nova-incus driver reports cpu_info.features as a space separated string
// and the topology counts as strings, exactly as captured from a real cloud.
const mixedHypervisorsResponse = `{"hypervisors": [
  {"id": "1", "hypervisor_hostname": "libvirt-host", "hypervisor_type": "QEMU", "hypervisor_version": 9002000,
   "state": "up", "status": "enabled", "vcpus": 8, "vcpus_used": 2, "memory_mb": 16000, "memory_mb_used": 4000,
   "service": {"host": "libvirt-host", "id": "a"},
   "cpu_info": {"arch": "x86_64", "features": ["pge", "clflush"], "topology": {"cells": 1, "sockets": 1, "cores": 4, "threads": 2}}},
  {"id": "2", "hypervisor_hostname": "incus-host", "hypervisor_type": "lxd", "hypervisor_version": 6000000,
   "state": "up", "status": "enabled", "vcpus": 32, "vcpus_used": 4, "memory_mb": 64000, "memory_mb_used": 8000,
   "service": {"host": "incus-host", "id": "b"},
   "cpu_info": {"arch": "x86_64", "vendor": "AuthenticAMD", "model": "EPYC", "features": "fpu vme de pse", "topology": {"sockets": "2", "cores": "8", "threads": "2"}}},
  {"id": "3", "hypervisor_hostname": "ironic-node", "hypervisor_type": "ironic", "hypervisor_version": 1,
   "state": "up", "status": "enabled", "vcpus": 0, "vcpus_used": 0, "memory_mb": 0, "memory_mb_used": 0,
   "service": {"host": "ironic", "id": "c"},
   "cpu_info": {}}
]}`

// The request goes through the real hypervisors pager: gophercloud calls
// ExtractHypervisors from HypervisorPage.IsEmpty while paging, so a fix that
// only touches the final extraction is not enough.
func (suite *NovaTestSuite) TestHypervisorsWithStringCPUFeatures() {
	// The pager only decodes a JSON body and reads the page URL from the
	// response request, so both must be set.
	httpmock.RegisterResponder("GET", suite.MakeURL("/compute/os-hypervisors/detail", ""),
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, mixedHypervisorsResponse)
			resp.Header.Set("Content-Type", "application/json")
			resp.Request = req
			return resp, nil
		})

	exporter := (*suite.Exporter).(*NovaExporter)
	ch := make(chan prometheus.Metric, 1000)
	err := ListHypervisors(context.Background(), &exporter.BaseOpenStackExporter, ch)
	close(ch)
	assert.NoError(suite.T(), err)

	workloads := 0
	for m := range ch {
		if strings.Contains(m.Desc().String(), "openstack_nova_current_workload") {
			workloads++
		}
	}
	assert.Equal(suite.T(), 3, workloads)
}
