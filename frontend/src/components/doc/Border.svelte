<script lang="ts">
  import {Layer, t} from "svelte-canvas";

  const offset = 32;

  export let scrollX = 0;
  export let scrollY = 0;
  export let docWidth = 0;
  export let docHeight = 0;
  export let scale = 1;

  $: render = ({context: ctx}) => {
    // render the document border
    ctx.strokeStyle = "#ddd";
    ctx.lineWidth = 2;
    ctx.strokeRect(offset + scrollX, offset + scrollY, docWidth * scale, docHeight * scale);
  };
</script>

<Layer {render} />
