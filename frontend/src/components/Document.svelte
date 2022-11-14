<script lang="ts">
  import {onMount} from "svelte";
  import {Canvas} from "svelte-canvas";
  import DocumentCanvas from "./DocumentCanvas.svelte";

  export let doc;

  let canvasWidth = 1000;
  let canvasHeight = 1000;
  let mouseX = 0;
  let mouseY = 0;
  let docCanvas;

  onMount(() => {
    let can = docCanvas.getCanvas();
    console.log(can);
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
    <DocumentCanvas docWidth={doc.width} docHeight={doc.height} {mouseX} {mouseY} />
  </Canvas>
</div>

<style lang="scss">
  .document {
    height: 100%;
  }
</style>
