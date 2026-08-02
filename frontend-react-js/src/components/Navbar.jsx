import { Link, useNavigate } from "react-router-dom";
import { useAuthStore } from "../store/useAuthStore";

export default function Navbar() {
  const logout = useAuthStore((state) => state.logout);
  const accessToken = useAuthStore((state) => state.accessToken);
  const navigate = useNavigate();

  const handleLogout = (e) => {
    e.preventDefault();
    logout();
    navigate("/login");
  };

  return (
    <nav className="bg-slate-900 text-white shadow-md w-full">
      <div className="max-w-6xl mx-auto w-full p-4 flex justify-center items-center">
        <div className="flex gap-4 items-center">
          {!accessToken ? (
            <>
              <Link to="/" className="hover:underline">
                Pricing
              </Link>
              <Link to="/login" className="hover:underline">
                Login
              </Link>
            </>
          ) : (
            <>
              <Link to="/home" className="hover:underline">
                Home
              </Link>
              <Link to="/" className="hover:underline">
                Pricing
              </Link>
              <Link to="/profile" className="hover:underline">
                Profile
              </Link>
              <Link to="/add-balance" className="hover:underline">
                Add Balance
              </Link>

              <button
                onClick={handleLogout}
                className="hover:underline bg-transparent border-none text-white cursor-pointer font-sans text-base"
              >
                Logout
              </button>
            </>
          )}
        </div>
      </div>
    </nav>
  );
}
