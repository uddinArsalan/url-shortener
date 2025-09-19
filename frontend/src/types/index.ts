export interface UserType {
  id: number;
  username: string;
  email: string;
  created_at: string;
}

export interface UrlType {
  id: string;
  shortcode: string;
  original_url: string;
  created_at: string;
}

interface PaginationType {
  next_cursor: string;
  has_more: boolean;
}

export interface UrlResponseType {
  urls: UrlType[];
  pagination: PaginationType;
}

export interface UserStoreType {
  user: UserType | null;
  isLoggedIn: boolean;
  isLoading : boolean
}

type DeviceType = "desktop" | "mobile" | "bot";

export interface ClickAnalyticsType {
  id: number;
  browser: string;
  os: string;
  country: string;
  city: string;
  device: DeviceType;
  referrer: string;
  timestamp: string;
}

export interface ClickAnalyticsResponseType {
  click_analytics: ClickAnalyticsType[];
  total_clicks: number;
  unique_clicks: number;
}

export type HourlyClick = {
  hour: string;
  count: number;
};
export type CountryClick = {
  country: string;
  count: number;
};
export type CityClick = {
  city: string;
  count: number;
};

export type DeviceClick = {
  device: DeviceType;
  count: number;
};

export type BrowserClick = {
  browser: string;
  count: number;
};

export type ReferrerClick = {
  referrer: string;
  count: number;
};
