<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import type { Terminal } from "xterm";

  import { fetchContainer } from "$lib";
  import ContainerDetailEntry from "$lib/components/ContainerDetailEntry.svelte";

  export let data;

  let usedMemory,
    availableMemory,
    memoryUsagePercent: number = 0,
    cpuDelta,
    systemCpuDelta,
    numCpu,
    cpuUsagePercent = 0;

  // let terminalContainer: HTMLDivElement;

  // let logsEventSource: EventSource;
  let statsEventSource: EventSource;

  // let term: Terminal;

  onMount(async () => {
    // if (typeof window !== "undefined") {
    //   const { Terminal } = await import("xterm");
    //   import("xterm/css/xterm.css");

    //   term = new Terminal({
    //     cursorBlink: true,
    //     rows: 24,
    //     cols: 80,
    //     theme: {
    //       background: "#262C36",
    //     },
    //   });
    //   term.open(terminalContainer);
    // }

    await loadStats(data.slug);
    // await loadLogs(data.slug);
  });

  onDestroy(() => {
    // if (logsEventSource) {
    //   logsEventSource.close();
    // }

    if (statsEventSource) {
      statsEventSource.close();
    }
  });

  // // Reactive statement to handle containerId changes
  // $: if (term && data.slug) {
  //   term.clear();
  //   loadLogs(data.slug);
  // }

  const loadStats = async (id: any) => {
    return new Promise<void>((resolve, reject) => {
      if (statsEventSource) statsEventSource.close();

      statsEventSource = new EventSource(
        `http://localhost:8888/api/v1/containers/${id}/stats/test`
      );

      statsEventSource.onmessage = (event) => {
        setMetricsValues(JSON.parse(event.data));
        resolve();
      };
    });
  };

  // const loadLogs = async (id: any) => {
  //   return new Promise<void>((resolve, reject) => {
  //     if (logsEventSource) logsEventSource.close();

  //     logsEventSource = new EventSource(
  //       `http://localhost:8888/api/v1/containers/${id}/logs/events`
  //     );
  //     logsEventSource.onmessage = (event) => {
  //       term.write(event.data + "\r\n");
  //       resolve();
  //     };
  //   });
  // };

  const setMetricsValues = (stats: any) => {
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

<div class="container mx-auto">
  {#await fetchContainer(data.slug)}
    <p>loading</p>
  {:then container}
    <div class="grid grid-cols-3 gap-4">
      <div><ContainerDetailEntry title="Id" content={container.Id} /></div>
      <div>
        <ContainerDetailEntry title="Image" content={container.Image} />
      </div>
      <div>
        <ContainerDetailEntry title="Status" content={container.State.Status} />
      </div>
      <div><ContainerDetailEntry title="Name" content={container.Name} /></div>
      <div>
        <ContainerDetailEntry title="Created" content={container.Created} />
      </div>
    </div>

    <div class="divider" />

    <div class="grid grid-cols-3 gap-4">
      <div>
        <ContainerDetailEntry
          title="Memory Usage %"
          content={memoryUsagePercent}
        />
      </div>
      <div>
        <ContainerDetailEntry title="CPU Usage %" content={cpuUsagePercent} />
      </div>
    </div>
  {:catch error}
    <p>{error.message}</p>
  {/await}
</div>

<!-- <div bind:this={terminalContainer} id="terminal"></div> -->
