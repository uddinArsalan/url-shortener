<script>
  let { children } = $props();
  import "../app.css";
  import { onMount } from "svelte";
  import { userStore } from "$lib/store/userStore";
  import { fetchUser } from "$lib/api/auth";

  onMount(async () => {
    userStore.update((s) => ({ ...s, isLoading: true }));
    try {
      const user = await fetchUser();
      userStore.update((s) => ({
        ...s,
        user,
        isLoggedIn: true,
        isLoading: false,
      }));
    } catch {
      userStore.update((s) => ({
        ...s,
        user: null,
        isLoggedIn: false,
        isLoading: false,
      }));
    }
  });
</script>

{@render children()}
