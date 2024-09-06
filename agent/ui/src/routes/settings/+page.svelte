<script>
  import { onMount } from "svelte";
  import * as d3 from "d3";

  onMount(() => {
    const svg = d3
      .select("#gauge")
      .append("svg")
      .attr("width", 200)
      .attr("height", 200);

    const arc = d3
      .arc()
      .innerRadius(50)
      .outerRadius(70)
      .startAngle(0)
      .endAngle(Math.PI);

    svg
      .append("path")
      .attr("d", arc)
      .attr("fill", "steelblue")
      .attr("transform", "translate(100,100)");
  });
</script>

<h1>settings</h1>

<div id="gauge"></div>
