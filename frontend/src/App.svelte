<script lang="ts">
  import Editor from "./components/Editor.svelte";
  import {loginStore, profileStore} from "./stores/login";
  import {getEnv} from "./utils/env";

  let doc;

  async function getDocList() {
    let f = await fetch(getEnv("API_URL") + "/docs");
    return await f.json();
  }

  function clickLoginButton() {
    popupCenterScreen(`${getEnv("API_URL")}/login`, "Login with Discord", 800, 800, false);
  }

  function handleMessage(event) {
    console.log(event);
    if (event.origin !== getEnv("API_ORIGIN")) return;
    if (isObject(event.data)) {
      console.log(event.data);
      if (isObject(event.data.member)) {
        let d = Object.assign({user: {id: null, username: null, discriminator: null}}, event.data.member);
        if (d.username === null || d.id === null || d.discriminator === null) {
          alert("Failed to log user in: the login data is structured correctly but probably corrupted");
          return;
        }
        loginStore.set(event.data.token);
        profileStore.set(event.data.member);
        return;
      }
    }
    alert("Failed to log user in: the login data was probably corrupted");
  }

  function isObject(obj) {
    return obj != null && obj.constructor.name === "Object";
  }

  function popupCenterScreen(url, title, w, h, focus) {
    const top = (screen.availHeight - h) / 4,
      left = (screen.availWidth - w) / 2;
    const popup = openWindow(url, title, `scrollbars=yes,width=${w},height=${h},top=${top},left=${left}`);
    if (focus === true && window.focus) popup.focus();
    return popup;
  }

  function openWindow(url, winnm, options) {
    var wTop = firstAvailableValue([window.screen.availTop, window.screenY, window.screenTop, 0]);
    var wLeft = firstAvailableValue([window.screen.availLeft, window.screenX, window.screenLeft, 0]);
    var top = 0,
      left = 0;
    var result;
    let w;
    if ((result = /top=(\d+)/g.exec(options))) top = parseInt(result[1]);
    if ((result = /left=(\d+)/g.exec(options))) left = parseInt(result[1]);
    if (options) {
      options = options.replace("top=" + top, "top=" + (parseInt(top) + wTop));
      options = options.replace("left=" + left, "left=" + (parseInt(left) + wLeft));
      w = window.open(url, winnm, options);
    } else w = window.open(url, winnm);
    return w;
  }

  function firstAvailableValue(arr) {
    for (var i = 0; i < arr.length; i++) if (typeof arr[i] != "undefined") return arr[i];
  }

  function generateAvatarUrl(id: string, discriminator: string, hash: string, guild: boolean): string {
    if (hash === "") {
      let i = parseInt(discriminator);
      return `https://cdn.discordapp.com/embed/avatars/${i % 5}.png?size=512`;
    } else if (hash.startsWith("a_")) {
      if (guild) return `https://cdn.discordapp.com/guilds/${getEnv("API_GUILD")}/users/${id}/avatars/${hash}.gif?size=512`;
      else return `https://cdn.discordapp.com/avatars/${id}/${hash}.gif?size=512`;
    } else {
      if (guild) return `https://cdn.discordapp.com/guilds/${getEnv("API_GUILD")}/users/${id}/avatars/${hash}.png?size=512`;
      else return `https://cdn.discordapp.com/avatars/${id}/${hash}.png?size=512`;
    }
  }

  const t = new Date().getTime();
</script>

<svelte:head>
  <script src="zlib.js"></script>
  <script src="png.js"></script>
</svelte:head>

<svelte:window on:message={handleMessage} />

<main>
  {#if doc}
    <Editor {doc} on:exit={_ => (doc = undefined)} />
  {:else}
    <header>
      <h1>Prosperity r/place</h1>
      <div class="flex-gap" />
      {#if $profileStore}
        {#if $profileStore.avatar}
          <img
            src={generateAvatarUrl($profileStore.user.id, $profileStore.user.discriminator, $profileStore.avatar, true)}
            width="32"
            height="32"
            alt="Discord Avatar"
          />
        {:else}
          <img
            src={generateAvatarUrl($profileStore.user.id, $profileStore.user.discriminator, $profileStore.user.avatar ?? "", false)}
            width="32"
            height="32"
            alt="Discord Avatar"
          />
        {/if}
        <div>
          <div class="profile-user">{$profileStore.user.username}</div>
          {#if $profileStore.nick}
            <div class="profile-nick">{$profileStore.nick}</div>
          {/if}
        </div>
        <button class="yellow-button" on:click={() => (loginStore.set(undefined), profileStore.set(undefined))}>Logout</button>
      {:else}
        <button class="yellow-button" on:click={clickLoginButton}>Login</button>
      {/if}
    </header>
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

  header {
    display: flex;
    padding: 0 32px;
    align-items: center;
    gap: 16px;
    background-color: lighten($theme-bg, 3);
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);

    > .flex-gap {
      flex-grow: 1;
    }

    img {
      border-radius: 50%;
    }
  }

  .yellow-button {
    background-color: rgb(255, 193, 11);
    box-shadow: rgba(0, 0, 0, 0) 0px 0px 0px 0px, rgba(0, 0, 0, 0) 0px 0px 0px 0px, rgba(255, 193, 11, 0.2) 0px 4px 6px -1px,
      rgba(255, 193, 11, 0.2) 0px 2px 4px -2px;
    box-sizing: border-box;
    color: rgb(36, 37, 39);
    cursor: pointer;
    font-size: 16px;
    font-weight: 700;
    line-height: 24px;
    padding: 8px 24px;
    border-radius: 0.375rem;
  }

  .docs-view {
    display: flex;
    flex-wrap: wrap;
    padding: 16px;
    gap: 16px;

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
