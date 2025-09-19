import { onMount } from "svelte";
import { fetchUser } from "$lib/api/auth";
import { userStore } from "$lib/store/userStore";
// import { redirect } from "@sveltejs/kit";
// import { API_BASE_URL } from "../constants";

onMount(async () => {
  try {
    const userDetails = await fetchUser();
    userStore.update((s) => ({ ...s, user: userDetails, isLoggedIn: true }));
  } catch (err) {
    userStore.update((s) => ({ ...s, user: null, isLoggedIn: false }));
  }
});
