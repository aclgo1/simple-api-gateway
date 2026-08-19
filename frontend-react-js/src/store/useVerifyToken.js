// hooks/useVerifyToken.js
import { useQuery } from "@tanstack/react-query";
import api from "../services/api";

export const useVerifyToken = (enabled) => {
  return useQuery({
    queryKey: ["check-token-validity"],
    queryFn: async () => {
      const resp = await api.get(`/api/auth/is-valid`);
      return resp.data;
    },
    enabled,
    retry: false,
    staleTime: 1000 * 60 * 5, // 5 minutos
  });
};
