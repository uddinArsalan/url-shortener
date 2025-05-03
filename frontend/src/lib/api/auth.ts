import axios from "axios";

export async function fetchUser() {
  try {
    const user = await axios.get(`http://localhost:4000/me`, {
      withCredentials: true,
    });
    return user.data;
  } catch (error) {
    console.log("Error Fetching User ", error);
    throw error;
  }
}
