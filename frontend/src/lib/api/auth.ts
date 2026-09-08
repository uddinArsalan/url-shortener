import { API_BASE_URL } from "../../constants";
import axios, { setLoggingOut } from "../axios";

export async function fetchUser() {
  try {
    const user = await axios.get(`/me`);
    return user.data;
  } catch (error) {
    // console.log("Error Fetching User ", error);
    throw new Error(`Error Fetching User`);
  }
}

export function logout() {
  setLoggingOut(true); 
  window.location.href = `${API_BASE_URL}/auth/logout`;
}