export interface UserType {
  ID: number;
  Username: string;
  Email: string;
  CreatedAt: string;
}
export interface UserStoreType {
  user: UserType | null;
  isLoggedIn: boolean;
}