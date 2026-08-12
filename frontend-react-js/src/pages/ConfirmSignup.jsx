import { Mail, ArrowRight, CheckCircle2 } from "lucide-react";
import { useSearchParams, useNavigate } from "react-router-dom";

export default function ConfirmSignup() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const email = searchParams.get("email") || "seu e-mail";

  return (
    <div className="flex-1 w-full flex items-center justify-center p-4 bg-slate-950 text-slate-100">
      <div className="max-w-md w-full bg-slate-900 border border-slate-800 rounded-2xl p-6 md:p-8 text-center space-y-6 shadow-2xl">
        {/* Ícone de Sucesso */}
        <div className="mx-auto w-16 h-16 rounded-full bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center text-indigo-400">
          <Mail className="w-8 h-8" />
        </div>

        {/* Conteúdo Principal */}
        <div className="space-y-2">
          <h1 className="text-2xl font-bold text-white tracking-tight">
            Verifique seu e-mail
          </h1>
          <p className="text-sm text-slate-400 leading-relaxed">
            Enviamos um link de confirmação para:
          </p>
          <p className="text-sm font-semibold text-indigo-400 bg-indigo-500/10 py-1.5 px-3 rounded-lg border border-indigo-500/20 break-all inline-block">
            {email}
          </p>
        </div>

        {/* Instruções */}
        <div className="bg-slate-950/60 p-4 rounded-xl border border-slate-800 text-left space-y-2 text-xs text-slate-300">
          <div className="flex items-start gap-2">
            <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0 mt-0.5" />
            <span>Abra sua caixa de entrada e clique no link de ativação.</span>
          </div>
          <div className="flex items-start gap-2">
            <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0 mt-0.5" />
            <span>
              Se não encontrar, verifique a pasta de <strong>Spam</strong> ou
              Lixo Eletrônico.
            </span>
          </div>
        </div>

        {/* Ação */}
        <div className="pt-2 border-t border-slate-800">
          <button
            onClick={() => navigate("/login")}
            className="w-full py-2.5 px-4 bg-indigo-600 hover:bg-indigo-500 text-white font-semibold text-sm rounded-xl transition-all shadow-lg shadow-indigo-600/20 flex items-center justify-center gap-2"
          >
            <span>Ir para o Login</span>
            <ArrowRight className="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  );
}
