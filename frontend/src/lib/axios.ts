import axios from "axios";
import { API_BASE_URL } from "../constants";
import { userStore } from "./store/userStore";

const axiosInstance = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
});

axiosInstance.interceptors.response.use(
  function (response) {
    return response;
  },
  async function (error) {
    // console.log("Error in Axios Interceptor: ", error);
    if (error.response && error.response.status === 401) {
      userStore.set({
        user: {
          id: 0,
          username: "",
          email: "",
          created_at: "",
        },
        isLoggedIn: false,
      });
      if (typeof window !== "undefined") {
        window.location.href = `${API_BASE_URL}/auth/login`;
      }
    }
    return Promise.reject(error);
  }
);

export default axiosInstance;
