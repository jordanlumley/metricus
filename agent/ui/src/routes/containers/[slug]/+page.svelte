<script>
  import ContainerDetails from "$lib/components/ContainerDetails.svelte";
  import ContainerStats from "$lib/components/ContainerStats.svelte";
  import ContainerLogs from "$lib/components/ContainerLogs.svelte";

  export let data;
  let logs = [];
</script>

<ContainerDetails containerId={data.slug} />
<ContainerStats containerId={data.slug} />
<ContainerLogs containerId={data.slug} />
