<script lang="ts">
  import {Layer, t} from "svelte-canvas";
  import {BufferImage} from "~/lib/BufferImage";

  const offset = 32;

  export let scrollX;
  export let scrollY;
  export let doc;
  export let scale = 1;
  export let updateImage = x => (internalImage = x);

  let internalImage = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==";
  export let docImage = new BufferImage(
    "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==",
  ,doc.width,doc.height);
  let lastImage = 0;

  $: render = ({context: ctx}) => {
    // every 2s update the image using "_=timestamp"
    if ($t > lastImage + 2000) {
      lastImage = $t;
      docImage.update(internalImage);
    }

    // render the current image
    ctx.imageSmoothingEnabled = false;
    ctx.drawImage(docImage.main, offset + scrollX, offset + scrollY, doc.width * scale, doc.height * scale);
  };
</script>

<Layer {render} />
