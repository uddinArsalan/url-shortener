import axios from "$lib/axios";
import type { UrlResponseType } from "../../types";

export async function fetchUserUrls(cursor: string) {
  try {
    const response = await axios.get(`/url?cursor=${cursor}`);
    return response.data as UrlResponseType;
  } catch (error) {
    console.log("Error Fetching User URLs ", error);
    if (error instanceof Error) {
      throw new Error(`Error Fetching User URLs: ${error.message}`);
    } else {
      throw new Error(`Error Fetching User URLs: ${String(error)}`);
    }
  } 
}