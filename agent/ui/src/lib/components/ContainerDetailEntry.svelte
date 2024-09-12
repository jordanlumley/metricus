<script lang="ts">
  export let title, content;

  let truncatedContent: any;

  // if content is longer than 35 characters, truncate it
  $: if (content) {
    truncatedContent =
      content.length > 35 ? `${content.slice(0, 35)}...` : content;
  }
</script>

<p>{title}</p>
<p>{truncatedContent}</p>
