<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import type { Terminal } from "xterm";

  export let containerId;

  let terminalContainer;
  let term: Terminal;
  let logsEventSource;

  onMount(async () => {
    if (typeof window !== "undefined") {
      const { Terminal } = await import("xterm");
      import("xterm/css/xterm.css");

      term = new Terminal({
        cursorBlink: true,
        rows: 24,
        cols: 80,
        theme: {
          background: "#262C36",
        },
      });
      term.open(terminalContainer);
    }
  });

  // Reactive statement to handle containerId changes
  $: if (term && containerId) {
    term.clear();
    loadLogs(containerId);
  }

  function loadLogs(id) {
    // If an EventSource is already open, close it
    if (logsEventSource) {
      logsEventSource.close();
    }

    logsEventSource = new EventSource(
      `http://localhost:8888/api/v1/containers/${id}/logs/events`
    );
    logsEventSource.onmessage = (event) => {
      term.write(event.data + "\r\n");
    };
  }

  onDestroy(() => {
    if (logsEventSource) {
      logsEventSource.close();
    }
  });
</script>

<div bind:this={terminalContainer} id="terminal"></div>
