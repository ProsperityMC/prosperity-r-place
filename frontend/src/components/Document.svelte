<script lang="ts">
  import {onMount} from "svelte";
  import {Canvas} from "svelte-canvas";
  import type {Pixel} from "~/lib/Pixel";
  import {snapInsideLine} from "~/lib/SnapToBox";
  import Border from "./doc/Border.svelte";
  import Cursor from "./doc/Cursor.svelte";
  import Doc from "./doc/Doc.svelte";

  const offset = 32;

  export let doc;
  export let menuSel;
  export let zoomSel;
  export let paletteSel;
  export let scale;

  let canvasWidth = 0;
  let canvasHeight = 0;
  let scrollX = 0;
  let scrollY = 0;
  let startScrollX = 0;
  let startScrollY = 0;
  let mouseX = 0;
  let mouseY = 0;
  let startMouseX = 0;
  let startMouseY = 0;
  let holdMouse = false;
  let docCanvas;
  let docOverflow;
  let clientPixels: Pixel[] = [];

  // if the document overflows the canvas area
  $: docOverflow = offset * 2 + doc.width * scale > canvasWidth || offset * 2 + doc.height * scale > canvasHeight;

  // magic to calculate new scrollX position
  $: scrollX = docOverflow
    ? holdMouse && menuSel == "pan"
      ? -snapInsideLine(startMouseX - mouseX - startScrollX, canvasWidth, offset * 2 + doc.width * scale)
      : startScrollX
    : 0;

  // magic to calculate new scrollY position
  $: scrollY = docOverflow
    ? holdMouse && menuSel == "pan"
      ? -snapInsideLine(startMouseY - mouseY - startScrollY, canvasHeight, offset * 2 + doc.height * scale)
      : startScrollY
    : 0;

  let cellX = 0;
  let cellY = 0;

  // calculate the cellX and cellY of the mouse
  $: cellX = Math.floor((mouseX - scrollX - offset) / scale);
  $: cellY = Math.floor((mouseY - scrollY - offset) / scale);

  $: scale = zoomSel
    ? (() => {
        // don't auto scale if canvasWidth and canvasHeight aren't set
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

    // set the start dragging mouse position
    can.addEventListener("mousedown", e => {
      startMouseX = e.layerX;
      startMouseY = e.layerY;
      holdMouse = true;
    });

    // set the scroll position when releasing
    can.addEventListener("mouseup", () => {
      startScrollX = scrollX;
      startScrollY = scrollY;
      holdMouse = false;
    });

    can.addEventListener("mouseout", () => {
      startScrollX = scrollX;
      startScrollY = scrollY;
      mouseX = 0;
      mouseY = 0;
      holdMouse = false;
    });
  });
</script>

<div class="document {menuSel == 'pan' ? 'grab' : ''} {holdMouse ? 'grabbing' : ''}" bind:clientWidth={canvasWidth} bind:clientHeight={canvasHeight}>
  {#if scale >= 0}
    <Canvas width={canvasWidth} height={canvasHeight} style="position:absolute;" bind:this={docCanvas}>
      <Border {scrollX} {scrollY} docWidth={doc.width} docHeight={doc.height} {scale} />
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
