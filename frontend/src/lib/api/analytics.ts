import axios from "$lib/axios";
import type {
  BrowserClick,
  CityClick,
  ClickAnalyticsResponseType,
  CountryClick,
  DeviceClick,
  HourlyClick,
  ReferrerClick,
} from "../../types";

export async function getAnalyticsData(urlId: string) {
  try {
    const res = await axios.get(`/analytics?urlId=${urlId}`);
    return res.data as ClickAnalyticsResponseType;
  } catch (error) {
    console.error("Error fetching analytics data:", error);
    throw error;
  }
}
export async function getHourlyClicksData(urlId: string, from: Date, to: Date) {
  try {
    const res = await axios.get(
      `/analytics/${urlId}/hourly?from=${from.toISOString()}&to=${to.toISOString()}`
    );
    return res.data as HourlyClick[];
  } catch (error) {
    console.error("Error fetching analytics data:", error);
    throw error;
  }
}

export async function getCountryClicksData(urlId: string) {
  try {
    const res = await axios.get(`/analytics/${urlId}/country`);
    return res.data as CountryClick[];
  } catch (error) {
    console.error("Error fetching analytics data:", error);
    throw error;
  }
}
export async function getCityClicksData(urlId: string) {
  try {
    const res = await axios.get(`/analytics/${urlId}/city`);
    return res.data as CityClick[];
  } catch (error) {
    console.error("Error fetching analytics data:", error);
    throw error;
  }
}

export async function getReferrerClicksData(urlId: string) {
  try {
    const res = await axios.get(`/analytics/${urlId}/referrer`);
    return res.data as ReferrerClick[];
  } catch (error) {
    console.error("Error fetching analytics data:", error);
    throw error;
  }
}

export async function getBrowserClicksData(urlId: string) {
  try {
    const res = await axios.get(`/analytics/${urlId}/browser`);
    return res.data as BrowserClick[];
  } catch (error) {
    console.error("Error fetching analytics data:", error);
    throw error;
  }
}

export async function getDeviceClicksData(urlId: string) {
  try {
    const res = await axios.get(`/analytics/${urlId}/device`);
    return res.data as DeviceClick[];
  } catch (error) {
    console.error("Error fetching analytics data:", error);
    throw error;
  }
}
