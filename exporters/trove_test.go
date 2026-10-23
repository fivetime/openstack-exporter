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
# HELP openstack_trove_instance_operating_status instance_operating_status
# TYPE openstack_trove_instance_operating_status gauge
openstack_trove_instance_operating_status{cluster_id="38f82ead-a7b6-46b5-ac47-fcff8f4a6909",datastore_type="cassandra",datastore_version="5.0",health_status="",id="9413e539-2138-4df6-b524-a064b49b13a6",name="cs2-member-dc1-rack1-2",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 3
openstack_trove_instance_operating_status{cluster_id="38f82ead-a7b6-46b5-ac47-fcff8f4a6909",datastore_type="cassandra",datastore_version="5.0",health_status="",id="a3c83da9-3d6c-43e0-bff5-addf1e100e37",name="cs2-member-dc1-rack1-3",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 3
openstack_trove_instance_operating_status{cluster_id="38f82ead-a7b6-46b5-ac47-fcff8f4a6909",datastore_type="cassandra",datastore_version="5.0",health_status="",id="b7d29e9e-8c62-4c59-8acd-4a00af9013f6",name="cs2-member-dc1-rack1-1",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 3
# HELP openstack_trove_instance_status instance_status
# TYPE openstack_trove_instance_status gauge
openstack_trove_instance_status{cluster_id="38f82ead-a7b6-46b5-ac47-fcff8f4a6909",datastore_type="cassandra",datastore_version="5.0",health_status="",id="9413e539-2138-4df6-b524-a064b49b13a6",name="cs2-member-dc1-rack1-2",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 2
openstack_trove_instance_status{cluster_id="38f82ead-a7b6-46b5-ac47-fcff8f4a6909",datastore_type="cassandra",datastore_version="5.0",health_status="",id="a3c83da9-3d6c-43e0-bff5-addf1e100e37",name="cs2-member-dc1-rack1-3",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 2
openstack_trove_instance_status{cluster_id="38f82ead-a7b6-46b5-ac47-fcff8f4a6909",datastore_type="cassandra",datastore_version="5.0",health_status="",id="b7d29e9e-8c62-4c59-8acd-4a00af9013f6",name="cs2-member-dc1-rack1-1",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 2
# HELP openstack_trove_instance_volume_size_gb instance_volume_size_gb
# TYPE openstack_trove_instance_volume_size_gb gauge
openstack_trove_instance_volume_size_gb{cluster_id="38f82ead-a7b6-46b5-ac47-fcff8f4a6909",datastore_type="cassandra",datastore_version="5.0",health_status="",id="9413e539-2138-4df6-b524-a064b49b13a6",name="cs2-member-dc1-rack1-2",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 4
openstack_trove_instance_volume_size_gb{cluster_id="38f82ead-a7b6-46b5-ac47-fcff8f4a6909",datastore_type="cassandra",datastore_version="5.0",health_status="",id="a3c83da9-3d6c-43e0-bff5-addf1e100e37",name="cs2-member-dc1-rack1-3",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 4
openstack_trove_instance_volume_size_gb{cluster_id="38f82ead-a7b6-46b5-ac47-fcff8f4a6909",datastore_type="cassandra",datastore_version="5.0",health_status="",id="b7d29e9e-8c62-4c59-8acd-4a00af9013f6",name="cs2-member-dc1-rack1-1",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 4
# HELP openstack_trove_instance_volume_used_gb instance_volume_used_gb
# TYPE openstack_trove_instance_volume_used_gb gauge
openstack_trove_instance_volume_used_gb{cluster_id="38f82ead-a7b6-46b5-ac47-fcff8f4a6909",datastore_type="cassandra",datastore_version="5.0",health_status="",id="9413e539-2138-4df6-b524-a064b49b13a6",name="cs2-member-dc1-rack1-2",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 0
openstack_trove_instance_volume_used_gb{cluster_id="38f82ead-a7b6-46b5-ac47-fcff8f4a6909",datastore_type="cassandra",datastore_version="5.0",health_status="",id="a3c83da9-3d6c-43e0-bff5-addf1e100e37",name="cs2-member-dc1-rack1-3",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 0
openstack_trove_instance_volume_used_gb{cluster_id="38f82ead-a7b6-46b5-ac47-fcff8f4a6909",datastore_type="cassandra",datastore_version="5.0",health_status="",id="b7d29e9e-8c62-4c59-8acd-4a00af9013f6",name="cs2-member-dc1-rack1-1",operating_status="HEALTHY",region="",status="ACTIVE",tenant_id="ac19c812536b4b06bd8c394353d683f3"} 0
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
