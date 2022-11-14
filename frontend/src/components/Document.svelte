<script lang="ts">
  import {onMount} from "svelte";
  import {Canvas} from "svelte-canvas";
  import DocumentCanvas from "./doc/Setup.svelte";

  export let doc;
  export let menuSel;
  export let scale;

  let canvasWidth = 1000;
  let canvasHeight = 1000;
  let mouseX = 0;
  let mouseY = 0;
  let holdMouse = false;
  let docCanvas;

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
  <Canvas width={canvasWidth} height={canvasHeight} style="position:absolute;" bind:this={docCanvas}>
    <DocumentCanvas docWidth={doc.width} docHeight={doc.height} {mouseX} {mouseY} bind:scale {...$$restProps} />
  </Canvas>
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
