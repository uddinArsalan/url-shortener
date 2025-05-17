<script lang="ts">
  import type { ClickAnalyticsResponseType } from "../../../../types";
  import {
    Monitor,
    Smartphone,
    Bot,
    Chrome,
    MapPin,
    Clock,
    Link2,
    MousePointerClick,
    Users,
  } from "@lucide/svelte";
  import type { AnalyticsType } from "./+page";

  export let data : AnalyticsType;

  const getDeviceIcon = (device: string) => {
    switch (device) {
      case "desktop":
        return Monitor;
      case "mobile":
        return Smartphone;
      case "bot":
        return Bot;
      default:
        return Monitor;
    }
  };
</script>

<div class="min-h-screen bg-gray-50 p-6">
  <div class="max-w-7xl mx-auto">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
      <div class="bg-white rounded-lg shadow p-6">
        <div class="flex items-center gap-4">
          <MousePointerClick class="w-8 h-8 text-blue-600" />
          <div>
            <p class="text-gray-600">Total Clicks</p>
            <h3 class="text-2xl font-bold">{data.analyticsData.total_clicks}</h3>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow p-6">
        <div class="flex items-center gap-4">
          <Users class="w-8 h-8 text-green-600" />
          <div>
            <p class="text-gray-600">Unique Clicks</p>
            <h3 class="text-2xl font-bold">{data.analyticsData.unique_clicks}</h3>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-lg shadow overflow-hidden">
      <h2 class="text-xl font-semibold p-6 border-b">Click Analytics</h2>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50">
            <tr>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase"
                >Device</th
              >
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase"
                >Browser</th
              >
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase"
                >Location</th
              >
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase"
                >Referer</th
              >
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase"
                >Time</th
              >
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            {#each data.analyticsData.click_analytics as click}
              <tr class="hover:bg-gray-50">
                <td class="px-6 py-4">
                  <div class="flex items-center gap-2">
                    <svelte:component
                      this={getDeviceIcon(click.device)}
                      class="w-4 h-4"
                    />
                    <span class="capitalize">{click.device}</span>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <div class="flex items-center gap-2">
                    <Chrome class="w-4 h-4" />
                    <span>{click.browser}</span>
                    <span class="text-gray-500">({click.os})</span>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <div class="flex items-center gap-2">
                    <MapPin class="w-4 h-4" />
                    <span>{click.city}, {click.country}</span>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <div class="flex items-center gap-2">
                    <Link2 class="w-4 h-4" />
                    <span class="truncate max-w-xs"
                      >{click.referrer || "Direct"}</span
                    >
                  </div>
                </td>
                <td class="px-6 py-4">
                  <div class="flex items-center gap-2">
                    <Clock class="w-4 h-4" />
                    <span>{new Date(click.timestamp).toLocaleString()}</span>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  </div>
</div>
