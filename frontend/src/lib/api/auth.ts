import axios from "axios";
export async function login() {
  try {
    await axios.get(`http://localhost:4000/auth/login`);
  } catch (error) {
    console.log("Error Logging ", error);
    throw error;
  }
}
