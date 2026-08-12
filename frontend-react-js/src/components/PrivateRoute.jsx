import { Navigate } from "react-router-dom";
import { useAuthStore } from "../store/useAuthStore";
import { useVerifyPremium } from "../store/useVerifyPremium";

function RoutePrivate({ children, requirePremiun = false }) {
  const accessToken = useAuthStore((state) => state.accessToken);

  const isPremiumEnabled = Boolean(accessToken && requirePremiun);
  const { isLoading, isError } = useVerifyPremium(isPremiumEnabled);

  if (!accessToken) {
    return <Navigate to="/login"></Navigate>;
  }

  if (requirePremiun && isLoading) {
    return (
      <div
        style={{ display: "flex", justifyContent: "center", padding: "2rem" }}
      >
        <p>Validando assinatura premium...</p>
      </div>
    );
  }

  if (requirePremiun && isError) {
    return null;
  }

  return children;
}

export default RoutePrivate;
