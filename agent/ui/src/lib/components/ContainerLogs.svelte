<script>
  import { onMount } from "svelte";

  export let containerId;

  let logs = [];
  let terminalContainer;
  let term;

  onMount(async () => {
    term = new Terminal();
    term.open(terminalContainer);

    const logsEventSource = new EventSource(
      `http://localhost:8888/api/v1/containers/${containerId}/logs/events`
    );
    logsEventSource.onmessage = (event) => {
      //   logs = [event.data, ...logs];
      term.write(event.data + "\r\n");
    };

    return () => {
      logsEventSource.close();
    };
  });
</script>

<svelte:head>
  <script
    src="https://cdn.jsdelivr.net/npm/xterm@5.3.0/lib/xterm.min.js"
  ></script>
  <link
    href="https://cdn.jsdelivr.net/npm/xterm@5.3.0/css/xterm.min.css"
    rel="stylesheet"
  />
</svelte:head>

<div class="card">
  <header class="card-header">
    <p class="card-header-title">Container</p>
  </header>
  <div class="card-content" id="logsView">
    <div bind:this={terminalContainer} id="terminal"></div>
    {#each logs as log}
      <p>{log}</p>
    {/each}
  </div>
</div>
