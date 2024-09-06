<script>
  import { onMount } from "svelte";

  let metrics = {};

  onMount(() => {
    const logsEventSource = new EventSource(
      `http://localhost:8888/api/v1/metrics/events`
    );
    logsEventSource.onmessage = (event) => {
      metrics = JSON.parse(event.data);
    };

    return () => {
      logsEventSource.close();
    };
  });
</script>

<div>
  {#if !Object.keys(metrics).length}
    <p>loading</p>
  {:else}
    {#each Object.entries(metrics) as [key, value]}
      <strong>{key}</strong>: {value}
      <br />
    {/each}
  {/if}
</div>
