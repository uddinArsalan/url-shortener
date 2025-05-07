<script lang="ts">
  import { API_BASE_URL } from "../constants";
  import { Menu, X, BarChart3, User, LogIn } from "@lucide/svelte";
  import { userStore } from "$lib/store/userStore";
  import { fetchUser } from "$lib/api/auth";
  import { onMount } from "svelte";
  let isMobileMenuOpen = $state(false);

  async function setUserStore() {
    try {
      const userDetails = await fetchUser();
      userStore.update((store) => {
        store.user = userDetails;
        store.isLoggedIn = true;
        return store;
      });
    } catch (error) {
      console.error("Error fetching user:", error);
    }
  }

  onMount(setUserStore);

  function toggleMobileMenu() {
    isMobileMenuOpen = !isMobileMenuOpen;
  }
  async function handleLogin() {
    window.location.href = `${API_BASE_URL}/auth/login`;
    setUserStore();
  }
  $inspect($userStore.user);
</script>

<nav class="bg-white border-b shadow-sm">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="flex justify-between h-16">
      <div class="flex items-center">
        <a href="/" class="flex-shrink-0 flex items-center">
          <span class="text-2xl font-bold text-blue-600 font-mono"
            >URL<span class="text-purple-600">Short</span></span
          >
        </a>
        <div class="hidden md:block ml-4">
          <p class="text-gray-500 italic text-lg">
            A distributed URL shortener
          </p>
        </div>
      </div>

      <div class="hidden md:flex items-center">
        {#if $userStore.isLoggedIn}
          <div class="flex space-x-4 mr-4">
            <a
              href="/dashboard"
              class="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md flex items-center transition-colors"
              aria-label="Go to dashboard"
            >
              <BarChart3 size={18} class="mr-1" />
              <span>Dashboard</span>
            </a>
            <a
              href="/analytics"
              class="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md flex items-center transition-colors"
              aria-label="View analytics"
            >
              <BarChart3 size={18} class="mr-1" />
              <span>Analytics</span>
            </a>
          </div>
          <div class="flex items-center space-x-2">
            <a
              href="/profile"
              class="text-gray-600 hover:text-blue-600 p-1 rounded-full"
              aria-label="View profile"
            >
              <User size={24} />
              <div>{$userStore.user?.Username}</div>
            </a>
          </div>
        {:else}
          <div class="flex items-center space-x-4">
            <button
              class="text-gray-700 font-medium py-2 px-4 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors flex items-center focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              aria-label="Log in"
              onclick={handleLogin}
            >
              <LogIn size={18} class="mr-1" />
              Login
            </button>
          </div>
        {/if}
      </div>

      <div class="flex items-center md:hidden">
        <button
          type="button"
          class="inline-flex items-center justify-center p-2 rounded-md text-gray-500 hover:text-gray-700 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500"
          aria-controls="mobile-menu"
          aria-expanded={isMobileMenuOpen}
          onclick={toggleMobileMenu}
        >
          <span class="sr-only"
            >{isMobileMenuOpen ? "Close menu" : "Open menu"}</span
          >
          {#if isMobileMenuOpen}
            <X size={24} />
          {:else}
            <Menu size={24} />
          {/if}
        </button>
      </div>
    </div>
  </div>

  {#if isMobileMenuOpen}
    <div class="md:hidden" id="mobile-menu">
      <div class="border-t pt-2 pb-3 space-y-1">
        <p class="px-3 py-1 text-gray-500 italic text-lg">
          A distributed URL shortener
        </p>

        {#if $userStore.isLoggedIn}
          <a
            href="/dashboard"
            class="text-gray-600 hover:bg-gray-50 hover:text-blue-600 px-3 py-2 rounded-md text-base font-medium flex items-center"
          >
            <BarChart3 size={18} class="mr-2" />
            Dashboard
          </a>
          <a
            href="/analytics"
            class="text-gray-600 hover:bg-gray-50 hover:text-blue-600 px-3 py-2 rounded-md text-base font-medium flex items-center"
          >
            <BarChart3 size={18} class="mr-2" />
            Analytics
          </a>
          <a
            href="/profile"
            class="text-gray-600 hover:bg-gray-50 hover:text-blue-600 px-3 py-2 rounded-md text-base font-medium flex items-center"
          >
            <User size={18} class="mr-2" />
            Profile
          </a>
        {:else}
          <div class="px-2 pt-2 pb-3 space-y-2">
            <button
              class="w-full text-gray-700 font-medium py-2 px-4 border border-gray-300 rounded-md hover:bg-gray-50 flex items-center justify-center"
              onclick={handleLogin}
            >
              <LogIn size={18} class="mr-1" />
              Login
            </button>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</nav>
