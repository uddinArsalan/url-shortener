<script lang="ts">
  import {
    Chrome,
    MapPin,
    Clock,
    Link2,
    MousePointerClick,
    Users,
  } from "@lucide/svelte";
  import type { AnalyticsType } from "./+page";
  import { getDeviceIcon } from "$lib/utils";
  import CountryClicks from "../../../../components/CountryClicks.svelte";
  import HourlyClicks from "../../../../components/HourlyClicks.svelte";
  import BrowserWiseClicks from "../../../../components/BrowserWiseClicks.svelte";
  import DeviceWiseClicks from "../../../../components/DeviceWiseClicks.svelte";
  import CityWiseClicks from "../../../../components/CityWiseClicks.svelte";
  import ReferrerClicks from "../../../../components/ReferrerClicks.svelte";
  export let data: AnalyticsType;
</script>

<div class="min-h-screen bg-gray-50 p-6">
  <div class="max-w-7xl mx-auto">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
      <div class="bg-white rounded-lg shadow p-6">
        <div class="flex items-center gap-4">
          <MousePointerClick class="w-8 h-8 text-blue-600" />
          <div>
            <p class="text-gray-600">Total Clicks</p>
            <h3 class="text-2xl font-bold">
              {data.analyticsData.total_clicks}
            </h3>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow p-6">
        <div class="flex items-center gap-4">
          <Users class="w-8 h-8 text-green-600" />
          <div>
            <p class="text-gray-600">Unique Clicks</p>
            <h3 class="text-2xl font-bold">
              {data.analyticsData.unique_clicks}
            </h3>
          </div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
      <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold mb-4">Clicks by Country</h2>
        <CountryClicks data={data.countryData} />
      </div>
      <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold mb-4">Hourly Clicks</h2>
        <HourlyClicks data={data.hourlyClicksData} />
      </div>
      <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold mb-4">Clicks by Browser</h2>
        <BrowserWiseClicks data={data.browserAnalyticsData} />
      </div>
      <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold mb-4">Clicks by Device</h2>
        <DeviceWiseClicks data={data.deviceAnalyticsData} />
      </div>
      <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold mb-4">Clicks by City</h2>
        <CityWiseClicks data={data.cityAnalyticsData} />
      </div>
      <div class="bg-white rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold mb-4">Referrer Clicks</h2>
        <ReferrerClicks data={data.referrerAnalyticsData} />
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
                >OS</th
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
                    <span class="capitalize">{click.device || "Desktop"}</span>
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
                    <Link2 class="w-4 h-4" />
                    <span class="truncate max-w-xs"
                      >{click.os}</span
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
