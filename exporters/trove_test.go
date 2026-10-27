package exporters

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

type TroveTestSuite struct {
	BaseOpenStackTestSuite
}

var troveExpectedUp = `
# HELP openstack_trove_cluster_health cluster_health
# TYPE openstack_trove_cluster_health gauge
openstack_trove_cluster_health{datastore_type="valkey",datastore_version="9-kata",health_status="HEALTHY",id="73dfda58-c11a-499c-be16-370e30e20955",name="kw5-valkey-s",task="NONE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 0
# HELP openstack_trove_cluster_members_out cluster_members_out
# TYPE openstack_trove_cluster_members_out gauge
openstack_trove_cluster_members_out{datastore_type="valkey",datastore_version="9-kata",health_status="HEALTHY",id="73dfda58-c11a-499c-be16-370e30e20955",name="kw5-valkey-s",task="NONE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 0
# HELP openstack_trove_cluster_writer_changes cluster_writer_changes
# TYPE openstack_trove_cluster_writer_changes counter
openstack_trove_cluster_writer_changes{datastore_type="valkey",datastore_version="9-kata",id="73dfda58-c11a-499c-be16-370e30e20955",name="kw5-valkey-s",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 1
# HELP openstack_trove_instance_operating_status instance_operating_status
# TYPE openstack_trove_instance_operating_status gauge
openstack_trove_instance_operating_status{cluster_id="73dfda58-c11a-499c-be16-370e30e20955",datastore_type="valkey",datastore_version="9-kata",health_status="",id="04585468-d5b1-4c78-a7a4-138403bb8443",name="kw5-valkey-s-member-1793107093",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 3
openstack_trove_instance_operating_status{cluster_id="73dfda58-c11a-499c-be16-370e30e20955",datastore_type="valkey",datastore_version="9-kata",health_status="",id="16517140-bd04-4d87-ab75-09ca0c975367",name="kw5-valkey-s-member-1793107092",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 3
openstack_trove_instance_operating_status{cluster_id="73dfda58-c11a-499c-be16-370e30e20955",datastore_type="valkey",datastore_version="9-kata",health_status="",id="1915b943-f022-46a1-ac5a-57a1199a95e5",name="kw5-valkey-s-member-1793107091",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 3
# HELP openstack_trove_instance_status instance_status
# TYPE openstack_trove_instance_status gauge
openstack_trove_instance_status{cluster_id="73dfda58-c11a-499c-be16-370e30e20955",datastore_type="valkey",datastore_version="9-kata",health_status="",id="04585468-d5b1-4c78-a7a4-138403bb8443",name="kw5-valkey-s-member-1793107093",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 2
openstack_trove_instance_status{cluster_id="73dfda58-c11a-499c-be16-370e30e20955",datastore_type="valkey",datastore_version="9-kata",health_status="",id="16517140-bd04-4d87-ab75-09ca0c975367",name="kw5-valkey-s-member-1793107092",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 2
openstack_trove_instance_status{cluster_id="73dfda58-c11a-499c-be16-370e30e20955",datastore_type="valkey",datastore_version="9-kata",health_status="",id="1915b943-f022-46a1-ac5a-57a1199a95e5",name="kw5-valkey-s-member-1793107091",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 2
# HELP openstack_trove_instance_volume_size_gb instance_volume_size_gb
# TYPE openstack_trove_instance_volume_size_gb gauge
openstack_trove_instance_volume_size_gb{cluster_id="73dfda58-c11a-499c-be16-370e30e20955",datastore_type="valkey",datastore_version="9-kata",health_status="",id="04585468-d5b1-4c78-a7a4-138403bb8443",name="kw5-valkey-s-member-1793107093",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 2
openstack_trove_instance_volume_size_gb{cluster_id="73dfda58-c11a-499c-be16-370e30e20955",datastore_type="valkey",datastore_version="9-kata",health_status="",id="16517140-bd04-4d87-ab75-09ca0c975367",name="kw5-valkey-s-member-1793107092",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 2
openstack_trove_instance_volume_size_gb{cluster_id="73dfda58-c11a-499c-be16-370e30e20955",datastore_type="valkey",datastore_version="9-kata",health_status="",id="1915b943-f022-46a1-ac5a-57a1199a95e5",name="kw5-valkey-s-member-1793107091",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 2
# HELP openstack_trove_instance_volume_used_gb instance_volume_used_gb
# TYPE openstack_trove_instance_volume_used_gb gauge
openstack_trove_instance_volume_used_gb{cluster_id="73dfda58-c11a-499c-be16-370e30e20955",datastore_type="valkey",datastore_version="9-kata",health_status="",id="04585468-d5b1-4c78-a7a4-138403bb8443",name="kw5-valkey-s-member-1793107093",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 0
openstack_trove_instance_volume_used_gb{cluster_id="73dfda58-c11a-499c-be16-370e30e20955",datastore_type="valkey",datastore_version="9-kata",health_status="",id="16517140-bd04-4d87-ab75-09ca0c975367",name="kw5-valkey-s-member-1793107092",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 0
openstack_trove_instance_volume_used_gb{cluster_id="73dfda58-c11a-499c-be16-370e30e20955",datastore_type="valkey",datastore_version="9-kata",health_status="",id="1915b943-f022-46a1-ac5a-57a1199a95e5",name="kw5-valkey-s-member-1793107091",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 0
# HELP openstack_trove_total_clusters total_clusters
# TYPE openstack_trove_total_clusters gauge
openstack_trove_total_clusters 1
# HELP openstack_trove_total_instances total_instances
# TYPE openstack_trove_total_instances gauge
openstack_trove_total_instances 3
# HELP openstack_trove_up up
# TYPE openstack_trove_up gauge
openstack_trove_up 1
`

func (suite *TroveTestSuite) TestTroveExporter() {
	err := testutil.CollectAndCompare(*suite.Exporter, strings.NewReader(troveExpectedUp))
	assert.NoError(suite.T(), err)
}
