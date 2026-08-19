import { useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import {
  CreditCard,
  QrCode,
  ShieldCheck,
  Zap,
  ArrowLeft,
  Copy,
  Check,
  Wallet,
  AlertCircle,
  XCircle,
} from "lucide-react";
import { useAuthStore } from "../store/useAuthStore";
import api from "../services/api";

const PLAN_DETAILS = {
  "7_days": { name: "Plano Semanal", price: "19,99", rawPrice: 19.99 },
  "1_month": { name: "Plano Mensal", price: "34,99", rawPrice: 34.99 },
  "1_year": { name: "Plano Anual", price: "289,00", rawPrice: 289.0 },
};

const VALID_ACTIONS = ["add-balance", "new-subscription", "buy-product"];
const VALID_METHOD_PAYMENT = ["pix", "credit-card", "internal-balance"];

export default function Checkout() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const userId = useAuthStore((state) => state.userId);
  const [userBalanceCents, setUserBalanceCents] = useState(0);

  useEffect(() => {
    sessionStorage.removeItem("pendingCheckoutPlan");
    sessionStorage.removeItem("pendingCheckoutMethod");

    const fetchUser = async () => {
      try {
        const resp = await api.get("/api/user/find");
        const data = resp.data;
        setUserBalanceCents(data?.user?.balance || 0);
      } catch (error) {
        console.error("fetch get user balance: ", error);
      }
    };

    fetchUser();
  }, []);

  const userBalance = new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(userBalanceCents / 100);

  const actionParam = searchParams.get("action");
  const plan = searchParams.get("plan");
  const initialMethod = searchParams.get("method");
  const amountParam = searchParams.get("amount");
  const productIdParam = searchParams.get("product_id");

  const isValidAction = VALID_ACTIONS.includes(actionParam);
  const isValidPaymentMethod = VALID_METHOD_PAYMENT.includes(initialMethod);
  const isBalanceMode = actionParam === "add-balance";

  const currentPlan = PLAN_DETAILS[plan];

  const rawPriceInReais = isBalanceMode
    ? parseFloat(amountParam || "0")
    : currentPlan?.rawPrice || 0;

  const rawPriceInCents = Math.round(rawPriceInReais * 100);

  const hasEnoughBalance = userBalanceCents >= rawPriceInCents;

  const displayTitle = isBalanceMode
    ? "Adição de Saldo"
    : currentPlan?.name || "Assinatura Premium";

  const displayPrice = rawPriceInReais.toLocaleString("pt-BR", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });

  const [paymentMethod, setPaymentMethod] = useState(() => {
    if (isBalanceMode && initialMethod === "internal-balance") {
      return "pix";
    }
    return initialMethod === "card" ? "credit-card" : initialMethod;
  });

  const activePaymentMethod =
    isBalanceMode && paymentMethod === "internal-balance"
      ? "pix"
      : paymentMethod;

  const [loading, setLoading] = useState(false);
  const [pixData, setPixData] = useState(null);
  const [copied, setCopied] = useState(false);

  // Formulário do Cartão
  const [cardForm, setCardForm] = useState({
    number: "",
    name: "",
    expiry: "",
    cvv: "",
  });

  // Formatadores de Entrada (Máscaras)
  const handleInputChange = (e) => {
    let { name, value } = e.target;

    if (name === "number") {
      value = value
        .replace(/\D/g, "")
        .replace(/(\d{4})(?=\d)/g, "$1 ")
        .trim()
        .slice(0, 19);
    } else if (name === "expiry") {
      value = value
        .replace(/\D/g, "")
        .replace(/(\d{2})(\d)/, "$1/$2")
        .slice(0, 5);
    } else if (name === "cvv") {
      value = value.replace(/\D/g, "").slice(0, 4);
    }

    setCardForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleCopyPix = () => {
    if (pixData?.payload) {
      navigator.clipboard.writeText(pixData.payload);
      setCopied(true);
      setTimeout(() => setCopied(false), 3000);
    }
  };

  // Montador do corpo da requisição
  const buildRequestBody = () => {
    let payload;

    if (actionParam === "add-balance") {
      payload = {
        amount: rawPriceInCents,
        payment_method: activePaymentMethod,
        user_id: userId,
      };
    } else if (actionParam === "buy-product") {
      payload = {
        products: [{ product_id: productIdParam }],
        user_id: userId,
        payment_method: activePaymentMethod,
      };
    } else {
      payload = {
        plan: plan,
        method_payment: activePaymentMethod,
        user_id: userId,
      };
    }

    return { action: actionParam, payload };
  };

  const handleProcessPayment = async (e) => {
    e.preventDefault();
    setLoading(true);

    const requestBody = buildRequestBody();

    try {
      const response = await api.post("/api/orders", requestBody);

      console.log(response);

      if (paymentMethod === "pix") {
        console.log("aqui mostra o pixs");
      } else {
        alert("Operação realizada com sucesso!");
        navigate("/home");
      }
    } catch (error) {
      alert(error.message || "Erro ao processar pagamento.");
    } finally {
      setLoading(false);
    }
  };

  // TELA DE ERRO: Ação Inválida ou Ausente
  if (!isValidAction || !isValidPaymentMethod) {
    return (
      <div className="flex-1 bg-slate-950 text-slate-100 flex items-center justify-center p-4 min-h-screen">
        <div className="max-w-md w-full bg-slate-900 border border-slate-800 rounded-2xl p-6 md:p-8 text-center space-y-6 shadow-2xl">
          <div className="w-16 h-16 rounded-2xl bg-red-500/10 border border-red-500/20 flex items-center justify-center text-red-400 mx-auto shadow-inner">
            <XCircle className="w-8 h-8" />
          </div>
          <div className="space-y-2">
            <h2 className="text-xl font-bold text-white">
              Ação de Checkout Inválida
            </h2>
            <p className="text-xs text-slate-400 bg-red-500/10 border border-red-500/20 p-3 rounded-xl">
              Não foi possível identificar o tipo de operação. Certifique-se de
              acessar o checkout a partir de um plano ou recarga válida.
            </p>
          </div>
          <button
            onClick={() => navigate("/home")}
            className="w-full py-3 px-4 bg-slate-800 hover:bg-slate-700 text-white font-bold text-xs rounded-xl transition-all border border-slate-700 cursor-pointer"
          >
            Voltar para o Início
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 bg-slate-950 text-slate-100 flex flex-col items-center justify-center p-4 md:p-6">
      <div className="max-w-4xl w-full grid grid-cols-1 md:grid-cols-12 gap-6 bg-slate-900 border border-slate-800 rounded-2xl p-5 md:p-8 shadow-2xl">
        {/* Resumo do Pedido & Seleção de Método */}
        <div className="md:col-span-5 space-y-6 border-b md:border-b-0 md:border-r border-slate-800 pb-6 md:pb-0 md:pr-6 flex flex-col justify-between">
          <div className="space-y-4">
            <button
              onClick={() => navigate(-1)}
              className="inline-flex items-center gap-1.5 text-xs text-slate-400 hover:text-white transition-colors cursor-pointer"
            >
              <ArrowLeft className="w-3.5 h-3.5" /> Voltar
            </button>

            <div>
              <span className="text-xs font-semibold text-indigo-400 uppercase tracking-wider">
                Resumo do Pedido
              </span>
              <h1 className="text-xl font-bold text-white mt-1">
                {displayTitle}
              </h1>
            </div>

            <div className="bg-slate-950/60 p-4 rounded-xl border border-slate-800 space-y-2">
              <div className="flex justify-between text-xs text-slate-400">
                <span>Item</span>
                <span className="text-white font-medium">{displayTitle}</span>
              </div>
              <div className="flex justify-between items-baseline pt-2 border-t border-slate-800/80">
                <span className="text-sm font-semibold text-white">Total</span>
                <div className="text-right">
                  <span className="text-xs text-slate-400">R$ </span>
                  <span className="text-2xl font-extrabold text-white">
                    {displayPrice}
                  </span>
                </div>
              </div>
            </div>

            {/* Seleção do Método de Pagamento */}
            <div className="space-y-2 pt-2">
              <label className="text-xs font-medium text-slate-300">
                Forma de Pagamento
              </label>
              <div className="grid grid-cols-1 gap-2 p-1 bg-slate-950 rounded-xl border border-slate-800">
                {/* Opção Pix */}
                <button
                  type="button"
                  onClick={() => {
                    setPaymentMethod("pix");
                    setPixData(null);
                  }}
                  className={`flex items-center justify-between px-3 py-2.5 rounded-lg text-xs font-bold transition-all cursor-pointer ${
                    paymentMethod === "pix"
                      ? "bg-emerald-500 text-white shadow-md shadow-emerald-500/20"
                      : "text-slate-400 hover:text-white"
                  }`}
                >
                  <span className="flex items-center gap-2">
                    <QrCode className="w-4 h-4" /> Pix
                  </span>
                  <span className="text-[10px] opacity-80 font-normal">
                    Aprovação Instantânea
                  </span>
                </button>

                {/* Opção Cartão de Crédito */}
                <button
                  type="button"
                  onClick={() => setPaymentMethod("credit-card")}
                  className={`flex items-center justify-between px-3 py-2.5 rounded-lg text-xs font-bold transition-all cursor-pointer ${
                    paymentMethod === "credit-card"
                      ? "bg-indigo-600 text-white shadow-md shadow-indigo-600/20"
                      : "text-slate-400 hover:text-white"
                  }`}
                >
                  <span className="flex items-center gap-2">
                    <CreditCard className="w-4 h-4" /> Cartão de Crédito
                  </span>
                </button>

                {/* Opção Saldo Interno (Apenas se NÃO for 'add-balance') */}
                {!isBalanceMode && (
                  <button
                    type="button"
                    onClick={() => setPaymentMethod("internal-balance")}
                    className={`flex items-center justify-between px-3 py-2.5 rounded-lg text-xs font-bold transition-all cursor-pointer ${
                      paymentMethod === "internal-balance"
                        ? "bg-amber-500 text-slate-950 shadow-md shadow-amber-500/20"
                        : "text-slate-400 hover:text-white"
                    }`}
                  >
                    <span className="flex items-center gap-2">
                      <Wallet className="w-4 h-4" /> Saldo Interno
                    </span>
                    <span className="text-[11px]">R$ {userBalance}</span>
                  </button>
                )}
              </div>
            </div>
          </div>

          <div className="flex items-center gap-2 text-xs text-slate-500 pt-4">
            <ShieldCheck className="w-4 h-4 text-emerald-400 shrink-0" />
            <span>Pagamento 100% seguro com criptografia de ponta.</span>
          </div>
        </div>

        {/* Formulário / Processador */}
        <div className="md:col-span-7 flex flex-col justify-center pl-0 md:pl-2">
          {/* PAINEL 1: CARTÃO DE CRÉDITO */}
          {paymentMethod === "credit-card" && (
            <form onSubmit={handleProcessPayment} className="space-y-4">
              <h2 className="text-lg font-bold text-white mb-2">
                Dados do Cartão
              </h2>

              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">
                  Número do Cartão
                </label>
                <input
                  type="text"
                  name="number"
                  placeholder="0000 0000 0000 0000"
                  required
                  value={cardForm.number}
                  onChange={handleInputChange}
                  className="w-full h-10 px-3 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white placeholder-slate-500 outline-none focus:border-indigo-500"
                />
              </div>

              <div>
                <label className="block text-xs font-medium text-slate-300 mb-1">
                  Nome do Titular
                </label>
                <input
                  type="text"
                  name="name"
                  placeholder="NOME COMPLETO"
                  required
                  value={cardForm.name}
                  onChange={handleInputChange}
                  className="w-full h-10 px-3 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white placeholder-slate-500 outline-none focus:border-indigo-500 uppercase"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-medium text-slate-300 mb-1">
                    Validade
                  </label>
                  <input
                    type="text"
                    name="expiry"
                    placeholder="MM/AA"
                    required
                    value={cardForm.expiry}
                    onChange={handleInputChange}
                    className="w-full h-10 px-3 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white placeholder-slate-500 outline-none focus:border-indigo-500"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-slate-300 mb-1">
                    CVV
                  </label>
                  <input
                    type="text"
                    name="cvv"
                    placeholder="123"
                    required
                    value={cardForm.cvv}
                    onChange={handleInputChange}
                    className="w-full h-10 px-3 bg-slate-950 border border-slate-800 rounded-xl text-sm text-white placeholder-slate-500 outline-none focus:border-indigo-500"
                  />
                </div>
              </div>

              <button
                type="submit"
                disabled={loading}
                className="w-full mt-2 py-3 px-4 bg-indigo-600 hover:bg-indigo-500 text-white font-bold rounded-xl text-sm transition-all shadow-lg shadow-indigo-600/25 flex items-center justify-center gap-2 disabled:opacity-50 cursor-pointer"
              >
                <Zap className="w-4 h-4 fill-current" />
                {loading ? "Processando..." : `Pagar R$ ${displayPrice}`}
              </button>
            </form>
          )}

          {/* PAINEL 2: PIX */}
          {paymentMethod === "pix" && (
            <div className="space-y-5 text-center">
              <h2 className="text-lg font-bold text-white">
                Pagamento via Pix
              </h2>

              {!pixData ? (
                <div className="py-6 space-y-4">
                  <p className="text-xs text-slate-400 max-w-xs mx-auto">
                    Aprovação instantânea. Clique no botão abaixo para gerar o
                    código QR Code e Copia e Cola.
                  </p>
                  <button
                    type="button"
                    onClick={handleProcessPayment}
                    disabled={loading}
                    className="w-full py-3 px-4 bg-emerald-500 hover:bg-emerald-400 text-white font-bold rounded-xl text-sm transition-all shadow-lg shadow-emerald-500/25 flex items-center justify-center gap-2 disabled:opacity-50 cursor-pointer"
                  >
                    <QrCode className="w-4 h-4" />
                    {loading
                      ? "Gerando PIX..."
                      : `Gerar PIX de R$ ${displayPrice}`}
                  </button>
                </div>
              ) : (
                <div className="space-y-4 animate-fade-in">
                  <div className="bg-white p-3 rounded-xl inline-block border border-slate-700 shadow-md">
                    <img
                      src={pixData.qrCodeUrl}
                      alt="QR Code Pix"
                      className="w-44 h-44 object-contain mx-auto"
                    />
                  </div>

                  <div className="space-y-2">
                    <p className="text-xs text-slate-400">
                      Ou copie a chave Pix abaixo:
                    </p>
                    <div className="flex items-center gap-2 max-w-sm mx-auto">
                      <input
                        type="text"
                        readOnly
                        value={pixData.payload}
                        className="w-full h-9 px-3 bg-slate-950 border border-slate-800 rounded-lg text-xs text-slate-300 outline-none truncate"
                      />
                      <button
                        type="button"
                        onClick={handleCopyPix}
                        className="h-9 px-3 bg-emerald-500 hover:bg-emerald-400 text-white font-bold rounded-lg text-xs flex items-center gap-1 transition-colors shrink-0 cursor-pointer"
                      >
                        {copied ? (
                          <Check className="w-3.5 h-3.5" />
                        ) : (
                          <Copy className="w-3.5 h-3.5" />
                        )}
                        {copied ? "Copiado" : "Copiar"}
                      </button>
                    </div>
                  </div>
                </div>
              )}
            </div>
          )}

          {/* PAINEL 3: SALDO INTERNO (Apenas se NÃO for 'add-balance') */}
          {!isBalanceMode && paymentMethod === "internal-balance" && (
            <div className="space-y-5 text-center">
              <h2 className="text-lg font-bold text-white">
                Pagamento com Saldo Interno
              </h2>

              <div className="bg-slate-950 border border-slate-800 p-4 rounded-xl space-y-3 text-left">
                <div className="flex justify-between text-xs">
                  <span className="text-slate-400">Seu saldo disponível:</span>
                  <span className="text-white font-bold">{userBalance}</span>
                </div>
                <div className="flex justify-between text-xs">
                  <span className="text-slate-400">Valor da compra:</span>
                  <span className="text-amber-400 font-bold">
                    R$ {displayPrice}
                  </span>
                </div>
              </div>

              {!hasEnoughBalance ? (
                <div className="flex items-center gap-2 p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-red-400 text-xs">
                  <AlertCircle className="w-4 h-4 shrink-0" />
                  <span>
                    Saldo insuficiente para realizar esta compra. Faça uma
                    recarga via Pix ou Cartão.
                  </span>
                </div>
              ) : (
                <button
                  type="button"
                  onClick={handleProcessPayment}
                  disabled={loading}
                  className="w-full py-3 px-4 bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold rounded-xl text-sm transition-all shadow-lg shadow-amber-500/25 flex items-center justify-center gap-2 disabled:opacity-50 cursor-pointer"
                >
                  <Wallet className="w-4 h-4" />
                  {loading
                    ? "Debitando Saldo..."
                    : "Confirmar e Pagar com Saldo"}
                </button>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
