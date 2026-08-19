import { useQuery } from "@tanstack/react-query";
import api from "../services/api";

export const useVerifyPremium = (enabled) => {
  return useQuery({
    queryKey: ["check-premium"],
    queryFn: async () => {
      const resp = await api.get(`/api/user/is-premium`);
      return resp.data;
    },
    enabled,
    retry: false,
    staleTime: 1000 * 60 * 5,
  });
};
