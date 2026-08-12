import axios from "axios";
import { useAuthStore } from "../store/useAuthStore";
import { queryClient } from "../store/queryClient";

const api = axios.create({
  baseURL: "",
  headers: {
    "Content-Type": "application/json",
  },
});

let isRefreshing = false;
let failedQueue = [];

const processQueue = (error, token = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

const handleForceLogout = () => {
  useAuthStore.getState().logout();
  queryClient.clear();
  window.location.href = "/login";
};

api.interceptors.request.use(
  (config) => {
    const accessToken = useAuthStore.getState().accessToken;

    if (accessToken) {
      config.headers["Authorization"] = `Bearer ${accessToken}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (!error.response) {
      return Promise.reject(error);
    }

    const status = error.response.status;

    if (status === 403) {
      const responseData = error.response.data;
      const message = responseData?.message || responseData?.error || "";

      if (message.includes("login in new dispositivy")) {
        console.warn("Sessão encerrada: login realizado em outro dispositivo.");
        handleForceLogout();
        return Promise.reject(error);
      }

      if (message.includes("user id not is premiun subscription")) {
        console.warn("Acesso negado: recurso exclusivo para usuários premium.");
        window.location.href = "/pricing";
        return Promise.reject(error);
      }
    }

    if (status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            originalRequest.headers["Authorization"] = `Bearer ${token}`;
            return api(originalRequest);
          })
          .catch((err) => Promise.reject(err));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      const { accessToken, refreshToken, setTokens } = useAuthStore.getState();

      try {
        const response = await axios.get("/api/refresh", {
          headers: {
            "access-token": `Bearer ${accessToken}`,
            "refresh-token": `Bearer ${refreshToken}`,
          },
        });

        const newAccessToken = response.data.access_token;
        const newRefreshToken = response.data.refresh_token;

        setTokens(newAccessToken, newRefreshToken);

        processQueue(null, newAccessToken);

        originalRequest.headers["Authorization"] = `Bearer ${newAccessToken}`;
        return api(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);
        handleForceLogout();
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  },
);

export default api;
