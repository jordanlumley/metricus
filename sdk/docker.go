package metricus

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/ahmetalpbalkan/dlog"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerService struct {
	client *client.Client
}

func NewDockerService() (*DockerService, error) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &DockerService{client: apiClient}, nil
}

func (d *DockerService) Close() error {
	return d.client.Close()
}

func (d *DockerService) GetContainers(ctx context.Context) ([]types.Container, error) {
	return d.client.ContainerList(ctx, container.ListOptions{All: true})
}

func (d *DockerService) GetContainer(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	return d.client.ContainerInspect(ctx, containerID)
}

func (d *DockerService) GetLogs(ctx context.Context, containerID string) (string, error) {
	out, err := d.client.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "10",
	})
	if err != nil {
		return "", err
	}
	defer out.Close()

	p := make([]byte, 8)
	_, err = out.Read(p)
	if err != nil {
		return "", err
	}

	logs, err := io.ReadAll(out)
	if err != nil {
		return "", err
	}

	fmt.Println(string(logs))

	return string(logs), nil
}

func (d *DockerService) StreamLogs(ctx context.Context, containerID string, stream chan []byte) error {
	options := container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Tail: "15"}
	out, err := d.client.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return err
	}
	defer out.Close()

	r := dlog.NewReader(out)
	s := bufio.NewScanner(r)
	for s.Scan() {
		select {
		case <-ctx.Done():
			return nil
		default:
			stream <- s.Bytes()
		}
	}
	if err := s.Err(); err != nil {
		return err
	}

	return nil
}

type Stats struct {
	BlkioStats struct {
		IoMergedRecursive       any `json:"io_merged_recursive"`
		IoQueueRecursive        any `json:"io_queue_recursive"`
		IoServiceBytesRecursive []struct {
			Major int    `json:"major"`
			Minor int    `json:"minor"`
			Op    string `json:"op"`
			Value int    `json:"value"`
		} `json:"io_service_bytes_recursive"`
		IoServiceTimeRecursive any `json:"io_service_time_recursive"`
		IoServicedRecursive    any `json:"io_serviced_recursive"`
		IoTimeRecursive        any `json:"io_time_recursive"`
		IoWaitTimeRecursive    any `json:"io_wait_time_recursive"`
		SectorsRecursive       any `json:"sectors_recursive"`
	} `json:"blkio_stats"`
	CPUStats struct {
		CPUUsage struct {
			TotalUsage        int64 `json:"total_usage"`
			UsageInKernelmode int64 `json:"usage_in_kernelmode"`
			UsageInUsermode   int64 `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		OnlineCpus     int   `json:"online_cpus"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
		ThrottlingData struct {
			Periods          int `json:"periods"`
			ThrottledPeriods int `json:"throttled_periods"`
			ThrottledTime    int `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"cpu_stats"`
	ID          string `json:"id"`
	MemoryStats struct {
		Limit int64 `json:"limit"`
		Stats struct {
			ActiveAnon            int `json:"active_anon"`
			ActiveFile            int `json:"active_file"`
			Anon                  int `json:"anon"`
			AnonThp               int `json:"anon_thp"`
			File                  int `json:"file"`
			FileDirty             int `json:"file_dirty"`
			FileMapped            int `json:"file_mapped"`
			FileWriteback         int `json:"file_writeback"`
			InactiveAnon          int `json:"inactive_anon"`
			InactiveFile          int `json:"inactive_file"`
			KernelStack           int `json:"kernel_stack"`
			Pgactivate            int `json:"pgactivate"`
			Pgdeactivate          int `json:"pgdeactivate"`
			Pgfault               int `json:"pgfault"`
			Pglazyfree            int `json:"pglazyfree"`
			Pglazyfreed           int `json:"pglazyfreed"`
			Pgmajfault            int `json:"pgmajfault"`
			Pgrefill              int `json:"pgrefill"`
			Pgscan                int `json:"pgscan"`
			Pgsteal               int `json:"pgsteal"`
			Shmem                 int `json:"shmem"`
			Slab                  int `json:"slab"`
			SlabReclaimable       int `json:"slab_reclaimable"`
			SlabUnreclaimable     int `json:"slab_unreclaimable"`
			Sock                  int `json:"sock"`
			ThpCollapseAlloc      int `json:"thp_collapse_alloc"`
			ThpFaultAlloc         int `json:"thp_fault_alloc"`
			Unevictable           int `json:"unevictable"`
			WorkingsetActivate    int `json:"workingset_activate"`
			WorkingsetNodereclaim int `json:"workingset_nodereclaim"`
			WorkingsetRefault     int `json:"workingset_refault"`
		} `json:"stats"`
		Usage int `json:"usage"`
	} `json:"memory_stats"`
	Name     string `json:"name"`
	Networks struct {
		Eth0 struct {
			RxBytes   int `json:"rx_bytes"`
			RxDropped int `json:"rx_dropped"`
			RxErrors  int `json:"rx_errors"`
			RxPackets int `json:"rx_packets"`
			TxBytes   int `json:"tx_bytes"`
			TxDropped int `json:"tx_dropped"`
			TxErrors  int `json:"tx_errors"`
			TxPackets int `json:"tx_packets"`
		} `json:"eth0"`
	} `json:"networks"`
	NumProcs  int `json:"num_procs"`
	PidsStats struct {
		Current int `json:"current"`
		Limit   int `json:"limit"`
	} `json:"pids_stats"`
	PrecpuStats struct {
		CPUUsage struct {
			TotalUsage        int64 `json:"total_usage"`
			UsageInKernelmode int64 `json:"usage_in_kernelmode"`
			UsageInUsermode   int64 `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		OnlineCpus     int   `json:"online_cpus"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
		ThrottlingData struct {
			Periods          int `json:"periods"`
			ThrottledPeriods int `json:"throttled_periods"`
			ThrottledTime    int `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"precpu_stats"`
	Preread      time.Time `json:"preread"`
	Read         time.Time `json:"read"`
	StorageStats struct {
	} `json:"storage_stats"`
}

func (d *DockerService) GetContainerMetrics(ctx context.Context, containerID string) (*Stats, error) {
	stats, err := d.client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()

	var statsObj Stats
	if err := json.NewDecoder(stats.Body).Decode(&statsObj); err != nil {
		return nil, err
	}

	return &statsObj, nil
}

func (d *DockerService) StreamContainerMetrics(ctx context.Context, containerID string, stream chan []byte) error {
	stats, err := d.client.ContainerStats(ctx, containerID, true)
	if err != nil {
		return err
	}
	defer stats.Body.Close()

	reader := bufio.NewReader(stats.Body)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			line, err := reader.ReadBytes('\n')
			if err != nil && err != io.EOF {
				return err
			}
			if len(line) > 0 {
				stream <- line
			}
		}
	}
}
