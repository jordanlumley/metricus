<script>
  export let title, content;

  // if content is longer than 35 characters, truncate it
  $: truncatedContent =
    content.length > 35 ? `${content.slice(0, 35)}...` : content;
</script>

<p>{title}</p>
<p>{truncatedContent}</p>
