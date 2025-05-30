import axios from "../axios";

export async function fetchUser() {
  try {
    const user = await axios.get(`/me`);
    return user.data;
  } catch (error) {
    // console.log("Error Fetching User ", error);
    throw new Error(`Error Fetching User`);
  }
}

export async function logout() {
  try {
    await axios.get(`/auth/logout`);
  } catch (error) {
    console.log("Error Logging out User ", error);
    if (error instanceof Error) {
      throw new Error(`Error Logging out User: ${error.message}`);
    } else {
      throw new Error(`Error Logging out User: ${String(error)}`);
    }
  }
}
