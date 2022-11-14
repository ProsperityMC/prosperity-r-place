<script lang="ts">
  import {onMount} from "svelte";
  import colourPalette from "~/assets/colours.json";
  import {getEnv} from "~/utils/env";
  import Document from "./Document.svelte";

  export let doc;
  let ws;
  let showMenu = true;
  let showPalette = true;

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
    <div class="tool-gap" />
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
          <div class="tool-button" data-icon="pan_tool" />
          <div class="tool-button" data-icon="edit" />
          <div class="tool-button" data-icon="format_color_fill" />
          <div class="tool-button" data-icon="shape_line" />
          <div class="tool-button" data-icon="select_all" />
          <div class="tool-button" data-icon="deselect" />
        </div>
        <div id="editor-shapes" />
      {/if}
      <div id="editor-doc">
        <Document {doc} />
      </div>
      {#if showPalette}
        <div id="editor-palette">
          {#each colourPalette as palette}
            <div class="palette-panel">
              {#each palette.options as option}
                <div class="palette-button" title="{option.name} {palette.name}">
                  <div class="palette-button-blob" style="background-color:{option.hex};" />
                </div>
              {/each}
            </div>
          {/each}
        </div>
      {/if}
    {:else}
      <div>Connecting to live editor</div>
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
      background-color: darken($theme-bg, 5);
      height: 48px;
    }

    > #editor-main {
      display: flex;
      flex-direction: row;
      flex-grow: 1;
      min-height: 0;

      > #editor-tools {
        width: 48px;
        height: 100%;
        background-color: darken($theme-bg, 5);
      }

      > #editor-doc {
        overflow: hidden;
        flex-grow: 1;
      }

      > #editor-palette {
        display: flex;
        flex-direction: row;
        height: 100%;
        background-color: darken($theme-bg, 5);
        overflow-y: auto;

        > .palette-panel {
          width: 48px;
          height: 100%;
        }
      }
    }
  }

  .tool-gap {
    flex-grow: 1;
  }

  .tool-button {
    width: 48px;
    height: 48px;
    display: flex;

    &::before {
      @include mso;
      width: 32px;
      height: 32px;
      display: block;
      content: attr(data-icon);
      margin: auto;
      font-size: 32px;
      line-height: 32px;
      white-space: pre;
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
