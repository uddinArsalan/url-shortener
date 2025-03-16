<script lang="ts">
  import { shortenUrl } from "$lib/api/shorten";
  import {
    ExternalLink,
    Link,
    Loader2,
    AlertCircle,
    Clipboard,
  } from "@lucide/svelte";
  let originalUrl = "";
  let error: string | null;
  let loading: boolean;
  let shortUrl: string | null;
  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!originalUrl.trim()) {
      error = "Please enter a valid URL.";
      return;
    }
    error = null;
    loading = true;
    try {
      shortUrl = await shortenUrl(originalUrl);
    } catch (error) {
      error = "Failed to shorten URL. Please try again.";
    } finally {
      loading = false;
    }
  }

  function copyToClipboard() {
    if (shortUrl) {
      navigator.clipboard.writeText(shortUrl);
    }
  }
</script>

<main class="w-full max-w-md mx-auto my-8">
  <div
    class="bg-white rounded-lg shadow-lg overflow-hidden border border-gray-200"
  >
    <div class="bg-gradient-to-r from-blue-600 to-purple-600 p-4">
      <h2 class="text-xl font-bold text-white text-center">URL Shortener</h2>
    </div>

    <form
      on:submit|preventDefault={handleSubmit}
      class="flex flex-col p-6 space-y-6"
    >
      <div class="space-y-2">
        <label for="url" class="text-gray-700 font-medium block"
          >Enter your long URL</label
        >
        <div class="relative">
          <input
            type="text"
            id="url"
            name="url"
            placeholder="https://example.com/very/long/url/to/shorten..."
            class="w-full p-3 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all pl-10"
            bind:value={originalUrl}
            aria-label="URL to shorten"
          />
          <div class="absolute left-3 top-3.5 text-gray-400">
            <Link size={20} />
          </div>
        </div>
      </div>

      <button
        type="submit"
        class={`w-full py-3 px-4 rounded-md font-medium text-white transition-all transform hover:scale-[1.02] focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 ${loading ? "bg-blue-400 cursor-not-allowed" : "bg-blue-600 hover:bg-blue-700 shadow-md"}`}
        disabled={loading}
        aria-label="Shorten URL"
      >
        {#if loading}
          <span class="flex items-center justify-center">
            <Loader2 class="animate-spin mr-2" size={18} />
            Shortening...
          </span>
        {:else}
          Shorten URL
        {/if}
      </button>

      {#if error}
        <div
          class="p-3 bg-red-50 border border-red-200 rounded-md text-red-600 text-sm flex items-start"
          role="alert"
        >
          <AlertCircle class="mr-2 flex-shrink-0" size={20} />
          <span>{error}</span>
        </div>
      {/if}

      {#if shortUrl}
        <div class="mt-4 p-4 bg-blue-50 border border-blue-100 rounded-md">
          <p class="text-gray-700 font-medium mb-2">
            Your shortened URL is ready!
          </p>
          <div class="flex items-center">
            <a
              href={shortUrl}
              target="_blank"
              class="text-blue-600 font-medium hover:text-blue-800 truncate mr-2 flex-grow flex items-center"
              aria-label="Open shortened URL in new tab"
            >
              {shortUrl}
              <ExternalLink class="ml-1 inline" size={16} />
            </a>
            <button
              class="p-2 bg-gray-100 hover:bg-gray-200 rounded-md transition-colors"
              on:click={copyToClipboard}
              title="Copy to clipboard"
              aria-label="Copy shortened URL to clipboard"
            >
              <Clipboard size={18} class="text-gray-600" />
            </button>
          </div>
        </div>
      {/if}
    </form>
  </div>
</main>
