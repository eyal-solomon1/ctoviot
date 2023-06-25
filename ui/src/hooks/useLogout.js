import axios from "../api/axios";
import useAuth from "./useAuth";
import Cookies from "js-cookie";

const useLogout = () => {
  const { setAuth } = useAuth();

  const logout = async () => {
    setAuth({});
    try {
      await axios.post(
        "/users/logout",
        {},
        {
          withCredentials: true,
        }
      );
      Cookies.remove("refresh_token");
    } catch (err) {
      console.error(err);
    }
  };

  return logout;
};

export default useLogout;
