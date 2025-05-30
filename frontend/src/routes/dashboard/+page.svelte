<script lang="ts">
  import { userStore } from "$lib/store/userStore";
  import { logout } from "$lib/api/auth";
  import { fetchUserUrls } from "$lib/api/url";
  import type { UrlType } from "../../types";
  import { onMount } from "svelte";
  import { API_BASE_URL } from "../../constants";
  import {
    BarChart2,
    ChevronRight,
    Link2,
    Loader2,
    ChevronLeft,
  } from "@lucide/svelte";
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

<div class="min-h-screen bg-gray-100">
  <header class="bg-white shadow">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-5">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-2">
          <a
            href="/"
            class="flex items-center text-gray-700 hover:text-gray-900"
          >
            <ChevronLeft class="w-5 h-5 mt-1" />
          </a>
          <h1 class="text-2xl font-semibold text-gray-900">Your Links</h1>
        </div>

        <div class="flex items-center space-x-4">
          <span class="text-sm font-medium text-gray-800">
            {$userStore.user.username}
          </span>
          <button
            onclick={() => logout()}
            class="text-sm font-medium text-indigo-600 hover:text-indigo-800 transition-colors"
          >
            Sign out
          </button>
        </div>
      </div>
    </div>
  </header>

  <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    {#if isLoading && userUrls.length === 0}
      <div class="flex justify-center items-center h-64">
        <Loader2 class="h-8 w-8 text-indigo-600 animate-spin" />
      </div>
    {:else if userUrls.length === 0}
      <div class="text-center py-12 bg-white rounded-lg shadow">
        <Link2 class="mx-auto h-12 w-12 text-gray-400" />
        <h3 class="mt-2 text-lg font-medium text-gray-900">No links yet</h3>
        <p class="mt-1 text-sm text-gray-500">
          Get started by creating a new shortened URL.
        </p>
        <div class="mt-6">
          <a
            href="/"
            class="inline-flex items-center px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-md hover:bg-indigo-700"
          >
            Create Link
          </a>
        </div>
      </div>
    {:else}
      <div class="bg-white shadow rounded-lg overflow-hidden">
        <ul class="divide-y divide-gray-200">
          {#each userUrls as url}
            <li>
              <div class="block hover:bg-gray-50 transition-colors">
                <div
                  class="px-4 py-5 sm:px-6 flex items-center justify-between"
                >
                  <div class="flex items-center space-x-4">
                    <Link2 class="h-6 w-6 text-gray-400" />
                    <div class="min-w-0">
                      <a
                        href={`${API_BASE_URL}/url/${url.shortcode}`}
                        target="_blank"
                        class="text-sm font-medium text-indigo-600 truncate"
                      >
                        {API_BASE_URL}/url/{url.shortcode}
                      </a>
                      <p class="text-sm text-gray-500 truncate max-w-md">
                        {url.original_url}
                      </p>
                    </div>
                  </div>
                  <div class="flex items-center space-x-4">
                    <span class="text-sm text-gray-500">
                      {new Date(url.created_at).toLocaleDateString()}
                    </span>
                    <BarChart2 class="h-5 w-5 text-gray-400" />
                    <a href={`/dashboard/analytics/${url.id}`}
                      ><ChevronRight class="h-5 w-5 text-gray-400" /></a
                    >
                  </div>
                </div>
              </div>
            </li>
          {/each}
        </ul>
      </div>
      {#if hasMoreUrls}
        <div class="mt-6 flex justify-center">
          <button
            onclick={loadMoreUrls}
            disabled={isLoading}
            class="inline-flex items-center px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-md hover:bg-indigo-700 disabled:bg-indigo-400 disabled:cursor-not-allowed"
          >
            {#if isLoading}
              <Loader2 class="mr-2 h-5 w-5 animate-spin" />
            {/if}
            Load More
          </button>
        </div>
      {/if}
    {/if}
  </main>
</div>
