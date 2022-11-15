<script lang="ts">
  import {Layer, t} from "svelte-canvas";
  import {BufferImage} from "~/lib/BufferImage";

  const offset = 32;

  export let scrollX;
  export let scrollY;
  export let docWidth;
  export let docHeight;
  export let scale = 1;

  let docImage = new BufferImage();
  let lastImage = 0;

  $: render = ({context: ctx}) => {
    // every 2s update the image using "_=timestamp"
    if ($t > lastImage + 2000) {
      lastImage = $t;
      docImage.update("https://localhost:5444/doc/r-place?raw=image&_=" + lastImage);
    }

    // render the current image
    ctx.drawImage(docImage.main, offset + scrollX, offset + scrollY, docWidth * scale, docHeight * scale);
  };
</script>

<Layer {render} />
