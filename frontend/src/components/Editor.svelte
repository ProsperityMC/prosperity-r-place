<script lang="ts">
  import {onMount} from "svelte";
  import colourPalette from "~/assets/colours.json";
  import {getEnv} from "~/utils/env";
  import Document from "./Document.svelte";

  export let doc;
  let ws;
  let showMenu = true;
  let showPalette = true;
  let menuSel = "pan";
  let shapeSel = "circle";
  let paletteSel = {name: "transparent", hex: "#ffffffff"};
  let scale = 1;

  async function connectToWebsocket() {
    let wsUrl = `${getEnv("API_URL").replace("https://", "wss://")}/doc/${doc.name}`;
    console.log(`Connecting to WS: ${wsUrl}`);
    let openWS = new WebSocket(wsUrl);
    openWS.onopen = function () {
      ws = openWS;
    };
  }

  onMount(async () => {
    console.log(`Starting editor for ${doc.name} - ${doc.width}x${doc.height}`);
    connectToWebsocket();
  });

  const menuButtons = [
    {key: "pan", icon: "pan_tool", select: true},
    {key: "pencil", icon: "edit", select: true},
    {key: "fill", icon: "format_color_fill", select: true},
    {key: "shape", icon: "shape_line", select: true},
    {key: "select", icon: "select_all", select: true},
    {key: "deselect", icon: "deselect", select: false},
  ];

  const shapeButtons = [
    {key: "circle", icon: "circle", select: true},
    {key: "triagle", icon: "change_history", select: true},
    {key: "square", icon: "square", select: true},
    {key: "pentagon", icon: "pentagon", select: true},
    {key: "hexagon", icon: "hexagon", select: true},
  ];
</script>

<div id="editor">
  <div id="editor-navbar">
    <div
      class="tool-button"
      data-icon="menu"
      on:click={() => {
        showMenu = !showMenu;
      }}
      on:keydown={() => {
        showMenu = !showMenu;
      }}
    />
    <div class="flex-gap" />
    <div
      class="tool-button"
      data-icon="palette"
      on:click={() => {
        showPalette = !showPalette;
      }}
      on:keydown={() => {
        showPalette = !showPalette;
      }}
    />
  </div>
  <div id="editor-main">
    {#if ws}
      {#if showMenu}
        <div id="editor-tools">
          {#each menuButtons as b (b.key)}
            <div
              class="tool-button {menuSel == b.key ? 'tool-button-sel' : ''}"
              data-icon={b.icon}
              on:click={() => {
                if (b.select) menuSel = b.key;
              }}
              on:keypress={() => {
                if (b.select) menuSel = b.key;
              }}
            />
          {/each}
        </div>
        {#if menuSel == "shape"}
          <div id="editor-shapes">
            {#each shapeButtons as b (b.key)}
              <div
                class="tool-button {menuSel == b.key ? 'tool-button-sel' : ''}"
                data-icon={b.icon}
                on:click={() => {
                  if (b.select) shapeSel = b.key;
                }}
                on:keypress={() => {
                  if (b.select) shapeSel = b.key;
                }}
              />
            {/each}
          </div>
        {/if}
      {/if}
      <div id="editor-doc">
        <Document {doc} {menuSel} {paletteSel} bind:scale />
      </div>
      {#if showPalette}
        <div id="editor-palette">
          {#each colourPalette as palette}
            <div class="palette-panel">
              {#each palette.options as option}
                <div
                  class="palette-button"
                  title="{option.name} {palette.name}"
                  on:click={() => {
                    paletteSel = option;
                  }}
                  on:keypress={() => {
                    paletteSel = option;
                  }}
                >
                  <div class="palette-button-blob" style="background-color:{option.hex};" />
                </div>
              {/each}
            </div>
          {/each}
        </div>
      {/if}
    {/if}
  </div>
  <div id="editor-statusbar">
    {#if ws}
      <div class="flex-gap" />
      <div id="editor-zoom">
        <div class="icon" data-icon="zoom_out" />
        <div id="zoom-value">{Math.floor(scale * 100)}%</div>
        <div class="icon" data-icon="zoom_in" />
        <div class="icon" data-icon="fit_screen" />
      </div>
    {:else}
      <div id="editor-connecting">Connecting to live editor</div>
    {/if}
  </div>
</div>

<style lang="scss">
  @import "../assets/material-symbols.scss";
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
        width: 48px;
        background-color: darken($theme-bg, 3);
      }

      > #editor-shapes {
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
          @include mso;
          width: 32px;
          height: 32px;
          display: flex;

          &::before {
            @include mso;
            width: 24px;
            height: 24px;
            display: block;
            content: attr(data-icon);
            margin: auto;
            font-size: 24px;
            line-height: 24px;
            white-space: pre;
          }

          &:hover {
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

    &::before {
      @include mso;
      width: 24px;
      height: 24px;
      display: block;
      content: attr(data-icon);
      margin: auto;
      font-size: 24px;
      line-height: 24px;
      white-space: pre;
    }

    &:hover,
    &.tool-button-sel {
      background-color: lighten($theme-bg, 3);
    }
  }

  .palette-button {
    width: 48px;
    height: 48px;
    display: flex;

    > .palette-button-blob {
      width: 32px;
      height: 32px;
      margin: auto;
    }
  }
</style>
