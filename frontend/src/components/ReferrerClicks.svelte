<script lang="ts">
  import {
    Chart,
    Legend,
    Title,
    BarController,
    ArcElement,
    BarElement,
  } from "chart.js";
  import type { ChartConfiguration } from "chart.js";
  import { onMount } from "svelte";
  import type { ReferrerClick } from "../types";
  Chart.register(Title, Legend, BarController, ArcElement, BarElement);

  let props = $props();
  let { data }: { data: ReferrerClick[] } = props;
  let labels = data.map((d) => d.referrer);
  let clicks = data.map((d) => d.count);
  let chartData = {
    labels,
    datasets: [
      {
        label: "Clicks by Referrer",
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
  const config: ChartConfiguration<"bar", number[], string> = {
    type: "bar",
    data: chartData,
    options: {
      responsive: true,
      scales: {
        y: {
          beginAtZero: true,
        },
      },
      plugins: {
        legend: {
          position: "top",
        },
        title: {
          display: true,
          text: "Referrer Wise Clicks",
        },
      },
    },
  };
  onMount(() => {
    const ctx = document.getElementById("referrer") as HTMLCanvasElement;
    if (ctx) {
      new Chart(ctx, config);
    }
  });
</script>

<div class="bg-white rounded-lg shadow p-6">
  <h2 class="text-xl font-semibold mb-4">Clicks by Referrer</h2>
  <canvas id="referrer"></canvas>
</div>
