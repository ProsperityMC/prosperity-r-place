<script lang="ts">
  import {Layer} from "svelte-canvas";

  export let names: Map<string, string>;

  $: render = ({context: ctx}) => {
    ctx.font = "20px Verdana";
    ctx.fillStyle = "#dddddd";
    ctx.fillText(`Hello World!`, 0, 0);
    console.log("fill text");
    names.forEach((v, k) => {
      ctx.fillText(`${v} - ${k}`, 0, 0);
    });
  };
</script>

<Layer {render} />
