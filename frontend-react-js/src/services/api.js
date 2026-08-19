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
  failedQueue.forEach((promise) => {
    if (error) {
      promise.reject(error);
    } else {
      promise.resolve(token);
    }
  });
  failedQueue = [];
};

// Adicione a palavra 'export' na declaração da função
export const handleForceLogout = async () => {
  console.log("--> Entrei na handleForceLogout!"); // Log de teste

  const { accessToken, refreshToken, logout } = useAuthStore.getState();

  if (accessToken || refreshToken) {
    const headers = {};
    if (accessToken) headers["access-token"] = `Bearer ${accessToken}`;
    if (refreshToken) headers["refresh-token"] = `Bearer ${refreshToken}`;

    try {
      const resp = await axios.get("/api/logout", { headers });
      console.log("Sucesso backend:", resp.data);
    } catch (error) {
      console.warn("Erro no logout backend:", error);
    }
  }

  logout();
  queryClient.clear();
  window.location.href = "/login";
};
// 1. INTERCEPTOR DE REQUISIÇÃO
api.interceptors.request.use(
  (config) => {
    // Só injeta do Zustand se o cabeçalho AINDA NÃO FOI preenchido manualmente pelo retry
    const existingHeader =
      config.headers["access-token"] || config.headers.get?.("access-token");

    if (!existingHeader) {
      const accessToken = useAuthStore.getState().accessToken;
      if (accessToken) {
        if (config.headers.set) {
          config.headers.set("access-token", `Bearer ${accessToken}`);
        } else {
          config.headers["access-token"] = `Bearer ${accessToken}`;
        }
      }
    }
    return config;
  },
  (error) => Promise.reject(error),
);

// 2. INTERCEPTOR DE RESPOSTA
api.interceptors.response.use(
  (response) => response,
  (error) => {
    const originalRequest = error.config;

    if (!error.response) {
      return Promise.reject(error);
    }

    const status = error.response.status;

    if (status === 403) {
      const message =
        error.response.data?.message || error.response.data?.error || "";

      if (message.includes("login in new dispositivy")) {
        handleForceLogout();
        return Promise.reject(error);
      }

      return Promise.reject(error);
    }

    // if (message.includes("user id not is premiun subscription")) {
    //   console.log("redirect user not premiun");
    //   window.location.href = "/";

    //   return Promise.reject(error);
    // }

    if (error.response.status === 401 && !originalRequest._retry) {
      // Evita loop se o 401 vier do próprio endpoint de refresh
      if (originalRequest.url?.includes("/api/refresh")) {
        handleForceLogout();
        return Promise.reject(error);
      }

      originalRequest._retry = true;

      const { accessToken, refreshToken, setTokens } = useAuthStore.getState();

      if (!refreshToken) {
        handleForceLogout();
        return Promise.reject(error);
      }

      if (!isRefreshing) {
        isRefreshing = true;

        axios
          .get("/api/refresh", {
            headers: {
              "access-token": `Bearer ${accessToken}`,
              "refresh-token": `Bearer ${refreshToken}`,
            },
          })
          .then((response) => {
            // Acessa a propriedade 'tokens' retornada pelo Go
            const newAccessToken = response.data?.tokens?.access_token;
            const newRefreshToken = response.data?.tokens?.refresh_token;

            if (!newAccessToken) {
              throw new Error(
                "Novo access_token não encontrado na chave 'tokens'.",
              );
            }

            // 1. Atualiza Zustand
            setTokens(newAccessToken, newRefreshToken);

            // 2. Libera a fila com o novo token
            processQueue(null, newAccessToken);
          })
          .catch((refreshError) => {
            processQueue(refreshError, null);
            handleForceLogout();
          })
          .finally(() => {
            isRefreshing = false;
          });
      }

      return new Promise((resolve, reject) => {
        failedQueue.push({
          resolve: (newToken) => {
            // Injeta o novo token diretamente nos headers da requisição original
            const authHeaderValue = `Bearer ${newToken}`;

            if (originalRequest.headers.set) {
              originalRequest.headers.set("access-token", authHeaderValue);
            } else {
              originalRequest.headers["access-token"] = authHeaderValue;
            }

            // Executa a requisição novamente
            resolve(api(originalRequest));
          },
          reject: (err) => {
            reject(err);
          },
        });
      });
    }

    return Promise.reject(error);
  },
);

export default api;
