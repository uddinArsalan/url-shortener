<script lang="ts">
  import { userStore } from "$lib/store/userStore";
  import { fetchUserUrls } from "$lib/api/url";
  import type { UrlType } from "../../types";
  import { onMount } from "svelte";
  import { API_BASE_URL } from "../../constants";
  let userUrls = $state<UrlType[]>([]);
  let hasMoreUrls = $state(false);
  let isLoading = $state(false);
  let nextCursor = $state<string>("");
  onMount(async () => {
    try {
      isLoading = true;
      const urls = await fetchUserUrls(nextCursor);
      userUrls = urls.urls;
      hasMoreUrls = urls.pagination.has_more;
      nextCursor = urls.pagination.next_cursor;
    } catch (error) {
      console.error("Error fetching user URLs:", error);
    } finally {
      isLoading = false;
    }
  });

  function loadMoreUrls() {
    isLoading = true;
    fetchUserUrls(nextCursor)
      .then((urls) => {
        userUrls = [...userUrls, ...urls.urls];
        hasMoreUrls = urls.pagination.has_more;
        nextCursor = urls.pagination.next_cursor;
      })
      .catch((error) => {
        console.error("Error fetching more URLs:", error);
      })
      .finally(() => {
        isLoading = false;
      });
  }
</script>

<h1>Dashboard</h1>
<p>Welcome to the dashboard! {$userStore.user.username}</p>
<div>
  <h2>Your URLs</h2>
  {#if isLoading}
    <p>Loading...</p>
  {:else if !userUrls || userUrls.length === 0}
    <p>No URLs found.</p>
  {:else}
    <ul>
      {#each userUrls as url}
        <li>
          <a href={`${API_BASE_URL}/${url.shortcode}`} target="_blank">
            Link
          </a>
          <span> - </span>
          <span>{url.original_url}</span>
          <span> - </span>
          <span>{url.created_at}</span>
        </li>
      {/each}
    </ul>
    {#if hasMoreUrls}
      <button onclick={loadMoreUrls}>Load more</button>
    {/if}
  {/if}
</div>
