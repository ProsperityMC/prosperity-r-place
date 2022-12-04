import {writable} from "svelte/store";

export interface LoginStore {
  access: string;
  refresh: string;
}

export const loginStore = writable<LoginStore>(
  (() => {
    try {
      return JSON.parse(localStorage.getItem("login-session"));
    } catch (_) {
      return null;
    }
  })(),
);

loginStore.subscribe(value => localStorage.setItem("login-session", JSON.stringify(value)));

export interface DiscordMember {
  avatar: string | null;
  nick: string | null;
  user: DiscordUser;
}

export interface DiscordUser {
  id: string;
  username: string;
  avatar: string | null;
  discriminator: string;
}

export const profileStore = writable<DiscordMember>(
  (() => {
    try {
      return JSON.parse(localStorage.getItem("login-profile"));
    } catch (_) {
      return null;
    }
  })(),
);

profileStore.subscribe(value => localStorage.setItem("login-profile", JSON.stringify(value)));
