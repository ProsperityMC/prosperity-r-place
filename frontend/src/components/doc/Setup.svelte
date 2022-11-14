<script lang="ts">
  import {Layer, t} from "svelte-canvas";
  import {BufferImage} from "~/lib/BufferImage";

  const offset = 32;

  export let docWidth;
  export let docHeight;
  export let mouseX;
  export let mouseY;
  export let scale = 1;

  let docImage = new BufferImage();
  let lastImage = 0;
  let scrollX = 0;
  let scrollY = 0;
  let frame = 0;

  $: render = ({context: ctx, width, height}) => {
    // SECTION: setup initial scale
    if (frame >= 1) {
      // size of canvas without offset margin
      let w = width - offset * 2;
      let h = height - offset * 2;

      // max scale in each axis
      let maxScaleX = w / docWidth;
      let maxScaleY = h / docHeight;

      // use the smaller scale so the document is scaled perfectly
      scale = maxScaleX < maxScaleY ? maxScaleX : maxScaleY;
      // set frame to -1 to stop
      frame = -1;
      return;
    }
    // increment to find the 2nd frame
    // when the clientWidth and clientHeight have been defined
    if (frame >= 0) frame++;
    // !SECTION

    // every 2s update the image using "_=timestamp"
    if ($t > lastImage + 2000) {
      lastImage = $t;
      docImage.update("https://localhost:5444/doc/r-place?raw=image&_=" + lastImage);
    }

    // calculate the cellX and cellY of the mouse
    let cellX = Math.floor((mouseX - scrollX - offset) / scale);
    let cellY = Math.floor((mouseY - scrollY - offset) / scale);

    // render the document border
    ctx.strokeStyle = "#ddd";
    ctx.lineWidth = 2;
    ctx.strokeRect(offset + scrollX, offset + scrollY, docWidth * scale, docHeight * scale);

    // render the current image
    ctx.drawImage(docImage.main, offset + scrollX, offset + scrollY, docWidth * scale, docHeight * scale);

    // if the cellX and cellY are within the document area
    if (cellX >= 0 && cellY >= 0 && cellX < docWidth && cellY < docHeight) {
      // render the cell selection box
      ctx.strokeStyle = "#aaa";
      ctx.lineWidth = 2;
      ctx.strokeRect(offset + scrollX + cellX * scale, offset + scrollY + cellY * scale, scale, scale);
    }
  };
</script>

<Layer {render} />
