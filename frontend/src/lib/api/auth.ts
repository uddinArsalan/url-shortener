import axios from "axios";
export async function login() {
  try {
    await axios.get(`http://localhost:4000/auth/login`);
  } catch (error) {
    console.log("Error Logging ", error);
    throw error;
  }
}

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
