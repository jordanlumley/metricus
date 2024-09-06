<script>
  import { fetchContainer } from "$lib";
  import ContainerDetailEntry from "$lib/components/ContainerDetailEntry.svelte";

  export let containerId;
</script>

<div class="card">
  <header class="card-header">
    <p class="card-header-title">Container</p>
  </header>
  <div class="card-content">
    {#await fetchContainer(containerId)}
      <div class="skeleton-block"><div></div></div>
      <div class="skeleton-block"><div></div></div>
      <div class="skeleton-block"><div></div></div>
      <div class="skeleton-block"><div></div></div>
      <div class="skeleton-block"><div></div></div>
    {:then container}
      <div class="fixed-grid">
        <div class="grid">
          <div class="cell">
            <ContainerDetailEntry title="Id" content={container.Id} />
          </div>
          <div class="cell">
            <ContainerDetailEntry title="Image" content={container.Image} />
          </div>
          <div class="cell">
            <ContainerDetailEntry
              title="Status"
              content={container.State.Status}
            />
          </div>
          <div class="cell">
            <ContainerDetailEntry title="Name" content={container.Name} />
          </div>
          <div class="cell">
            <ContainerDetailEntry title="Created" content={container.Created} />
          </div>
        </div>
      </div>
    {:catch error}
      <p>{error.message}</p>
    {/await}
  </div>
</div>
