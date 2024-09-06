<script>
  import { onMount } from "svelte";
  import ContainerDetailEntry from "./ContainerDetailEntry.svelte";

  export let containerId;

  let stats = null;
  let usedMemory = 0;
  let availableMemory = 0;
  let memoryUsagePercent = 0;
  let cpuDelta = 0;
  let systemCpuDelta = 0;
  let numCpu = 0;
  let cpuUsagePercent = 0;

  onMount(() => {
    const statsEventSource = new EventSource(
      `http://localhost:8888/api/v1/containers/${containerId}/stats/test`
    );
    statsEventSource.onmessage = (event) => {
      stats = JSON.parse(event.data);
      setMetricsValues(stats);
    };

    return () => {
      statsEventSource.close();
    };
  });

  const setMetricsValues = (stats) => {
    usedMemory = stats.memory_stats.usage;
    availableMemory = stats.memory_stats.limit;
    memoryUsagePercent = (usedMemory / availableMemory) * 100.0;
    cpuDelta =
      stats.cpu_stats.cpu_usage.total_usage -
      stats.precpu_stats.cpu_usage.total_usage;
    systemCpuDelta =
      stats.cpu_stats.system_cpu_usage - stats.precpu_stats.system_cpu_usage;
    numCpu = stats.cpu_stats.online_cpus;
    cpuUsagePercent = (cpuDelta / systemCpuDelta) * numCpu * 100.0;
  };
</script>

<svelte:head>
  <script src="https://cdn.jsdelivr.net/npm/c3@0.7.20/c3.min.js"></script>
  <link
    href="https://cdn.jsdelivr.net/npm/c3@0.7.20/c3.min.css"
    rel="stylesheet"
  />
</svelte:head>
<div class="card">
  <header class="card-header">
    <p class="card-header-title">Container</p>
  </header>
  <div class="card-content">
    {#if !stats}
      <div class="skeleton-block"><div></div></div>
      <div class="skeleton-block"><div></div></div>
      <div class="skeleton-block"><div></div></div>
      <div class="skeleton-block"><div></div></div>
      <div class="skeleton-block"><div></div></div>
    {:else}
      <div class="fixed-grid">
        <div class="grid">
          <div class="cell">
            <ContainerDetailEntry
              title="Memory Usage %"
              content={memoryUsagePercent}
            />
          </div>
          <div class="cell">
            <ContainerDetailEntry
              title="Memory Limit"
              content={stats.memory_stats.limit}
            />
          </div>
          <div class="cell">
            <ContainerDetailEntry
              title="CPU Usage %"
              content={cpuUsagePercent}
            />
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>
