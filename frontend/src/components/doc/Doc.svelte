<script lang="ts">
  import {Layer} from "svelte-canvas";
  import {BufferImage} from "~/lib/BufferImage";

  const offset = 32;

  export let scrollX;
  export let scrollY;
  export let doc;
  export let scale = 1;
  export const updateImage = x => {
    docImage.update(x);
    docImage = docImage;
  };

  export let docImage = new BufferImage(
    "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==",
    doc.width,
    doc.height,
  );

  $: render = ({context: ctx}) => {
    // render the current image
    ctx.imageSmoothingEnabled = false;
    ctx.drawImage(docImage.main, offset + scrollX, offset + scrollY, doc.width * scale, doc.height * scale);
  };
</script>

<Layer {render} />
