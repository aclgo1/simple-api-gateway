import { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import {
  Loader2,
  KeyRound,
  RefreshCw,
  ArrowRight,
  AlertCircle,
  CheckCircle2,
} from "lucide-react";
import api from "../services/api";

const resetPassSchema = z.object({
  email: z
    .string()
    .min(1, "E-mail é obrigatório")
    .email("Formato de e-mail inválido"),
  captchaAnswer: z.string().min(1, "Digite o código da imagem"),
});

export default function ResetPass() {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState(null);
  const [isSubmitted, setIsSubmitted] = useState(false);

  // Hook Form com validação Zod
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    resolver: zodResolver(resetPassSchema),
  });

  // Query do Captcha (com chave única e sem cache)
  const {
    data: captcha,
    isLoading: loadingCaptcha,
    refetch: refetchCaptcha,
  } = useQuery({
    queryKey: ["captcha", "resetpass"], // Chave única para não conflitar com o Login
    queryFn: async () => {
      const resp = await api.get("/api/captcha");
      return resp.data;
    },
    staleTime: 0, // Garante que o dado sempre seja considerado obsoleto
    gcTime: 0, // Remove do cache imediatamente ao desmontar
    refetchOnMount: "always", // Força nova busca sempre que a tela abrir
    refetchOnWindowFocus: false,
  });

  // Mutation do Reset Pass
  const resetPassMutation = useMutation({
    mutationFn: async (formData) => {
      setServerError(null);
      const resp = await api.get(
        `/api/user/resetpass/${formData.email}/${captcha?.id}/${formData.captchaAnswer}`,
      );
      return { responseData: resp.data, userEmail: formData.email };
    },
    onSuccess: (data) => {
      setIsSubmitted(true);
      setTimeout(() => {
        // Redireciona passando o e-mail via query param
        navigate(
          `/resetpass-confirm?email=${encodeURIComponent(data.userEmail)}`,
        );
      }, 1500);
    },
    onError: (err) => {
      const msg =
        err.response?.data?.message ||
        "Erro ao solicitar redefinição de senha.";
      setServerError(msg);
      refetchCaptcha(); // Atualiza o captcha em caso de erro
    },
  });

  const onSubmit = (data) => {
    resetPassMutation.mutate(data);
  };

  return (
    <div className="flex-1 bg-slate-950 text-slate-100 flex items-center justify-center p-4 md:p-6 min-h-screen">
      <div className="max-w-md w-full bg-slate-900 border border-slate-800 rounded-2xl p-6 md:p-8 space-y-6 shadow-2xl">
        {/* CABEÇALHO */}
        <div className="text-center space-y-2">
          <div className="w-16 h-16 rounded-2xl bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center text-indigo-400 mx-auto shadow-inner">
            <KeyRound className="w-8 h-8" />
          </div>
          <h2 className="text-xl font-bold text-white">Recuperar Senha</h2>
          <p className="text-xs text-slate-400">
            Informe seu e-mail e o código da imagem para receber o link de
            redefinição.
          </p>
        </div>

        {/* FEEDBACK DE SUCESSO */}
        {isSubmitted ? (
          <div className="space-y-4 py-4 text-center">
            <div className="w-12 h-12 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400 mx-auto">
              <CheckCircle2 className="w-6 h-6" />
            </div>
            <div className="space-y-1">
              <h3 className="text-base font-bold text-white">
                E-mail Enviado!
              </h3>
              <p className="text-xs text-slate-400">
                Redirecionando para a página de confirmação...
              </p>
            </div>
          </div>
        ) : (
          /* FORMULÁRIO */
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            {/* CAMPO: E-MAIL */}
            <div className="space-y-1.5">
              <label className="text-xs font-medium text-slate-300">
                E-mail ou Usuário
              </label>
              <input
                type="email"
                placeholder="seu@email.com"
                {...register("email")}
                className="w-full px-3.5 py-2.5 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all"
              />
              {errors.email && (
                <p className="text-xs text-red-400 flex items-center gap-1">
                  <AlertCircle className="w-3.5 h-3.5" />
                  {errors.email.message}
                </p>
              )}
            </div>

            {/* CAMPO: CAPTCHA */}
            <div className="space-y-2">
              <label className="text-xs font-medium text-slate-300">
                Código de Verificação
              </label>

              {/* Contêiner Flex: Captcha esticado + Botão do lado de fora */}
              <div className="flex items-center gap-2">
                {/* Área da Imagem do Captcha (Comprida e totalmente preenchida) */}
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

                {/* Botão de recarregar (Lado de fora) */}
                <button
                  type="button"
                  onClick={() => refetchCaptcha()}
                  className="h-16 w-14 shrink-0 bg-slate-900 hover:bg-slate-800 text-slate-300 hover:text-white rounded-xl transition-all border border-slate-800 flex items-center justify-center cursor-pointer shadow-md"
                  title="Atualizar código"
                >
                  <RefreshCw
                    className={`w-5 h-5 ${loadingCaptcha ? "animate-spin text-indigo-400" : ""}`}
                  />
                </button>
              </div>

              {/* Campo de Entrada do Texto */}
              <input
                type="text"
                placeholder="Digite o código acima"
                {...register("captchaAnswer")}
                className="w-full px-3.5 py-2.5 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all font-mono tracking-widest text-center uppercase"
              />

              {errors.captchaAnswer && (
                <p className="text-xs text-red-400 flex items-center gap-1 mt-1">
                  <AlertCircle className="w-3.5 h-3.5" />
                  {errors.captchaAnswer.message}
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
              disabled={resetPassMutation.isPending}
              className="w-full py-3 px-4 bg-indigo-600 hover:bg-indigo-500 disabled:bg-indigo-600/50 text-white font-bold rounded-xl text-sm transition-all shadow-lg shadow-indigo-600/25 flex items-center justify-center gap-2 cursor-pointer disabled:cursor-not-allowed"
            >
              {resetPassMutation.isPending ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin" />
                  <span>Enviando...</span>
                </>
              ) : (
                <>
                  <span>Enviar Link de Recuperação</span>
                  <ArrowRight className="w-4 h-4" />
                </>
              )}
            </button>

            {/* VOLTAR AO LOGIN */}
            <button
              type="button"
              onClick={() => navigate("/login")}
              className="w-full py-2.5 text-xs text-slate-400 hover:text-white transition-all text-center cursor-pointer"
            >
              Voltar para o Login
            </button>
          </form>
        )}
      </div>
    </div>
  );
}
