import { useSearchParams, useNavigate } from "react-router-dom";
import { MailCheck, ArrowRight } from "lucide-react";

export default function ResetPassConfirm() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  // Pega o e-mail da URL (ex: /resetpass-confirm?email=usuario@email.com)
  const email = searchParams.get("email");

  return (
    <div className="flex-1 bg-slate-950 text-slate-100 flex items-center justify-center p-4 md:p-6 min-h-screen">
      <div className="max-w-md w-full bg-slate-900 border border-slate-800 rounded-2xl p-6 md:p-8 text-center space-y-6 shadow-2xl">
        {/* ÍCONE */}
        <div className="w-16 h-16 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400 mx-auto shadow-inner">
          <MailCheck className="w-8 h-8" />
        </div>

        {/* MENSAGEM */}
        <div className="space-y-2">
          <h2 className="text-xl font-bold text-white">E-mail Enviado!</h2>
          <p className="text-xs text-slate-400 leading-relaxed">
            Enviamos um link de redefinição de senha para{" "}
            <span className="font-semibold text-slate-200">
              {email || "o seu e-mail"}
            </span>
            . Verifique sua caixa de entrada e a pasta de spam.
          </p>
        </div>

        {/* BOTÃO */}
        <button
          onClick={() => navigate("/login")}
          className="w-full py-3 px-4 bg-indigo-600 hover:bg-indigo-500 text-white font-bold rounded-xl text-sm transition-all shadow-lg shadow-indigo-600/25 flex items-center justify-center gap-2 cursor-pointer"
        >
          <span>Ir para o Login</span>
          <ArrowRight className="w-4 h-4" />
        </button>
      </div>
    </div>
  );
}
