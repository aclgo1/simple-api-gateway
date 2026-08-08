import { useState, useEffect } from "react";
import { FaUser } from "react-icons/fa";

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
    <div className="text-center">
      <footer className="bg-slate-900 text-white p-6 text-center text-sm w-full mt-auto shadow-inner">
        <div className="flex flex-col items-center gap-2">
          <p>&copy; 2026 Simple Api Gateway. Todos os direitos reservados.</p>
          <p className="flex items-center gap-2 font-medium text-slate-300 bg-slate-800 px-3 py-1 rounded-full w-fit mx-auto">
            <FaUser className="text-blue-400"></FaUser>
            <span>
              {onlineUsers}{" "}
              {onlineUsers == 1 ? "usúario online" : "usúarios online"}
            </span>
          </p>
        </div>
      </footer>
    </div>
  );
}
