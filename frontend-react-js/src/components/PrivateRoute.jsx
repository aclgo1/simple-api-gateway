import { Navigate } from "react-router-dom";
import { useAuthStore } from "../store/useAuthStore";
import { useVerifyPremium } from "../store/useVerifyPremium";
import { useVerifyToken } from "../store/useVerifyToken";

function RoutePrivate({ children, requirePremium = false }) {
  const accessToken = useAuthStore((state) => state.accessToken);
  const refreshToken = useAuthStore((state) => state.refreshToken);
  const isInitializing = useAuthStore((state) => state.isInitializing);

  const hasTokens = Boolean(accessToken || refreshToken);

  // 1. Validação global do Token (Roda em todas as rotas privadas)
  const { isLoading: isCheckingToken, isError: isTokenInvalid } =
    useVerifyToken(hasTokens);

  // 2. Validação Premium (Só roda se o token for válido e se a rota exigir premium)
  const isPremiumEnabled = hasTokens && !isTokenInvalid && requirePremium;
  const { isLoading: isCheckingPremium, isError: isNotPremium } =
    useVerifyPremium(isPremiumEnabled);

  // A. Aguarda inicialização do Zustand (localStorage)
  if (isInitializing) {
    return (
      <div className="flex items-center justify-center h-screen bg-gray-900 text-white">
        <p>Carregando sessão...</p>
      </div>
    );
  }

  // B. Se não há tokens salvos, manda para o login
  if (!hasTokens) {
    return <Navigate to="/login" replace />;
  }

  // C. Aguarda as requisições de validação (com suporte a refresh automático do Axios no 401)
  if (isCheckingToken || (requirePremium && isCheckingPremium)) {
    return (
      <div className="flex items-center justify-center h-screen bg-gray-900 text-white">
        <p>Verificando permissões de acesso...</p>
      </div>
    );
  }

  // D. Se a checagem de token falhou (e nem o refresh salvou) -> Login
  if (isTokenInvalid) {
    return <Navigate to="/login" replace />;
  }

  // E. Se a rota exige Premium e o usuário não possui -> Home/Dashboard
  if (requirePremium && isNotPremium) {
    return <Navigate to="/" replace />;
  }

  return children;
}

export default RoutePrivate;
