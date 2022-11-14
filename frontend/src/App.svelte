<script lang="ts">
  import Editor from "./components/Editor.svelte";
  import {getEnv} from "./utils/env";

  let doc;

  async function getDocList() {
    let f = await fetch(getEnv("API_URL") + "/docs");
    return await f.json();
  }
</script>

<main>
  {#if doc}
    <Editor {doc} />
  {:else}
    <h1>Prosperity r/place</h1>
    {#await getDocList()}
      Loading...
    {:then docs}
      <ul>
        {#each docs as x}
          <li><button on:click={() => (doc = x)}>{x.name} - {x.width}x{x.height}</button></li>
        {/each}
      </ul>
    {:catch err}
      <div>{err}</div>
    {/await}
  {/if}
</main>
