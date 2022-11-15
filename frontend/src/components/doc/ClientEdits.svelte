<script lang="ts">
  import colourPalette from "~/assets/colours.json";
  import {Layer} from "svelte-canvas";

  const offset = 32;

  export let scrollX = 0;
  export let scrollY = 0;
  export let docWidth = 0;
  export let docHeight = 0;
  export let scale = 1;
  export let clientPixels: Uint16Array;

  $: render = ({context: ctx}) => {
    if (clientPixels.length === docWidth * docHeight) {
      for (let x = 0; x < docWidth; x++) {
        for (let y = 0; y < docHeight; y++) {
          let n = clientPixels[y * docHeight + x];
          let upper = (n >> 8) & 0xff;
          let lower = n & 0xff;
          let pixel = colourPalette[upper].options[lower];
          ctx.fillStyle = pixel.hex;
          ctx.fillRect(offset + scrollX + x * scale, offset + scrollY + y * scale, scale, scale);
        }
      }
    } else {
      ctx.fillStyle = "#ff000066";
      ctx.fillRect(offset + scrollX, offset + scrollY, docWidth * scale + docHeight * scale);
    }
  };
</script>

<Layer {render} />
