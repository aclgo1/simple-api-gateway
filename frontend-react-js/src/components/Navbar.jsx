import { Link } from "react-router-dom";
import { useAuthStore } from "../store/useAuthStore";
import { LogIn, LogOut, CreditCard, User, Home, Shield } from "lucide-react";
import { handleForceLogout } from "../services/api";
export default function Navbar() {
  // const logout = useAuthStore((state) => state.logout);
  const accessToken = useAuthStore((state) => state.accessToken);
  // const navigate = useNavigate();

  const handleLogout = async (e) => {
    e.preventDefault();
    await handleForceLogout();
    // navigate("/login");
  };

  return (
    <nav className="bg-slate-900 border-b border-slate-800 text-slate-100 w-full shadow-md">
      <div className="max-w-6xl mx-auto px-4 py-3 flex flex-wrap justify-between items-center gap-4">
        {/* Brand / Logo */}
        <Link
          to="/"
          className="flex items-center gap-2 font-extrabold text-base tracking-wide text-white hover:text-indigo-400 transition-colors"
        >
          <Shield className="w-5 h-5 text-indigo-500" />
          <span>Simple API Gateway</span>
        </Link>

        {/* Navigation Links Lado a Lado */}
        <div className="flex items-center gap-1 sm:gap-2 text-xs font-semibold">
          {!accessToken ? (
            <>
              <Link
                to="/"
                className="px-3 py-1.5 rounded-lg text-slate-300 hover:text-white hover:bg-slate-800 transition-all"
              >
                Pricing
              </Link>
              <Link
                to="/login"
                className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-indigo-600 hover:bg-indigo-500 text-white font-bold transition-all shadow-sm shadow-indigo-600/20"
              >
                <LogIn className="w-3.5 h-3.5" />
                Login
              </Link>
            </>
          ) : (
            <>
              <Link
                to="/home"
                className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-slate-300 hover:text-white hover:bg-slate-800 transition-all"
              >
                <Home className="w-3.5 h-3.5 text-indigo-400" />
                Home
              </Link>

              <Link
                to="/"
                className="px-3 py-1.5 rounded-lg text-slate-300 hover:text-white hover:bg-slate-800 transition-all"
              >
                Pricing
              </Link>

              <Link
                to="/profile"
                className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-slate-300 hover:text-white hover:bg-slate-800 transition-all"
              >
                <User className="w-3.5 h-3.5 text-indigo-400" />
                Profile
              </Link>

              <Link
                to="/add-balance"
                className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-amber-400 hover:text-amber-300 hover:bg-slate-800 transition-all font-bold"
              >
                <CreditCard className="w-3.5 h-3.5" />
                Add Balance
              </Link>

              <button
                onClick={handleLogout}
                className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-red-400 hover:text-red-300 hover:bg-red-500/10 transition-all font-bold ml-1"
              >
                <LogOut className="w-3.5 h-3.5" />
                Logout
              </button>
            </>
          )}
        </div>
      </div>
    </nav>
  );
}
