<script lang="ts">
  import {onMount} from "svelte";
  import {Canvas} from "svelte-canvas";
  import DocumentCanvas from "./DocumentCanvas.svelte";

  export let doc;
  export let menuSel;
  export let scale;

  let canvasWidth = 1000;
  let canvasHeight = 1000;
  let mouseX = 0;
  let mouseY = 0;
  let docCanvas;

  onMount(() => {
    // small hack to get passive mouse move events
    let can = docCanvas.getCanvas();
    can.addEventListener(
      "mousemove",
      function (e) {
        mouseX = e.layerX;
        mouseY = e.layerY;
      },
      {passive: true},
    );
  });
</script>

<div class="document" bind:clientWidth={canvasWidth} bind:clientHeight={canvasHeight}>
  <Canvas width={canvasWidth} height={canvasHeight} style="position:absolute;" bind:this={docCanvas}>
    <DocumentCanvas docWidth={doc.width} docHeight={doc.height} {mouseX} {mouseY} bind:scale {...$$restProps} />
  </Canvas>
</div>

<style lang="scss">
  .document {
    height: 100%;
  }
</style>
