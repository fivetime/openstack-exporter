package exporters

import (
	"context"
	"log/slog"
	"net/url"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/datastores"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"github.com/prometheus/client_golang/prometheus"
)

var knownDBInstanceStatuses = map[string]int{
	"NEW":              0,  // The database instance creation request is just received by Trove.
	"BUILD":            1,  // The database instance is being installed.
	"ACTIVE":           2,  // The database instance is up and running.
	"REBOOT":           3,  // The database instance is rebooting.
	"RESIZE":           4,  // The database instance is being resized.
	"UPGRADE":          5,  // The database instance is upgrading its datastore, e.g. from mysql 5.7.29 to mysql 5.7.30
	"RESTART_REQUIRED": 6,  // The database service needs to restart, e.g. due to the configuration change.
	"PROMOTE":          7,  // A replica instance in the replication cluster is being promoted to the primary.
	"EJECT":            8,  // The current primary instance in a replication cluster is being ejected, one of the replicas is going to be elected as the new primary.
	"DETACH":           9,  // One of the replicas in a replication cluster is being detached and will become a standalone instance.
	"SHUTDOWN":         10, // The database instance is being shutdown during deletion.
	"BACKUP":           11, // The database instance is being backed up.
	"ERROR":            12, // The database instance is in error.
}

// The operating status is the state of the database service inside the
// instance, as its guest agent last reported it (Trove's ServiceStatuses),
// distinct from the status of the instance itself.
var knownDBOperatingStatuses = map[string]int{
	"NEW":              0,  // The guest agent has not reported yet.
	"BUILD":            1,  // The database is being installed.
	"ACTIVE":           2,  // The database is running.
	"HEALTHY":          3,  // The database is running and answers.
	"BLOCKED":          4,  // The database is running but does not answer.
	"SHUTDOWN":         5,  // The database is stopped, paused or has crashed.
	"FAILED":           6,  // The database failed to start.
	"PROMOTE":          7,  // The replica is being promoted.
	"EJECT":            8,  // The replication source is being ejected.
	"LOGGING":          9,  // Guest logs are being transferred.
	"RESTART_REQUIRED": 10, // The database needs a restart, e.g. for a configuration change.
	"UPGRADING":        11, // The database is being upgraded.
	"RESTARTING":       12, // The database is restarting.
	"DELETED":          13, // The database is deleted.
	"ERROR":            14, // The guest agent is in error or has not reported for too long.
}

// The health of a cluster as a whole, as the task manager recorded it.
var knownDBClusterHealthStatuses = map[string]int{
	"HEALTHY":  0, // Every member answers.
	"DEGRADED": 1, // Some member is out; the cluster still serves.
	"DOWN":     2, // The cluster does not serve.
	"UNKNOWN":  3, // The task manager has not looked yet.
}

func mapDBInstanceStatus(current string) int {
	return mapStatus(knownDBInstanceStatuses, current)
}

func mapDBOperatingStatus(current string) int {
	return mapStatus(knownDBOperatingStatuses, current)
}

func mapDBClusterHealthStatus(current string) int {
	return mapStatus(knownDBClusterHealthStatuses, current)
}

type TroveExporter struct {
	BaseOpenStackExporter
}

type instanceAttributesExt struct {
	// Indicates the unique identifier for the instance resource.
	ID string

	// The human-readable name of the instance.
	Name string

	// The build status of the instance.
	Status string

	// The status of the database service inside the instance.
	OperatingStatus string `json:"operating_status"`

	// Information about the attached volume of the instance.
	Volume instances.Volume

	// Indicates how the instance stores data.
	Datastore datastores.DatastorePartial

	// Region is the region the instance is located in.
	Region string

	HealthStatus string `json:"health_status"`

	// TenantID is the id of the project that owns the instance.
	TenantID string `json:"tenant_id"`

	// ClusterID is the id of the cluster the instance is a member of; empty
	// for a standalone instance.
	ClusterID string `json:"cluster_id"`
}

type clusterAttributes struct {
	// Indicates the unique identifier for the cluster resource.
	ID string

	// The human-readable name of the cluster.
	Name string

	// The task the cluster is on; NONE when idle.
	Task struct {
		Name string `json:"name"`
	} `json:"task"`

	// Indicates how the cluster stores data.
	Datastore datastores.DatastorePartial

	// The health of the cluster as the task manager recorded it.
	Health struct {
		Status     string   `json:"status"`
		MembersOut []string `json:"members_out"`
	} `json:"health"`

	// TenantID is the id of the project that owns the cluster.
	TenantID string `json:"tenant_id"`
}

// troveInstancePage is a page of the management instance list. Trove links
// the next page as "links" (rel "next"), not "instances_links", and the
// link is built from the request as the API saw it, so only its query is
// trusted: the next page is asked for at the URL the first one was.
type troveInstancePage struct {
	pagination.LinkedPageBase
	listURL string
}

func (page troveInstancePage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	next, err := gophercloud.ExtractNextURL(s.Links)
	if err != nil || next == "" {
		return next, err
	}
	nextURL, err := url.Parse(next)
	if err != nil {
		return "", err
	}
	listURL, err := url.Parse(page.listURL)
	if err != nil {
		return "", err
	}
	listURL.RawQuery = nextURL.RawQuery
	return listURL.String(), nil
}

func (page troveInstancePage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}
	instances, err := extractDBInstances(page)
	return len(instances) == 0, err
}

// troveClusterPage is the management cluster list, which comes in one page.
type troveClusterPage struct {
	pagination.SinglePageBase
}

func (page troveClusterPage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}
	clusters, err := extractDBClusters(page)
	return len(clusters) == 0, err
}

