package docker

import (
	"math/rand"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

type Sampler interface {
	Populate(*metric.Set) error
}

type ContainerSampler struct {
}

func (c ContainerSampler) Populate(ms *metric.Set) error {
	rndCpu := 10 + rand.Float64()*10
	for _, metric := range []Metric{
		MetricCommandLine("/command-exec"),
		MetricUser("root"),
		MetricContainerImage("12345"),
		MetricContainerImageName("alpine:latest"),
		MetricContainerName("containername"),
		MetricContainerID("123456"),
		MetricState("running"),
		{"label.docker.meta", metric.ATTRIBUTE, "label-value"},
		{"label.docker.io.kubernetes.pod.namespace", metric.ATTRIBUTE, "default"},
		MetricCPUPercent(rndCpu),
		MetricCPUSystemPercent(rndCpu * 0.2),
		MetricCPUUserPercent(rndCpu * 0.8),
		MetricMemoryVirtualSizeBytes(10000000),
		MetricMemoryResidentSizeBytes(8000000),
		MetricIOReadCountPerSecond(rand.Intn(10000000)), // take from blkio_stats
		MetricIOWriteCountPerSecond(rand.Intn(10000000)),
		MetricIOReadBytesPerSecond(rand.Intn(10000000)),
		MetricIOWriteBytesPerSecond(rand.Intn(10000000)),
		MetricIOTotalReadCount(rand.Intn(10000000)),
		MetricIOTotalWriteCount(rand.Intn(10000000)),
		MetricIOTotalReadBytes(rand.Intn(10000000)),
		MetricIOTotalWriteBytes(rand.Intn(10000000)),
		MetricPIDs(1),
	} {
		if err := ms.SetMetric(metric.Name, metric.Value, metric.Type); err != nil {
			return err
		}
	}
	return nil
}

func NewContainerSampler() Sampler {
	return ContainerSampler{}
}
