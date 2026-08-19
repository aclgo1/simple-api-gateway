import { useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useQuery, useMutation } from "@tanstack/react-query";
import {
  Loader2,
  KeyRound,
  RefreshCw,
  ArrowRight,
  AlertCircle,
  CheckCircle2,
  XCircle,
} from "lucide-react";
import { useAuthStore } from "../store/useAuthStore";
import api from "../services/api";

// Schema de validação com Zod
const newPassSchema = z
  .object({
    new_pass: z.string().min(6, "A senha deve ter no mínimo 6 caracteres"),
    confirm_pass: z.string().min(1, "Confirme a nova senha"),
    captcha_awnser: z.string().min(1, "Digite o código do captcha"),
  })
  .refine((data) => data.new_pass === data.confirm_pass, {
    message: "As senhas não conferem",
    path: ["confirm_pass"],
  });

export default function NewPass() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const setTokens = useAuthStore((state) => state.setTokens);

  const code = searchParams.get("code");

  const [serverError, setServerError] = useState(null);
  const [isSuccess, setIsSuccess] = useState(false);

  // Form Hook
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    resolver: zodResolver(newPassSchema),
  });

  // Query do Captcha (Isolada e sem cache)
  const {
    data: captcha,
    isLoading: loadingCaptcha,
    refetch: refetchCaptcha,
  } = useQuery({
    queryKey: ["captcha", "newpass"],
    queryFn: async () => {
      const resp = await api.get("/api/captcha");
      return resp.data;
    },
    staleTime: 0,
    gcTime: 0,
    refetchOnMount: "always",
    refetchOnWindowFocus: false,
  });

  // Mutation do POST de Nova Senha
  const newPassMutation = useMutation({
    mutationFn: async (formData) => {
      setServerError(null);

      // Payload correspondente à struct ParamsNewPass em Go
      const payload = {
        NewPassCode: code,
        new_pass: formData.new_pass,
        confirm_pass: formData.confirm_pass,
        captcha_id: captcha?.id,
        captcha_awnser: formData.captcha_awnser,
      };

      const resp = await api.post("/api/user/newpass", payload);
      return resp.data;
    },
    onSuccess: (data) => {
      setIsSuccess(true);

      // Grava tokens no Zustand se retornados
      if (data?.tokens?.access_token) {
        setTokens(data.tokens.access_token, data.tokens.refresh_token);
      }

      // Redirecionamento após 1.5s
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
    },
    onError: (err) => {
      const msg =
        err.response?.data?.message ||
        err.response?.data?.error ||
        "Erro ao redefinir senha. Verifique os dados e o captcha.";
      setServerError(msg);
      refetchCaptcha(); // Atualiza o captcha se houver erro
    },
  });

  const onSubmit = (formData) => {
    if (!code) {
      setServerError("Código de recuperação ausente ou inválido.");
      return;
    }
    newPassMutation.mutate(formData);
  };

  // Se o link veio sem o parâmetro ?code=
  if (!code) {
    return (
      <div className="flex-1 bg-slate-950 text-slate-100 flex items-center justify-center p-4 md:p-6 min-h-screen">
        <div className="max-w-md w-full bg-slate-900 border border-slate-800 rounded-2xl p-6 text-center space-y-6 shadow-2xl">
          <div className="w-16 h-16 rounded-2xl bg-red-500/10 border border-red-500/20 flex items-center justify-center text-red-400 mx-auto">
            <XCircle className="w-8 h-8" />
          </div>
          <div className="space-y-2">
            <h2 className="text-xl font-bold text-white">Código Inválido</h2>
            <p className="text-xs text-slate-400">
              O link de redefinição de senha está incompleto ou expirou.
            </p>
          </div>
          <button
            onClick={() => navigate("/login")}
            className="w-full py-3 px-4 bg-slate-800 hover:bg-slate-700 text-white font-bold text-xs rounded-xl border border-slate-700 cursor-pointer"
          >
            Voltar para o Login
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 bg-slate-950 text-slate-100 flex items-center justify-center p-4 md:p-6 min-h-screen">
      <div className="max-w-md w-full bg-slate-900 border border-slate-800 rounded-2xl p-6 md:p-8 space-y-6 shadow-2xl">
        {/* CABEÇALHO */}
        <div className="text-center space-y-2">
          <div className="w-16 h-16 rounded-2xl bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center text-indigo-400 mx-auto shadow-inner">
            <KeyRound className="w-8 h-8" />
          </div>
          <h2 className="text-xl font-bold text-white">Criar Nova Senha</h2>
          <p className="text-xs text-slate-400">
            Informe e confirme sua nova senha abaixo.
          </p>
        </div>

        {/* FEEDBACK DE SUCESSO */}
        {isSuccess ? (
          <div className="space-y-4 py-4 text-center">
            <div className="w-16 h-16 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400 mx-auto">
              <CheckCircle2 className="w-8 h-8" />
            </div>
            <div className="space-y-1">
              <h3 className="text-base font-bold text-white">
                Senha Alterada com Sucesso!
              </h3>
              <p className="text-xs text-slate-400">
                Redirecionando para a plataforma...
              </p>
            </div>
          </div>
        ) : (
          /* FORMULÁRIO */
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            {/* NOVA SENHA */}
            <div className="space-y-1.5">
              <label className="text-xs font-medium text-slate-300">
                Nova Senha
              </label>
              <input
                type="password"
                placeholder="••••••••"
                {...register("new_pass")}
                className="w-full px-3.5 py-2.5 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all"
              />
              {errors.new_pass && (
                <p className="text-xs text-red-400 flex items-center gap-1">
                  <AlertCircle className="w-3.5 h-3.5" />
                  {errors.new_pass.message}
                </p>
              )}
            </div>

            {/* CONFIRMAR NOVA SENHA */}
            <div className="space-y-1.5">
              <label className="text-xs font-medium text-slate-300">
                Confirmar Nova Senha
              </label>
              <input
                type="password"
                placeholder="••••••••"
                {...register("confirm_pass")}
                className="w-full px-3.5 py-2.5 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all"
              />
              {errors.confirm_pass && (
                <p className="text-xs text-red-400 flex items-center gap-1">
                  <AlertCircle className="w-3.5 h-3.5" />
                  {errors.confirm_pass.message}
                </p>
              )}
            </div>

            {/* CAPTCHA */}
            <div className="space-y-2">
              <label className="text-xs font-medium text-slate-300">
                Código de Verificação
              </label>

              <div className="flex items-center gap-2">
                <div className="flex-1 h-16 bg-slate-950 border border-slate-800 rounded-xl overflow-hidden flex items-center justify-center shadow-inner">
                  {loadingCaptcha ? (
                    <Loader2 className="w-6 h-6 animate-spin text-indigo-400" />
                  ) : captcha?.base64_image ? (
                    <img
                      src={captcha.base64_image}
                      alt="Captcha"
                      className="w-full h-full object-fill contrast-125 brightness-110"
                    />
                  ) : (
                    <span className="text-xs text-slate-500 font-medium">
                      Falha ao carregar imagem
                    </span>
                  )}
                </div>

                <button
                  type="button"
                  onClick={() => refetchCaptcha()}
                  className="h-16 w-14 shrink-0 bg-slate-900 hover:bg-slate-800 text-slate-300 hover:text-white rounded-xl transition-all border border-slate-800 flex items-center justify-center cursor-pointer shadow-md"
                  title="Atualizar código"
                >
                  <RefreshCw
                    className={`w-5 h-5 ${
                      loadingCaptcha ? "animate-spin text-indigo-400" : ""
                    }`}
                  />
                </button>
              </div>

              <input
                type="text"
                placeholder="Digite o código acima"
                {...register("captcha_awnser")}
                className="w-full px-3.5 py-2.5 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all font-mono tracking-widest text-center uppercase"
              />

              {errors.captcha_awnser && (
                <p className="text-xs text-red-400 flex items-center gap-1">
                  <AlertCircle className="w-3.5 h-3.5" />
                  {errors.captcha_awnser.message}
                </p>
              )}
            </div>

            {/* MENSAGEM DE ERRO DO SERVIDOR */}
            {serverError && (
              <p className="text-xs text-red-400 bg-red-500/10 border border-red-500/20 p-3 rounded-xl flex items-center gap-2">
                <AlertCircle className="w-4 h-4 shrink-0" />
                <span>{serverError}</span>
              </p>
            )}

            {/* BOTÃO DE SUBMIT */}
            <button
              type="submit"
              disabled={newPassMutation.isPending}
              className="w-full py-3 px-4 bg-indigo-600 hover:bg-indigo-500 disabled:bg-indigo-600/50 text-white font-bold rounded-xl text-sm transition-all shadow-lg shadow-indigo-600/25 flex items-center justify-center gap-2 cursor-pointer disabled:cursor-not-allowed"
            >
              {newPassMutation.isPending ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin" />
                  <span>Redefinindo...</span>
                </>
              ) : (
                <>
                  <span>Redefinir Senha</span>
                  <ArrowRight className="w-4 h-4" />
                </>
              )}
            </button>
          </form>
        )}
      </div>
    </div>
  );
}
