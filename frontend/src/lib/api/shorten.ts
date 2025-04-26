import axios from "axios";

export const shortenUrl = async (originalUrl: string) => {
  try {
    const res = await axios.post(`http://localhost:4000/shorten?url=${originalUrl}`);
    return res.data as string;
  } catch (error) {
    console.log("Error Shortening Url ", error);
    throw error
  }
};
