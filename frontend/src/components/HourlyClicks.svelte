<script lang="ts">
  import type { HourlyClick } from "../types";
  import type { ChartConfiguration } from "chart.js";
  import {format} from "date-fns"
  import { Chart,LineController,CategoryScale,LinearScale,PointElement,LineElement } from "chart.js";
  import { onMount } from "svelte";

  Chart.register(LineController,CategoryScale,LinearScale,PointElement,LineElement)

  let props = $props();
  let { data }: { data: HourlyClick[] } = props;
  const labels = data.map((d) => format(d.hour, "MMM d, h aa"));
  const clicks = data.map((d) => d.count);

  let hourlyClicksData = {
    labels,
    datasets: [
      {
        label: "Hourly Clicks",
        data: clicks,
      },
    ],
  };

  const config: ChartConfiguration = {
    data: hourlyClicksData,
    type: "line",
    options: {},
  };
  onMount(() => {
    const ctx = document.getElementById("lineChart") as HTMLCanvasElement;
    if (ctx) {
      new Chart(ctx, config);
    }
  });
</script>

<div class="bg-white rounded-lg shadow p-6">
  <h2 class="text-xl font-semibold mb-4">Clicks per Hour</h2>
  <canvas id="lineChart"></canvas>
</div>
