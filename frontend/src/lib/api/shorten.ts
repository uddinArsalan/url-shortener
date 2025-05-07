import axios from "../axios";

export const shortenUrl = async (originalUrl: string) => {
  try {
    const res = await axios.post(
      `/shorten?url=${encodeURIComponent(originalUrl)}`
    );
    return res.data as string;
  } catch (error) {
    console.log("Error Shortening Url ", error);
    if (error instanceof Error) {
      throw new Error(`Failed to shorten URL: ${error.message}`);
    } else {
      throw new Error("Failed to shorten URL: An unknown error occurred.");
    }
  }
};
