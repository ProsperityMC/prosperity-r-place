<script lang="ts">
  import {xlink_attr} from "svelte/internal";
  import Editor from "./components/Editor.svelte";
  import {getEnv} from "./utils/env";

  let doc;

  async function getDocList() {
    let f = await fetch(getEnv("API_URL") + "/docs");
    return await f.json();
  }

  const t = new Date().getTime();
</script>

<main>
  {#if doc}
    <Editor {doc} on:exit={_ => (doc = undefined)} />
  {:else}
    <h1>Prosperity r/place</h1>
    {#await getDocList()}
      Loading...
    {:then docs}
      <div class="docs-view">
        {#each docs as x}
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <div class="doc-item" on:click|preventDefault={() => (doc = x)}>
            <img src="{getEnv('API_URL')}/doc/{x.name}?raw=image&_={t}" alt="{x.name} image" />
            <h1>{x.name}</h1>
            <h2>{x.width} x {x.height}</h2>
          </div>
        {/each}
      </div>
    {:catch err}
      <div>{err}</div>
    {/await}
  {/if}
</main>

<style lang="scss">
  @import "assets/theme.scss";

  main > h1 {
    color: #d2d2d2;
    margin: 0;
    padding: 16px;
    line-height: normal;
  }

  .docs-view {
    display: flex;
    flex-wrap: wrap;
    padding: 16px;

    > .doc-item {
      min-width: 220px;
      min-height: 220px;
      border-radius: 16px;
      background-color: lighten($theme-bg, 3);
      box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
      overflow: hidden;
      cursor: pointer;
      transition: background-color 100ms;

      &:hover {
        background-color: lighten($theme-bg, 5);

        > img {
          background-color: lighten($theme-bg, 8);
        }
      }

      > img {
        width: calc(100% - 16px);
        margin: 8px;
        box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
        image-rendering: pixelated;
        border-radius: 8px;
        background-color: lighten($theme-bg, 6);
      }

      > h1 {
        color: #d2d2d2;
        margin: 4px 16px;
      }

      > h2 {
        color: #828282;
        margin: 4px 16px 8px;
      }
    }
  }
</style>
