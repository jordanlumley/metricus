<script>
  import { fetchContainers } from "$lib";
  import { onMount } from "svelte";

  onMount(() => {});
</script>

<div class="columns is-mobile">
  <div class="column is-3">
    {#await fetchContainers()}
      <p>Loading...</p>
    {:then containers}
      <ul>
        {#each containers as container}
          <li>
            <a data-sveltekit-reload href="/containers/{container.Id}"
              >{container.Image}</a
            >
          </li>
        {/each}
      </ul>
    {:catch error}
      <p>{error.message}</p>
    {/await}
  </div>
  <div class="column"><slot></slot></div>
</div>
