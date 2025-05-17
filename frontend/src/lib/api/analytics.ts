import axios from "$lib/axios";
import type { ClickAnalyticsResponseType } from "../../types";
export async function getAnalyticsData(urlId: number) {
  try {
    const res = await axios.get(`/analytics?urlId=${urlId}`);
    return res.data as ClickAnalyticsResponseType;
  } catch (error) {
    console.error("Error fetching analytics data:", error);
    throw error;
  }
}
