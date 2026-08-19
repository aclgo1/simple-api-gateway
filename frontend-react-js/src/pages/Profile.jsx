import { useAuthStore } from "../store/useAuthStore";
import { useEffect, useState } from "react";
import api from "../services/api";

export default function Profile() {
  const userId = useAuthStore((state) => state.userId);
  const [userInfo, setUserInfo] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchUser() {
      try {
        const resp = await api.get("/api/user/find");
        setUserInfo(resp.data.user);
      } catch (error) {
        console.error("error find user: ", error);
      } finally {
        setLoading(false);
      }
    }

    fetchUser();
  }, []);

  if (loading) {
    return <div>Carregando...</div>;
  }

  return (
    <div>
      <h1>{userId}</h1>
      <h2>{`${userInfo?.name} - ${userInfo?.role} (${userInfo?.email})`}</h2>
    </div>
  );
}
