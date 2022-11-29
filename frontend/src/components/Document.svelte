<script lang="ts">
  import {onDestroy, onMount} from "svelte";
  import {Canvas} from "svelte-canvas";
  import {AlwaysOnWS} from "~/lib/AlwaysOnWS";
  import {snapInsideLine} from "~/lib/SnapToBox";
  import {getEnv} from "~/utils/env";
  import Border from "./doc/Border.svelte";
  import ClientEdits from "./doc/ClientEdits.svelte";
  import Cursor from "./doc/Cursor.svelte";
  import Doc from "./doc/Doc.svelte";
  import colourPalette from "~/assets/colours.json";
  import {GenerateShapePixels} from "~/lib/GenerateShapePixels";

  const offset = 32;

  export let doc;
  export let menuSel;
  export let shapeSel;
  export let zoomSel;
  export let paletteSel: number;
  export let scale;
  export let desel = () => (selArea = {x1: -1, y1: -1, x2: -1, y2: -1});

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
  let clientPixels: Uint16Array = new Uint16Array(doc.width * doc.height);
  let shapeArea = {x1: 0, y1: 0, x2: 0, y2: 0};
  let selArea = {x1: -1, y1: -1, x2: -1, y2: -1};
  $: lockArea =
    selArea.x1 !== -1 && selArea.y1 !== -1 && selArea.x2 !== -1 && selArea.y2 !== -1
      ? {...selArea}
      : {
          x1: 0,
          y1: 0,
          x2: doc.width,
          y2: doc.height,
        };

  let ws: AlwaysOnWS;
  let clock: number;
  let eTag: string;
  let updateImage;

  async function connectToWebsocket() {
    let wsUrl = `${getEnv("API_URL").replace("https://", "wss://")}/doc/${doc.name}`;
    console.log(`Connecting to WS: ${wsUrl}`);
    let openWS = new AlwaysOnWS(wsUrl);
    openWS.onopen = function () {
      ws = openWS;
      clock = setInterval(() => _onclock(), 2000);
    };
    openWS.onmessage = function (x) {
      let args = x.data.split(" ");
      if (args.length < 1) return;
      switch (args[0]) {
        case "refresh":
          if (args.length !== 3) return;
          eTag = args[1];
          updateImage("data:image/png;base64," + args[2]);
          break;
      }
    };
    openWS.close = function () {
      clearInterval(clock);
    };
  }

  function _onclock() {
    ws.send("check " + eTag);
  }

  onMount(async () => {
    console.log(`Starting editor for ${doc.name} - ${doc.width}x${doc.height}`);
    connectToWebsocket();
  });

  onDestroy(() => {
    ws.close();
  });

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

        // don't auto scale if the max scale is less than 0
        if (maxScaleX < 0 || maxScaleY < 0) return 1;

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
        checkDraw();
      },
      {passive: true},
    );

    // set the start dragging mouse position
    can.addEventListener("mousedown", e => {
      startMouseX = e.layerX;
      startMouseY = e.layerY;
      mouseX = e.layerX;
      mouseY = e.layerY;
      holdMouse = true;
      switch (menuSel) {
        case "shape":
          if (cellX >= lockArea.x1 && cellX <= lockArea.x2) shapeArea.x1 = cellX;
          if (cellY >= lockArea.y1 && cellY <= lockArea.y2) shapeArea.y1 = cellY;
          shapeArea = shapeArea;
          break;
        case "select":
          if (cellX >= lockArea.x1 && cellX <= lockArea.x2) selArea.x1 = cellX;
          if (cellY >= lockArea.y1 && cellY <= lockArea.y2) selArea.y1 = cellY;
          selArea = selArea;
          break;
      }
      checkDraw();
    });

    // set the scroll position when releasing
    can.addEventListener("mouseup", () => {
      startScrollX = scrollX;
      startScrollY = scrollY;
      holdMouse = false;
      releaseMouse();
    });

    can.addEventListener("mouseout", () => {
      startScrollX = scrollX;
      startScrollY = scrollY;
      mouseX = 0;
      mouseY = 0;
      holdMouse = false;
      releaseMouse();
    });
  });

  function checkDraw() {
    if (holdMouse) {
      switch (menuSel) {
        case "pencil":
          if (cellX >= lockArea.x1 && cellY >= lockArea.y1 && cellX <= lockArea.x2 && cellY <= lockArea.y2) {
            clientPixels[cellY * doc.height + cellX] = paletteSel;
            clientPixels = clientPixels;
          }
          break;
        case "shape":
          if (cellX >= lockArea.x1 && cellY >= lockArea.y1 && cellX <= lockArea.x2 && cellY <= lockArea.y2) {
            shapeArea.x2 = cellX;
            shapeArea.y2 = cellY;
            shapeArea = shapeArea;
            clientPixels = GenerateShapePixels(doc.width, doc.height, shapeArea, shapeSel, paletteSel);
          }
          break;
        case "select":
          selArea.x2 = cellX;
          selArea.y2 = cellY;
          selArea = selArea;
          break;
      }
    }
  }

  function releaseMouse() {
    let pixels: Map<string, {x: number; y: number}[]> = new Map();
    for (let j = 0; j < doc.height; j++) {
      for (let i = 0; i < doc.width; i++) {
        let coord = clientPixels[j * doc.width + i];
        if (coord === 0) continue;

        let upper = (coord >> 8) & 0xff;
        let lower = coord & 0xff;

        let pixel = colourPalette[upper].options[lower];
        if (!pixels.has(pixel.hex)) pixels.set(pixel.hex, []);
        pixels.get(pixel.hex).push({x: i, y: j});
      }
    }

    [...pixels.keys()].forEach(x => {
      let colour = pixels.get(x);
      ws.send(`draw ${x} ${colour.map(x => `${x.x},${x.y}`).join(" ")}`);
    });
    clientPixels = new Uint16Array(doc.width * doc.height);
  }
</script>

<div class="document {menuSel == 'pan' ? 'grab' : ''} {holdMouse ? 'grabbing' : ''}" bind:clientWidth={canvasWidth} bind:clientHeight={canvasHeight}>
  <Canvas width={canvasWidth} height={canvasHeight} style="position:absolute;" bind:this={docCanvas}>
    {#if scale >= 0}
      <Border {scrollX} {scrollY} docWidth={doc.width} docHeight={doc.height} {scale} />
      <Doc {scrollX} {scrollY} {doc} {scale} bind:updateImage />
      <ClientEdits {scrollX} {scrollY} docWidth={doc.width} docHeight={doc.height} {scale} {clientPixels} {selArea} />
      <Cursor docWidth={doc.width} docHeight={doc.height} {cellX} {cellY} {scrollX} {scrollY} {scale} {menuSel} />
    {/if}
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
