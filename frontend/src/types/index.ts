export interface UserType {
  id: number;
  username: string;
  email: string;
  created_at: string;
}

export interface UrlType{
  id : string;
  shortcode : string
  original_url : string;
  created_at : string;
}

interface PaginationType{
  next_cursor: string;
  has_more : boolean;
}

export interface UrlResponseType{
  urls : UrlType[];
  pagination : PaginationType
}

export interface UserStoreType {
  user: UserType;
  isLoggedIn: boolean;
}

type DeviceType = 'desktop' | 'mobile' | 'bot';

export interface ClickAnalyticsType {
  id: number;
  browser : string;
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
  unique_clicks : number;
}
