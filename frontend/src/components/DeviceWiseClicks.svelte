<script lang="ts">
  import {
    Chart,
    Legend,
    Title,
    DoughnutController,
    ArcElement,
  } from "chart.js";
  import type { ChartConfiguration } from "chart.js";
  import { onMount } from "svelte";
  import type { DeviceClick } from "../types";
  Chart.register(Title, Legend, DoughnutController, ArcElement);

  let props = $props();
  let { data }: { data: DeviceClick[] } = props;
  let labels = data.map((d) => d.device);
  let clicks = data.map((d) => d.count);
  let chartData = {
    labels,
    datasets: [
      {
        label: "Clicks by Device Type",
        data: clicks,
        backgroundColor: [
          "#FF6384",
          "#36A2EB",
          "#FFCE56",
          "#4BC0C0",
          "#9966FF",
          "#FF9F40",
        ],
      },
    ],
  };
  const config: ChartConfiguration<"doughnut", number[], string> = {
    type: "doughnut",
    data: chartData,
    options: {
      responsive: true,
      plugins: {
        legend: {
          position: "top",
        },
        title: {
          display: true,
          text: "DeviceWise Clicks",
        },
      },
    },
  };
  onMount(() => {
    const ctx = document.getElementById("device") as HTMLCanvasElement;
    if (ctx) {
      new Chart(ctx, config);
    }
  });
</script>

<div class="bg-white rounded-lg shadow p-6">
  <h2 class="text-xl font-semibold mb-4">Clicks by Device</h2>
  <canvas id="device" class="w-full h-[400px]"></canvas>
</div>
