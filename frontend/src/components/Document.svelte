<script lang="ts">
  import {onMount} from "svelte";
  import {Canvas} from "svelte-canvas";
  import Border from "./doc/Border.svelte";
  import Cursor from "./doc/Cursor.svelte";
  import Doc from "./doc/Doc.svelte";

  const offset = 32;

  export let doc;
  export let menuSel;
  export let zoomSel;
  export let scale;

  let canvasWidth = 0;
  let canvasHeight = 0;
  let scrollX = 0;
  let scrollY = 0;
  let mouseX = 0;
  let mouseY = 0;
  let holdMouse = false;
  let docCanvas;

  let cellX = 0;
  let cellY = 0;

  // calculate the cellX and cellY of the mouse
  $: cellX = Math.floor((mouseX - scrollX - offset) / scale);
  $: cellY = Math.floor((mouseY - scrollY - offset) / scale);

  $: scale = zoomSel
    ? (() => {
        if (canvasWidth === 0 && canvasHeight === 0) return 1;

        // size of canvas without offset margin
        let w = canvasWidth - offset * 2;
        let h = canvasHeight - offset * 2;

        // max scale in each axis
        let maxScaleX = w / doc.width;
        let maxScaleY = h / doc.height;

        // use the smaller scale so the document is scaled perfectly
        return maxScaleX < maxScaleY ? maxScaleX : maxScaleY;
      })()
    : scale;

  onMount(() => {
    // small hack to get passive mouse move events
    let can = docCanvas.getCanvas();
    can.addEventListener(
      "mousemove",
      e => {
        mouseX = e.layerX;
        mouseY = e.layerY;
      },
      {passive: true},
    );
    can.addEventListener("mousedown", () => {
      holdMouse = true;
    });
    can.addEventListener("mouseup", () => {
      holdMouse = false;
    });
  });
</script>

<div class="document {menuSel == 'pan' ? 'grab' : ''} {holdMouse ? 'grabbing' : ''}" bind:clientWidth={canvasWidth} bind:clientHeight={canvasHeight}>
  {#if scale >= 0}
    <Canvas width={canvasWidth} height={canvasHeight} style="position:absolute;" bind:this={docCanvas}>
      <Border docWidth={doc.width} docHeight={doc.height} bind:scale {...$$restProps} />
      <Doc {scrollX} {scrollY} docWidth={doc.width} docHeight={doc.height} {scale} />
      <Cursor docWidth={doc.width} docHeight={doc.height} {cellX} {cellY} {scrollX} {scrollY} {scale} {menuSel} />
    </Canvas>
  {:else}
    <div>Invalid document scale</div>
  {/if}
</div>

<style lang="scss">
  .document {
    height: 100%;

    &.grab {
      cursor: grab;

      &.grabbing {
        cursor: grabbing;
      }
    }
  }
</style>
