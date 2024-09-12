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

<div class="container mx-auto px-4">
  <table class="table table-lg">
    <thead>
      <tr>
        <th>Name</th>
        <th>Value</th>
      </tr>
    </thead>
    <tbody>
      {#each Object.entries(metrics) as [key, value]}
        <tr>
          <td>{key}</td>
          <td>{value}</td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
