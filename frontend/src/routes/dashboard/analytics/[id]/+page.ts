import {
  getAnalyticsData,
  getHourlyClicksData,
  getCountryClicksData,
  getBrowserClicksData,
  getCityClicksData,
  getDeviceClicksData,
  getReferrerClicksData,
} from "$lib/api/analytics";
import type {
  BrowserClick,
  CityClick,
  ClickAnalyticsResponseType,
  CountryClick,
  DeviceClick,
  HourlyClick,
  ReferrerClick,
} from "../../../../types";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params }) => {
  try {
    const urlId = params.id;
    const startDate = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000);
    const endDate = new Date();
    const [
      analyticsData,
      hourlyClicksData,
      countryData,
      browserAnalyticsData,
      deviceAnalyticsData,
      referrerAnalyticsData,
      cityAnalyticsData,
    ] = await Promise.all([
      getAnalyticsData(urlId),
      getHourlyClicksData(urlId, startDate, endDate),
      getCountryClicksData(urlId),
      getBrowserClicksData(urlId),
      getDeviceClicksData(urlId),
      getReferrerClicksData(urlId),
      getCityClicksData(urlId),
    ]);
    return {
      analyticsData,
      hourlyClicksData,
      countryData,
      browserAnalyticsData,
      deviceAnalyticsData,
      referrerAnalyticsData,
      cityAnalyticsData,
    };
  } catch (error) {
    console.error("LOAD ERROR:", error);
    throw error;
  }
};

export interface AnalyticsType {
  analyticsData: ClickAnalyticsResponseType;
  hourlyClicksData: HourlyClick[];
  countryData: CountryClick[];
  browserAnalyticsData: BrowserClick[];
  deviceAnalyticsData: DeviceClick[];
  referrerAnalyticsData: ReferrerClick[];
  cityAnalyticsData: CityClick[];
}
