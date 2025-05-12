export interface UserType {
  id: number;
  username: string;
  email: string;
  created_at: string;
}

export interface UrlType{
  id : number;
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