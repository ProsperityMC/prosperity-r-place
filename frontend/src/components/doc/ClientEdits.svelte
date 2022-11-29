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
  export let selArea: {x1: number; y1: number; x2: number; y2: number};

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

      if (selArea.x1 !== -1 && selArea.y1 !== -1 && selArea.x2 !== -1 && selArea.y2 !== -1) {
        let selCheckX = selArea.x1 < selArea.x2;
        let selCheckY = selArea.y1 < selArea.y2;
        let selLowX = selCheckX ? selArea.x1 : selArea.x2;
        let selHighX = selCheckX ? selArea.x2 : selArea.x1;
        let selLowY = selCheckY ? selArea.y1 : selArea.y2;
        let selHighY = selCheckY ? selArea.y2 : selArea.y1;

        // select box
        ctx.fillStyle = "#ff0000";
        // top
        ctx.fillRect(offset + scrollX + selLowX * scale - 2, offset + scrollY + selLowY * scale - 2, Math.abs(selHighX - selLowX + 1) * scale + 4, 4);
        // left
        ctx.fillRect(offset + scrollX + selLowX * scale - 2, offset + scrollY + selLowY * scale - 2, 4, Math.abs(selHighY - selLowY + 1) * scale + 4);
        // right
        ctx.fillRect(
          offset + scrollX + (selHighX + 1) * scale - 2,
          offset + scrollY + selLowY * scale - 2,
          4,
          Math.abs(selHighY - selLowY + 1) * scale + 4,
        );
        // bottom
        ctx.fillRect(
          offset + scrollX + selLowX * scale - 2,
          offset + scrollY + (selHighY + 1) * scale - 2,
          Math.abs(selHighX - selLowX + 1) * scale + 4,
          4,
        );
      }
    } else {
      ctx.fillStyle = "#ff000066";
      ctx.fillRect(offset + scrollX, offset + scrollY, docWidth * scale + docHeight * scale);
    }
  };
</script>

<Layer {render} />
