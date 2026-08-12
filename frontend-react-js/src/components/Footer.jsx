import { useState, useEffect } from "react";
import { Users } from "lucide-react";

export default function Footer() {
  const [onlineUsers, setOnlineOusers] = useState(0);

  useEffect(() => {
    const src = new EventSource("/api/user/stats-sse");

    src.addEventListener("conns", (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log("coon sse active received data: ", data);
        setOnlineOusers(data.conns);
      } catch (parseError) {
        console.log("error parse json conns sse", parseError);
      }
    });

    src.addEventListener("error", (event) => {
      if (src.readyState === EventSource.CLOSED) {
        console.error("connection SSE encerrada");
      }

      console.log(event.data);
    });

    src.onerror = (error) => {
      console.error("error conn sse: ", error);
    };

    return () => {
      src.close();
    };
  }, []);

  return (
    <footer className="bg-slate-900 border-c border-slate-800 text-slate-400 py-4 px-6 text-xs w-full mt-auto shadow-inner">
      <div className="max-w-6xl mx-auto flex flex-col sm:flex-row items-center justify-center gap-3">
        {/* Direitos Reservados */}
        <p className="text-slate-400 font-medium">
          &copy; 2026{" "}
          <span className="text-slate-200 font-semibold">
            Simple Api Gateway
          </span>
          . Todos os direitos reservados.
        </p>

        {/* Usuários Online em Destaque Lado a Lado */}
        <div className="flex items-center gap-2 font-semibold text-slate-300 bg-slate-950 border border-slate-800 px-3 py-1.5 rounded-full shadow-sm">
          <span className="relative flex h-2 w-2">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
          </span>
          <Users className="w-3.5 h-3.5 text-indigo-400" />
          <span>
            {onlineUsers}{" "}
            {onlineUsers === 1 ? "usuário online" : "usuários online"}
          </span>
        </div>
      </div>
    </footer>
  );
}
