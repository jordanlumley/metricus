<script>
  import { fetchContainer } from "$lib";
  import ContainerDetailEntry from "$lib/components/ContainerDetailEntry.svelte";

  export let containerId;
</script>

{#await fetchContainer(containerId)}
  <p>loading</p>
{:then container}
  <ContainerDetailEntry title="Id" content={container.Id} />
  <ContainerDetailEntry title="Image" content={container.Image} />
  <ContainerDetailEntry title="Status" content={container.State.Status} />
  <ContainerDetailEntry title="Name" content={container.Name} />
  <ContainerDetailEntry title="Created" content={container.Created} />
{:catch error}
  <p>{error.message}</p>
{/await}
