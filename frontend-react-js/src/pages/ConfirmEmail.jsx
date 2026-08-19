import { useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { Loader2, CheckCircle2, XCircle, ArrowRight } from "lucide-react";
import { useAuthStore } from "../store/useAuthStore";
import api from "../services/api";

export default function ConfirmEmail() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const setTokens = useAuthStore((state) => state.setTokens);

  const code = searchParams.get("code");

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function confirmCode() {
      if (!code) {
        setError("Código de verificação ausente ou inválido.");
        setLoading(false);
        return;
      }

      try {
        const response = await api.get(`/api/user/confirm/${code}`);
        const data = response.data;

        if (response.status < 200 || response.status >= 300) {
          throw new Error(
            data.message || "Falha ao confirmar o código de cadastro.",
          );
        }

        if (data.tokens.access_token) {
          setTokens(data.tokens.access_token, data.tokens.refresh_tokens);
        }

        setTimeout(() => {
          const pendingPlan = sessionStorage.getItem("pendingCheckoutPlan");
          const pendingMethod = sessionStorage.getItem("pendingCheckoutMethod");

          if (pendingPlan && pendingMethod) {
            sessionStorage.removeItem("pendingCheckoutPlan");
            sessionStorage.removeItem("pendingCheckoutMethod");
            navigate(
              `/checkout?action=new-subscription&plan=${pendingPlan}&method=${pendingMethod}`,
            );
          } else {
            navigate("/home");
          }
        }, 1500);
      } catch (err) {
        setError(err.message || "Erro ao conectar com o servidor.");
      } finally {
        setLoading(false);
      }
    }

    confirmCode();
  }, [code, navigate, setTokens]);

  return (
    <div className="flex-1 bg-slate-950 text-slate-100 flex items-center justify-center p-4 md:p-6">
      <div className="max-w-md w-full bg-slate-900 border border-slate-800 rounded-2xl p-6 md:p-8 text-center space-y-6 shadow-2xl">
        {/* ESTADO 1: CARREGANDO */}
        {loading && (
          <div className="space-y-4 py-8">
            <div className="w-16 h-16 rounded-2xl bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center text-indigo-400 mx-auto shadow-inner">
              <Loader2 className="w-8 h-8 animate-spin" />
            </div>
            <div className="space-y-1">
              <h2 className="text-xl font-bold text-white">
                Validando Código...
              </h2>
              <p className="text-xs text-slate-400">
                Estamos ativando sua conta no gateway.
              </p>
            </div>
          </div>
        )}

        {/* ESTADO 2: ERRO */}
        {!loading && error && (
          <div className="space-y-6">
            <div className="w-16 h-16 rounded-2xl bg-red-500/10 border border-red-500/20 flex items-center justify-center text-red-400 mx-auto shadow-inner">
              <XCircle className="w-8 h-8" />
            </div>
            <div className="space-y-2">
              <h2 className="text-xl font-bold text-white">
                Falha na Ativação
              </h2>
              <p className="text-xs text-slate-400 bg-red-500/10 border border-red-500/20 p-3 rounded-xl">
                {error}
              </p>
            </div>
            <button
              onClick={() => navigate("/login")}
              className="w-full py-3 px-4 bg-slate-800 hover:bg-slate-700 text-white font-bold text-xs rounded-xl transition-all border border-slate-700 cursor-pointer"
            >
              Voltar para o Login
            </button>
          </div>
        )}

        {/* ESTADO 3: SUCESSO */}
        {!loading && !error && (
          <div className="space-y-6">
            <div className="w-16 h-16 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400 mx-auto shadow-inner">
              <CheckCircle2 className="w-8 h-8" />
            </div>
            <div className="space-y-2">
              <h2 className="text-xl font-bold text-white">
                E-mail Confirmado!
              </h2>
              <p className="text-xs text-slate-400">
                Sua conta foi ativada com sucesso. Você será redirecionado em
                instantes...
              </p>
            </div>
            <button
              onClick={() => navigate("/home")}
              className="w-full py-3 px-4 bg-indigo-600 hover:bg-indigo-500 text-white font-bold rounded-xl text-sm transition-all shadow-lg shadow-indigo-600/25 flex items-center justify-center gap-2 cursor-pointer"
            >
              <span>Acessar Painel Agora</span>
              <ArrowRight className="w-4 h-4" />
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
