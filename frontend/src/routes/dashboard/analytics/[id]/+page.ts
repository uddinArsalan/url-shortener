import { getAnalyticsData } from "$lib/api/analytics";
import type { ClickAnalyticsResponseType } from "../../../../types";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params }) => {
  console.log(params.id)
  const analyticsData = await getAnalyticsData(params.id as unknown as number);
  return {
    analyticsData,
  };
};

export interface AnalyticsType {
    analyticsData : ClickAnalyticsResponseType
}