// extractInstances will convert a generic pagination struct into a more
// relevant slice of instanceAttributesExt structs.
func extractDBInstances(r pagination.Page) ([]instanceAttributesExt, error) {
	var s struct {
		Instances []instanceAttributesExt `json:"instances"`
	}
	err := (r.(troveInstancePage)).ExtractInto(&s)
	return s.Instances, err
}

// extractDBClusters will convert a generic pagination struct into a slice
// of clusterAttributes structs.
func extractDBClusters(r pagination.Page) ([]clusterAttributes, error) {
	var s struct {
		Clusters []clusterAttributes `json:"clusters"`
	}
	err := (r.(troveClusterPage)).ExtractInto(&s)
	return s.Clusters, err
}

// list retrieves the status and information for all database instances,
// the members of clusters among them.
func listDBInstances(client *gophercloud.ServiceClient) pagination.Pager {
	listURL := client.ServiceURL("mgmt", "instances?include_clustered=True&deleted=False")
	return pagination.NewPager(client, listURL, func(r pagination.PageResult) pagination.Page {
		return troveInstancePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}, listURL: listURL}
	})
}

// list retrieves the status and health of all database clusters.
func listDBClusters(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, client.ServiceURL("mgmt", "clusters?deleted=False"), func(r pagination.PageResult) pagination.Page {
		return troveClusterPage{SinglePageBase: pagination.SinglePageBase(r)}
	})
}

var troveInstanceLabels = []string{"cluster_id", "datastore_type", "datastore_version", "health_status", "id", "name", "operating_status", "region", "status", "tenant_id"}

var troveClusterLabels = []string{"datastore_type", "datastore_version", "health_status", "id", "name", "task", "tenant_id"}

var defaultTroveMetrics = []Metric{
	{Name: "total_instances", Fn: ListAllInstances},
	{Name: "instance_status", Labels: troveInstanceLabels, Fn: nil},
	{Name: "instance_operating_status", Labels: troveInstanceLabels, Fn: nil},
	{Name: "instance_volume_size_gb", Labels: troveInstanceLabels, Fn: nil},
	{Name: "instance_volume_used_gb", Labels: troveInstanceLabels, Fn: nil},
	{Name: "total_clusters", Fn: ListAllDBClusters},
	{Name: "cluster_health", Labels: troveClusterLabels, Fn: nil},
	{Name: "cluster_members_out", Labels: troveClusterLabels, Fn: nil},
}

func NewTroveExporter(config *ExporterConfig, logger *slog.Logger) (*TroveExporter, error) {
	exporter := TroveExporter{
		BaseOpenStackExporter{
			Name:           "trove",
			ExporterConfig: *config,
			logger:         logger,
		},
	}

	for _, metric := range defaultTroveMetrics {
		if !exporter.isSlowMetric(&metric) {
			exporter.AddMetric(metric.Name, metric.Fn, metric.Labels, metric.DeprecatedVersion, nil)
		}
	}

	return &exporter, nil
}

func ListAllInstances(ctx context.Context, exporter *BaseOpenStackExporter, ch chan<- prometheus.Metric) error {
	var allInstances []instanceAttributesExt
	allPagesInstances, err := listDBInstances(exporter.ClientV2).AllPages(ctx)
	if err != nil {
		return err
	}

	allInstances, err = extractDBInstances(allPagesInstances)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(exporter.Metrics["total_instances"].Metric,
		prometheus.GaugeValue, float64(len(allInstances)))

	for _, instance := range allInstances {
		labelValues := []string{instance.ClusterID, instance.Datastore.Type, instance.Datastore.Version,
			instance.HealthStatus, instance.ID, instance.Name, instance.OperatingStatus, instance.Region,
			instance.Status, instance.TenantID}
		ch <- prometheus.MustNewConstMetric(exporter.Metrics["instance_status"].Metric,
			prometheus.GaugeValue, float64(mapDBInstanceStatus(instance.Status)), labelValues...)
		ch <- prometheus.MustNewConstMetric(exporter.Metrics["instance_operating_status"].Metric,
			prometheus.GaugeValue, float64(mapDBOperatingStatus(instance.OperatingStatus)), labelValues...)
		ch <- prometheus.MustNewConstMetric(exporter.Metrics["instance_volume_size_gb"].Metric,
			prometheus.GaugeValue, float64(instance.Volume.Size), labelValues...)
		ch <- prometheus.MustNewConstMetric(exporter.Metrics["instance_volume_used_gb"].Metric,
			prometheus.GaugeValue, instance.Volume.Used, labelValues...)
	}

	return nil
}

func ListAllDBClusters(ctx context.Context, exporter *BaseOpenStackExporter, ch chan<- prometheus.Metric) error {
	var allClusters []clusterAttributes
	allPagesClusters, err := listDBClusters(exporter.ClientV2).AllPages(ctx)
	if err != nil {
		return err
	}

	allClusters, err = extractDBClusters(allPagesClusters)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(exporter.Metrics["total_clusters"].Metric,
		prometheus.GaugeValue, float64(len(allClusters)))

	for _, cluster := range allClusters {
		labelValues := []string{cluster.Datastore.Type, cluster.Datastore.Version, cluster.Health.Status,
			cluster.ID, cluster.Name, cluster.Task.Name, cluster.TenantID}
		ch <- prometheus.MustNewConstMetric(exporter.Metrics["cluster_health"].Metric,
			prometheus.GaugeValue, float64(mapDBClusterHealthStatus(cluster.Health.Status)), labelValues...)
		ch <- prometheus.MustNewConstMetric(exporter.Metrics["cluster_members_out"].Metric,
			prometheus.GaugeValue, float64(len(cluster.Health.MembersOut)), labelValues...)
	}

	return nil
}
