import { useNavigate } from "react-router-dom";
import { useForm, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import {
  Wallet,
  QrCode,
  CreditCard,
  ArrowRight,
  AlertCircle,
  CheckCircle2,
} from "lucide-react";

const AMOUNT_OPTIONS = [
  { cents: 2500, label: "R$ 25,00", real: "25.00" },
  { cents: 5000, label: "R$ 50,00", real: "50.00" },
  { cents: 10000, label: "R$ 100,00", real: "100.00" },
  // { cents: 100000, label: "R$ 1000,00", real: "1000.00" },
];

const addBalanceSchema = z.object({
  amount: z.string().min(1, "Selecione um valor para a recarga"),
  method: z.enum(["pix", "credit-card"], {
    required_error: "Selecione a forma de pagamento",
  }),
});

export default function AddBalance() {
  const navigate = useNavigate();

  const {
    control,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm({
    resolver: zodResolver(addBalanceSchema),
    defaultValues: {
      amount: "25.00",
      method: "pix",
    },
  });

  const selectedAmount = useWatch({ control, name: "amount" });
  const selectedMethod = useWatch({ control, name: "method" });

  const onSubmit = (data) => {
    navigate(
      `/checkout?action=add-balance&method=${data.method}&amount=${data.amount}`,
    );
  };

  return (
    /* REMOVIDO min-h-screen */
    <div className="flex-1 bg-slate-950 text-slate-100 flex items-center justify-center p-4 md:p-6 my-auto">
      {/* CORRIGIDO: max-w-md w-full (removido flex-1 daqui) */}
      <div className="max-w-md w-full bg-slate-900 border border-slate-800 rounded-2xl p-6 md:p-8 space-y-6 shadow-2xl">
        <div className="text-center space-y-2">
          <div className="w-16 h-16 rounded-2xl bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center text-indigo-400 mx-auto shadow-inner">
            <Wallet className="w-8 h-8" />
          </div>
          <h2 className="text-xl font-bold text-white">Adicionar Saldo</h2>
          <p className="text-xs text-slate-400">
            Selecione um dos valores disponíveis para recarregar sua conta.
          </p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <div className="space-y-2">
            <label className="text-xs font-medium text-slate-300">
              Selecione o Valor
            </label>

            <div className="grid grid-cols-2 gap-3">
              {AMOUNT_OPTIONS.map((option) => {
                const isSelected = selectedAmount === option.real;
                return (
                  <button
                    key={option.cents}
                    type="button"
                    onClick={() => setValue("amount", option.real)}
                    className={`relative p-4 rounded-xl border text-center transition-all cursor-pointer flex flex-col items-center justify-center ${
                      isSelected
                        ? "bg-indigo-600/10 border-indigo-500 text-indigo-400 shadow-md shadow-indigo-500/10"
                        : "bg-slate-950 border-slate-800 text-slate-300 hover:border-slate-700 hover:bg-slate-900"
                    }`}
                  >
                    {isSelected && (
                      <CheckCircle2 className="w-4 h-4 text-indigo-400 absolute top-2 right-2" />
                    )}
                    <span className="text-lg font-bold font-mono">
                      {option.label}
                    </span>
                  </button>
                );
              })}
            </div>

            {errors.amount && (
              <p className="text-xs text-red-400 flex items-center gap-1">
                <AlertCircle className="w-3.5 h-3.5 shrink-0" />
                {errors.amount.message}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <label className="text-xs font-medium text-slate-300">
              Forma de Pagamento
            </label>

            <div className="grid grid-cols-2 gap-3">
              <button
                type="button"
                onClick={() => setValue("method", "pix")}
                className={`p-3.5 rounded-xl border transition-all flex flex-col items-center justify-center gap-2 cursor-pointer ${
                  selectedMethod === "pix"
                    ? "bg-indigo-600/10 border-indigo-500 text-indigo-400"
                    : "bg-slate-950 border-slate-800 text-slate-400 hover:border-slate-700"
                }`}
              >
                <QrCode className="w-6 h-6" />
                <span className="text-xs font-medium">PIX</span>
              </button>

              <button
                type="button"
                onClick={() => setValue("method", "credit-card")}
                className={`p-3.5 rounded-xl border transition-all flex flex-col items-center justify-center gap-2 cursor-pointer ${
                  selectedMethod === "credit-card"
                    ? "bg-indigo-600/10 border-indigo-500 text-indigo-400"
                    : "bg-slate-950 border-slate-800 text-slate-400 hover:border-slate-700"
                }`}
              >
                <CreditCard className="w-6 h-6" />
                <span className="text-xs font-medium">Cartão</span>
              </button>
            </div>

            {errors.method && (
              <p className="text-xs text-red-400 flex items-center gap-1">
                <AlertCircle className="w-3.5 h-3.5 shrink-0" />
                {errors.method.message}
              </p>
            )}
          </div>

          <button
            type="submit"
            className="w-full py-3.5 px-4 bg-indigo-600 hover:bg-indigo-500 text-white font-bold rounded-xl text-sm transition-all shadow-lg shadow-indigo-600/25 flex items-center justify-center gap-2 cursor-pointer"
          >
            <span>Ir para o Pagamento</span>
            <ArrowRight className="w-4 h-4" />
          </button>
        </form>
      </div>
    </div>
  );
}
