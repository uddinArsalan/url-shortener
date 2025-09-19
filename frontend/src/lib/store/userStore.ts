import { writable } from "svelte/store";
import type { UserStoreType } from "../../types";

export const userStore = writable<UserStoreType>({
  user: {
    id: 0,
    username: "",
    email: "",
    created_at: "",
  },
  isLoggedIn: false,
  isLoading: false,
});
