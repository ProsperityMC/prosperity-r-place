<script lang="ts">
  import {Layer} from "svelte-canvas";

  const offset = 32;

  export let scrollX = 0;
  export let scrollY = 0;
  export let docWidth = 0;
  export let docHeight = 0;
  export let cellX = 0;
  export let cellY = 0;
  export let scale = 1;
  export let menuSel;

  $: render = ({context: ctx, width, height}) => {
    // if the cellX and cellY are within the document area
    if (cellX >= 0 && cellY >= 0 && cellX < docWidth && cellY < docHeight) {
      switch (menuSel) {
        case "pencil":
        case "fill":
        case "shape":
        case "select":
          // render the cell selection box
          ctx.strokeStyle = "#aaa";
          ctx.lineWidth = 2;
          ctx.strokeRect(offset + scrollX + cellX * scale, offset + scrollY + cellY * scale, scale, scale);
          break;
      }
    }
  };
</script>

<Layer {render} />
