import axios from "$lib/axios";
import type { UrlResponseType } from "../../types";

export async function fetchUserUrls(cursor: string): Promise<UrlResponseType> {
  try {
    const response = await axios.get(`/url?cursor=${cursor}`);
    return response.data as UrlResponseType;
  } catch (error) {
    console.log("Error Fetching User URLs ", error);
    throw error;
  }
}
