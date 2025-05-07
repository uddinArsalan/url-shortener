import { writable } from "svelte/store";
import type { UserStoreType } from "../../types";

export const userStore = writable<UserStoreType>({
  user: null,
  isLoggedIn: false,
});
