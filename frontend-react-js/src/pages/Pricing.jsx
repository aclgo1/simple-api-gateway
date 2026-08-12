import { useState } from "react";
import { Check, CreditCard, QrCode, Zap, Sparkles } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { useAuthStore } from "../store/useAuthStore";

export default function Pricing() {
  const navigate = useNavigate();
  const [selectedMethod, setSelectedMethod] = useState("pix");

  const accessToken = useAuthStore((state) => state.accessToken);

  const plans = [
    {
      id: "7_days",
      name: "Semanal",
      price: "19,99",
      period: "/semana",
      description: "Ideal para testar ou demandas pontuais.",
      popular: false,
      features: [
        "Acesso total por 7 dias",
        "Suporte por e-mail",
        "Sem fidelidade",
      ],
    },
    {
      id: "1_month",
      name: "Mensal",
      price: "34,99",
      period: "/mês",
      description: "A escolha favorita para resultados consistentes.",
      popular: true,
      features: [
        "Acesso ilimitado",
        "Suporte prioritário",
        "Projetos ilimitados",
        "Cancele quando quiser",
      ],
    },
    {
      id: "1_year",
      name: "Anual",
      price: "289,00",
      period: "/ano",
      description: "Economize mais de 30% no acesso anual.",
      popular: false,
      features: [
        "Tudo do plano Mensal",
        "Desconto exclusivo de 31%",
        "Atendimento VIP dedicado",
      ],
    },
  ];

  const handleCheckout = (planId) => {
    if (!accessToken) {
      sessionStorage.setItem("pendingCheckoutPlan", planId);
      sessionStorage.setItem("pendingCheckoutMethod", selectedMethod);
      navigate("/login");
    }
    navigate(`/checkout?plan=${planId}&method=${selectedMethod}`);
  };

  return (
    // Removido min-h-[calc...] e overflow-hidden. Ajustado padding vertical para py-4
    <div className="flex-1 bg-slate-950 text-slate-100 flex flex-col justify-center px-4 py-4 md:px-6">
      {/* Cabeçalho Compacto */}
      <header className="text-center space-y-2 mb-4">
        <div className="space-y-1">
          <div className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 text-[11px] font-semibold uppercase tracking-wider">
            <Sparkles className="w-3 h-3" /> Assinaturas Premium
          </div>
          <h1 className="text-xl md:text-2xl font-extrabold text-white tracking-tight">
            Escolha seu plano e comece agora
          </h1>
          <p className="text-xs text-slate-400">
            Acesso instantâneo com ativação imediata na sua conta.
          </p>
        </div>

        {/* Alternador de Forma de Pagamento */}
        <div className="inline-flex items-center gap-1.5 p-1 bg-slate-900 border border-slate-800 rounded-xl shadow-lg">
          <button
            onClick={() => setSelectedMethod("pix")}
            className={`flex items-center gap-1.5 px-4 py-1.5 rounded-lg text-xs font-bold transition-all ${
              selectedMethod === "pix"
                ? "bg-emerald-500 text-white shadow-md shadow-emerald-500/30 scale-105"
                : "text-slate-400 hover:text-white hover:bg-slate-800/50"
            }`}
          >
            <QrCode className="w-4 h-4" />
            Pix
          </button>
          <button
            onClick={() => setSelectedMethod("credit-card")}
            className={`flex items-center gap-1.5 px-4 py-1.5 rounded-lg text-xs font-bold transition-all ${
              selectedMethod === "credit-card"
                ? "bg-indigo-600 text-white shadow-md shadow-indigo-600/30 scale-105"
                : "text-slate-400 hover:text-white hover:bg-slate-800/50"
            }`}
          >
            <CreditCard className="w-4 h-4" />
            Cartão de Crédito
          </button>
        </div>
      </header>

      {/* Grid Horizontal de Cards */}
      <main className="max-w-6xl w-full mx-auto grid grid-cols-1 md:grid-cols-3 gap-3 lg:gap-5 items-stretch">
        {plans.map((plan) => (
          <div
            key={plan.id}
            className={`relative rounded-xl bg-slate-900/80 border p-4 flex flex-col justify-between transition-all ${
              plan.popular
                ? "border-indigo-500 shadow-lg shadow-indigo-500/10 ring-1 ring-indigo-500/50 bg-gradient-to-b from-indigo-950/30 to-slate-900"
                : "border-slate-800"
            }`}
          >
            {plan.popular && (
              <span className="absolute -top-2.5 left-1/2 -translate-x-1/2 px-2.5 py-0.5 rounded-full bg-indigo-600 text-white text-[9px] font-bold uppercase tracking-wider shadow-sm">
                Mais Popular
              </span>
            )}

            {/* Conteúdo do Card */}
            <div>
              <div className="flex justify-between items-baseline">
                <h2 className="text-base font-bold text-white">{plan.name}</h2>
                <div className="flex items-baseline gap-0.5">
                  <span className="text-[10px] font-semibold text-slate-400">
                    R$
                  </span>
                  <span className="text-xl font-black text-white">
                    {plan.price}
                  </span>
                  <span className="text-[10px] text-slate-400">
                    {plan.period}
                  </span>
                </div>
              </div>
              <p className="text-[11px] text-slate-400 mt-0.5">
                {plan.description}
              </p>

              {/* Lista de Features */}
              <ul className="mt-3 space-y-1.5 text-[11px] text-slate-300">
                {plan.features.map((feature, idx) => (
                  <li key={idx} className="flex items-center gap-2">
                    <Check className="w-3 h-3 text-emerald-400 shrink-0" />
                    <span className="truncate">{feature}</span>
                  </li>
                ))}
              </ul>
            </div>

            {/* Botão de Ação */}
            <div className="mt-4 pt-3 border-t border-slate-800">
              <button
                onClick={() => handleCheckout(plan.id)}
                className={`w-full py-2 px-3 rounded-lg text-xs font-bold flex items-center justify-center gap-1.5 transition-all ${
                  plan.popular
                    ? "bg-indigo-600 hover:bg-indigo-500 text-white shadow-md shadow-indigo-600/25"
                    : "bg-slate-800 hover:bg-slate-700 text-white border border-slate-700"
                }`}
              >
                <Zap className="w-3.5 h-3.5 fill-current" />
                Assinar via {selectedMethod === "pix" ? "Pix" : "Cartão"}
              </button>
            </div>
          </div>
        ))}
      </main>
    </div>
  );
}
