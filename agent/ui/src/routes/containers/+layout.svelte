<script>
  import { fetchContainers } from "$lib";
</script>

<div class="grid">
  <div>
    {#await fetchContainers()}
      <p>loading</p>
    {:then containers}
      <ul>
        {#each containers as container}
          <li>
            <a href="/containers/{container.Id}">{container.Image}</a>
          </li>
        {/each}
      </ul>
    {:catch error}
      <p>{error.message}</p>
    {/await}
  </div>
  <div>
    <slot></slot>
  </div>
</div>
