import type { LayoutLoad } from "./$types";
import { fetchUser } from "$lib/api/auth";
import { userStore } from "$lib/store/userStore";
// import { redirect } from "@sveltejs/kit";
// import { API_BASE_URL } from "../constants";

export const load: LayoutLoad = async () => {
  try {
    const userDetails = await fetchUser();
    userStore.update((store) => {
      store.user = userDetails;
      store.isLoggedIn = true;
      // store.isLoading = false;
      return store;
    });
  } catch (error) {
    // console.log("Failed to fetch user details:", error);
    userStore.update((store) => {
      store.user = {
        id: 0,
        email: "",
        username: "",
        created_at: "",
      };
      store.isLoggedIn = false;
      // store.isLoading = false;
      return store;
    });
    // redirect(307, `${API_BASE_URL}/auth/login`);
  }
  return {};
};
