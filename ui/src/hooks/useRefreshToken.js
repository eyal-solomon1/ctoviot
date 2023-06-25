import axios from "../api/axios";
import useAuth from "./useAuth";

const useRefreshToken = () => {
  const { setAuth } = useAuth();

  const refresh = async () => {
    const response = await axios.post(
      "/users/refresh_token",
      {},
      {
        withCredentials: true,
      }
    );
    setAuth((prev) => {
      return {
        ...prev,
        accessToken: response.data.payload.access_token,
      };
    });
    return response.data.payload.access_token;
  };
  return refresh;
};

export default useRefreshToken;
