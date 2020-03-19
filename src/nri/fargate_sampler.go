package nri

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"

	"github.com/newrelic/nri-docker/src/raw/aws"
)

const clusterNameLabel = "com.amazonaws.ecs.cluster"

// FargateSampler is a sampler to get container data from an HTTP endpoint.
type FargateSampler struct {
}

// SampleAll will sample all the information this sampler can fetch and populate it in the given integration.
func (f *FargateSampler) SampleAll(i *integration.Integration) error {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	fetcher, err := aws.NewFargateFetcher(httpClient, i.Logger())
	if err != nil {
		return fmt.Errorf("could not create fargate fetcher: %v", err)
	}

	dockerMetrics, err := fetcher.GetContainerMetrics()
	if err != nil {
		return fmt.Errorf("error fetching fargate docker metrics: %v", err)
	}

	for containerID, stats := range *dockerMetrics {
		if stats == nil {
			i.Logger().Errorf("no stats found for container %s", containerID)
			continue
		}
		entity, err := i.Entity(containerID, "fargate")
		if err != nil {
			return fmt.Errorf("could not create entity: %v", err)
		}

		ms := entity.NewMetricSet(
			containerSampleName,
			metric.Attr(attrContainerID, containerID),
		)
		bizMetrics := fetcher.BizMetricsFromDocker(containerID, stats)
		populate(ms, misc(&bizMetrics))
		populate(ms, cpu(&bizMetrics.CPU))
		populate(ms, memory(&bizMetrics.Memory))
		populate(ms, pids(&bizMetrics.Pids))
		populate(ms, blkio(&bizMetrics.BlkIO))

		containerMetadata, err := fetcher.InspectContainer(containerID)
		if err != nil {
			i.Logger().Errorf("could not find metadata for container %s: %v", containerID, err)
			continue
		}
		// We have to standardize the cluster name label because in the EC2 launch type it's purely the cluster name. In
		// Fargate though it's the complete cluster's ARN.
		containerMetadata.Labels[clusterNameLabel] = clusterNameFromARN(containerMetadata.Labels[clusterNameLabel])
		populate(ms, attributes(containerMetadata))
		populate(ms, labels(containerMetadata))
	}
	return nil
}

func clusterNameFromARN(clusterARN string) string {
	return strings.Split(clusterARN, "/")[1]
}
