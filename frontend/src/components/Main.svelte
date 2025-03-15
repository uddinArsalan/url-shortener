<script lang="ts">
  import { shortenUrl } from "../api/shorten";
  let originalUrl = "";
  let shortUrl: string | null;
  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    shortUrl = await shortenUrl(originalUrl);
  }
</script>

<main class="w-1/2 border mx-auto">
  <form
    onsubmit={handleSubmit}
    class="flex flex-col p-6 min-h-56 justify-between"
  >
    <label for="url" class="text-xl font-semibold underline text-center"
      >Add Url to shorten</label
    >
    <input
      type="text"
      id="url"
      name="url"
      placeholder="Add Url to shorten..."
      class="p-4 border"
      bind:value={originalUrl}
    />
    <button class="bg-blue-600 text-white p-2">Shorten Url</button>
    <div class="font-semibold">
      Your Shortened Url : <span>{shortUrl ? shortUrl : ""}</span>
    </div>
  </form>
</main>
