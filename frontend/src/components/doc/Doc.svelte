<script lang="ts">
  import {Layer, t} from "svelte-canvas";
  import {BufferImage} from "~/lib/BufferImage";
  import {getEnv} from "~/utils/env";

  const offset = 32;

  export let scrollX;
  export let scrollY;
  export let doc;
  export let scale = 1;

  let docImage = new BufferImage();
  let lastImage = 0;

  $: render = ({context: ctx}) => {
    // every 2s update the image using "_=timestamp"
    if ($t > lastImage + 2000) {
      lastImage = $t;
      docImage.update(`${getEnv("API_URL")}/doc/${doc.name}?raw=image&_${lastImage}`);
    }

    // render the current image
    ctx.drawImage(docImage.main, offset + scrollX, offset + scrollY, doc.width * scale, doc.height * scale);
  };
</script>

<Layer {render} />
