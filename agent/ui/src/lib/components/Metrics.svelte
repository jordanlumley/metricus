<script>
  import { onMount } from "svelte";

  let metrics = [];

  onMount(() => {
    const logsEventSource = new EventSource(
      `http://localhost:8888/api/v1/metrics/events`
    );
    logsEventSource.onmessage = (event) => {
      metrics = [event.data, ...metrics];
    };

    return () => {
      logsEventSource.close();
    };
  });
</script>

<div class="card">
  <header class="card-header">
    <p class="card-header-title">Metrics</p>
  </header>
  <div class="card-content" id="logsView">
    <div class="content">
      {#each metrics as metric}
        <p>
          {#each Object.entries(JSON.parse(metric)) as [key, value]}
            <strong>{key}</strong>: {value}
            <br />
          {/each}
        </p>
      {/each}
    </div>
  </div>
</div>
