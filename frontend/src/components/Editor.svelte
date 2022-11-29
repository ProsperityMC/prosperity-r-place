<script lang="ts">
  import {createEventDispatcher} from "svelte";
  import colourPalette from "~/assets/colours.json";
  import Document from "./Document.svelte";
  import {
    Menu,
    Palette,
    X,
    Hand,
    Edit2,
    PaintBucket,
    Square,
    BoxSelect,
    Slash,
    Circle,
    Triangle,
    Diamond,
    Hexagon,
    Octagon,
    ZoomOut,
    ZoomIn,
    Monitor,
  } from "lucide-svelte";

  export let doc;

  const dispatch = createEventDispatcher();

  function triggerExit() {
    dispatch("exit", {});
  }

  let showMenu = true;
  let showPalette = true;
  let menuSel = "pan";
  let shapeSel = "circle";
  let paletteSel = 0; // numeric for special->transparent
  let zoomSel = -1;
  let scale = 1;
  let desel;

  const menuButtons = [
    {key: "pan", icon: Hand, select: true},
    {key: "pencil", icon: Edit2, select: true},
    {key: "fill", icon: PaintBucket, select: true},
    {key: "shape", icon: Square, select: true},
    {key: "select", icon: BoxSelect, select: true},
    {key: "deselect", icon: Slash, select: false, click: () => desel()},
  ];

  const shapeButtons = [
    {key: "circle", icon: Circle, select: true},
    {key: "triangle", icon: Triangle, select: true},
    {key: "square", icon: Square, select: true},
    {key: "diamond", icon: Diamond, select: true},
    {key: "hexagon", icon: Hexagon, select: true},
    {key: "octagon", icon: Octagon, select: true},
  ];
</script>

<div id="editor">
  <div id="editor-navbar">
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="tool-button" on:click={() => (showMenu = !showMenu)}>
      <Menu />
    </div>
    <div class="flex-gap" />
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="tool-button" on:click={() => (showPalette = !showPalette)}>
      <Palette />
    </div>
  </div>
  <div id="editor-main">
    {#if showMenu}
      <div id="editor-tools">
        {#each menuButtons as b (b.key)}
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <div class="tool-button {menuSel == b.key ? 'tool-button-sel' : ''}" on:click={() => (b.select ? (menuSel = b.key) : b.click())}>
            <svelte:component this={b.icon} />
          </div>
        {/each}
        <div class="flex-gap" />
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <div class="tool-button" on:click={() => triggerExit()}>
          <X />
        </div>
      </div>
      {#if menuSel == "shape"}
        <div id="editor-shapes">
          {#each shapeButtons as b (b.key)}
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div class="tool-button {shapeSel == b.key ? 'tool-button-sel' : ''}" on:click={() => b.select && (shapeSel = b.key)}>
              <svelte:component this={b.icon} />
            </div>
          {/each}
        </div>
      {/if}
    {/if}
    <div id="editor-doc">
      <Document {doc} {menuSel} {shapeSel} {zoomSel} {paletteSel} bind:scale bind:desel />
    </div>
    {#if showPalette}
      <div id="editor-palette">
        {#each colourPalette as palette, paletteI}
          <div class="palette-panel">
            {#each palette.options as option, optionI}
              <!-- svelte-ignore a11y-click-events-have-key-events -->
              <div
                class="palette-button {paletteSel === (paletteI << 8) + optionI ? 'palette-selected' : ''}"
                title="{option.name} {palette.name}"
                on:click={() => (paletteSel = (paletteI << 8) + optionI)}
              >
                {#if option.name == "Transparent" && palette.name == "Special"}
                  <div class="palette-button-blob" data-transparent />
                {:else}
                  <div class="palette-button-blob" style="background-color:{option.hex};" />
                {/if}
              </div>
            {/each}
          </div>
        {/each}
      </div>
    {/if}
  </div>
  <div id="editor-statusbar">
    <div class="flex-gap" />
    <div id="editor-zoom">
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <div class="icon" on:click={() => (scale -= 1) && (zoomSel = 0)}>
        <ZoomOut />
      </div>
      <div id="zoom-value">{Math.floor(scale * 100)}%</div>
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <div class="icon" on:click={() => (scale += 1) && (zoomSel = 0)}>
        <ZoomIn />
      </div>
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <div class="icon {zoomSel == -1 ? 'zoom-sel' : ''}" on:click={() => (zoomSel = -1)}>
        <Monitor />
      </div>
    </div>
  </div>
</div>

<style lang="scss">
  @import "../assets/theme.scss";

  #editor {
    display: flex;
    flex-direction: column;
    flex-basis: auto;
    width: 100vw;
    height: 100vh;

    > #editor-navbar {
      display: flex;
      flex-direction: row;
      background-color: darken($theme-bg, 3);
      height: 48px;
    }

    > #editor-main {
      display: flex;
      flex-direction: row;
      flex-grow: 1;
      min-height: 0;

      > #editor-tools {
        display: flex;
        flex-direction: column;
        width: 48px;
        background-color: darken($theme-bg, 3);
      }

      > #editor-shapes {
        display: flex;
        flex-direction: column;
        width: 48px;
        background-color: darken($theme-bg, 6);
      }

      > #editor-doc {
        overflow: hidden;
        flex-grow: 1;
      }

      > #editor-palette {
        display: flex;
        flex-direction: row;
        height: 100%;
        background-color: darken($theme-bg, 3);
        overflow-y: auto;

        > .palette-panel {
          display: flex;
          flex-direction: column;
          width: 48px;
          height: 100%;
        }
      }
    }

    > #editor-statusbar {
      height: 32px;
      line-height: 32px;
      background-color: darken($theme-bg, 3);
      display: flex;

      > #editor-zoom {
        display: flex;

        .icon {
          width: 32px;
          height: 32px;
          display: flex;
          align-items: center;
          justify-content: center;

          &:hover,
          &.zoom-sel {
            background-color: lighten($theme-bg, 3);
          }
        }

        > #zoom-value {
          margin: 0 4px;
        }
      }
    }
  }

  .flex-gap {
    flex-grow: 1;
  }

  .tool-button {
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;

    &:hover,
    &.tool-button-sel {
      background-color: lighten($theme-bg, 3);
    }
  }

  .palette-button {
    width: 32px;
    height: 32px;
    display: flex;
    padding: 8px;

    &.palette-selected {
      background-color: lighten($theme-bg, 3);
    }

    > .palette-button-blob {
      width: 32px;
      height: 32px;
      margin: auto;

      &[data-transparent] {
        background-color: #f2f2f2;
        $grey: #cdcdcd;
        background-image: linear-gradient(45deg, $grey 25%, transparent 25%, transparent 75%, $grey 75%, $grey),
          linear-gradient(45deg, $grey 25%, transparent 25%, transparent 75%, $grey 75%, $grey);
        background-size: 32px 32px;
        background-position: 0 0, 16px 16px;
      }
    }
  }
</style>
