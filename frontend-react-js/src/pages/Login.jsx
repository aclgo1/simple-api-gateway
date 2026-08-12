import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useQuery, useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";

import api from "../services/api";
import { useAuthStore } from "../store/useAuthStore";

const loginSchema = z.object({
  email: z.string().email(1, "E-mail ou usuário é obrigatório"),
  password: z.string().min(1, "A senha é obrigatória"),
  captchaAnswer: z.string().min(1, "Digite o código da imagem"),
});

const signUpSchema = z.object({
  name: z.string().min(1, "Nome é obrigatório"),
  lastname: z.string().min(1, "Sobrenome é obrigatório"),
  email: z.string().email("E-mail inválido"),
  password: z.string().min(6, "A senha deve ter no mínimo 6 caracteres"),
  captchaAnswer: z.string().min(1, "Digite o código da imagem"),
});

export default function Login() {
  const [isSignUp, setIsSignUp] = useState(false);

  // Corrigido: setTokens (com 'o' minúsculo)
  const setTokens = useAuthStore((state) => state.setTokens);

  const navigate = useNavigate();

  const {
    data: captcha,
    isLoading: loadingCaptcha,
    refetch: refetchCaptcha,
  } = useQuery({
    queryKey: ["captcha"],
    queryFn: async () => {
      const response = await api.get("/api/captcha");
      return response.data;
    },
    refetchOnWindowFocus: false,
  });

  const handleAuthSuccess = (accessToken, refreshToken) => {
    setTokens(accessToken, refreshToken);

    const pendingPlan = sessionStorage.getItem("pendingCheckoutPlan");
    const pendingMethod = sessionStorage.getItem("pendingCheckoutMethod");

    if (pendingPlan && pendingMethod) {
      sessionStorage.removeItem("pendingCheckoutPlan");
      sessionStorage.removeItem("pendingCheckoutMethod");
      navigate(`/checkout?plan=${pendingPlan}&method=${pendingMethod}`);
    } else {
      navigate("/home");
    }
  };

  const loginMutation = useMutation({
    mutationFn: async (formData) => {
      const response = await api.post("/api/login", {
        email: formData.email,
        password: formData.password,
        captcha_id: captcha?.id,
        // Corrigido: captchaAnswer
        captcha_awnser: formData.captchaAnswer,
      });

      return response.data;
    },
    onSuccess: (data) => {
      alert("Login realizado com sucesso!");
      handleAuthSuccess(data.tokens.access_token, data.tokens.refresh_token);
    },
    onError: (error) => {
      const message = error.message || "Erro ao realizar login.";
      alert(message);
      refetchCaptcha();
    },
  });

  const signUpMutation = useMutation({
    mutationFn: async (formData) => {
      const response = await api.post("/api/user/register", {
        name: formData.name,
        lastname: formData.lastname,
        email: formData.email,
        password: formData.password,
        captcha_id: captcha?.id,
        captcha_awnser: formData.captchaAnswer,
      });

      return { responseDat: response.data, userEmail: formData.email };
    },
    onSuccess: (data) => {
      // alert(data.message || "Conta criada com sucesso");
      navigate(`/confirm?email=${encodeURIComponent(data.userEmail)}`);
    },
    onError: (error) => {
      const message = error.message || "Erro ao criar conta.";
      alert(message);
      refetchCaptcha();
    },
  });

  const loginForm = useForm({ resolver: zodResolver(loginSchema) });
  const signUpForm = useForm({ resolver: zodResolver(signUpSchema) });

  const onLoginSubmit = (data) => loginMutation.mutate(data);
  const onSignUpSubmit = (data) => signUpMutation.mutate(data);

  return (
    <div
      className="flex-1 w-full flex items-center justify-center bg-cover bg-center bg-fixed"
      style={{
        backgroundImage: 'url("/background.jpg")',
      }}
    >
      <div className="w-[320px] p-6 rounded-2xl bg-black/40 backdrop-blur-md border border-white/20 shadow-2xl">
        {!isSignUp ? (
          <form
            onSubmit={loginForm.handleSubmit(onLoginSubmit)}
            className="flex flex-col gap-3"
          >
            <div>
              <input
                {...loginForm.register("email")}
                type="text"
                placeholder="Email"
                className="w-full h-10 px-3 bg-white/90 rounded-md text-sm outline-none focus:ring-2 focus:ring-indigo-500"
              ></input>
              {loginForm.formState.errors.email && (
                <span className="text=xs text-red-400 mt-1 block">
                  {loginForm.formState.errors.email.message}
                </span>
              )}
            </div>

            <div>
              <input
                {...loginForm.register("password")}
                type="password"
                placeholder="Password"
                className="w-full h-10 px-3 bg-white/90 rounded-md text-sm outline-none focus:ring-2 focus:ring-indigo-500"
              >
                {loginForm.formState.errors.password && (
                  <span className="text-xs text-red-400 mt-1 block">
                    {loginForm.formState.errors.password.message}
                  </span>
                )}
              </input>
            </div>

            <div className="text-center my-1">
              {loadingCaptcha ? (
                <div className="h-17.5 w-50 mx-auto bg-ray-300 animate-pulse rounded border border-gray-400"></div>
              ) : (
                <img
                  src={captcha?.base64_image}
                  alt="Captcha"
                  onClick={() => refetchCaptcha()}
                  className="h-17.5 w-50 mx-auto mb-2 border border-gray-300 bg-gray-100 cursor-pointer object-cover rounded"
                  title="Clique para trocar o Captcha"
                ></img>
              )}
              <input
                {...loginForm.register("captchaAnswer")}
                type="text"
                placeholder="Code"
                className="w-30.5 h-10 mx-auto px-3 bg-white/90 rounded-md text-sm outline-none block text-center"
              ></input>
              {loginForm.formState.errors.captchaAnswer && (
                <span className="text-xs text-red-400 mt-1 block">
                  {loginForm.formState.errors.captchaAnswer.message}
                </span>
              )}
            </div>

            <button
              type="submit"
              disabled={loginMutation.isPending}
              className="w-full h-10 bg-[#362982] hover:bg-[#2a2068] text-white font-bold rounded-md transition-colors disabled::opacity-50"
            >
              {loginMutation.isPending ? "Wait..." : "Login"}
            </button>

            <a
              href="/resetpass"
              className="text-center text-sm font-bolf text-white hover:underline mt-1 drop-shadow-[0_1.2px_1.2px_rgba(0,0,0,0.8)]"
            >
              Forgot Pass?
            </a>

            <button
              type="button"
              onClick={() => {
                setIsSignUp(true);
                signUpForm.reset();
              }}
              className="w-full h-10 bg-gray-600 hover:bg-gray-700 text-white font-bold rounded-md text-xs transition-colors mt-2"
            >
              Create Account
            </button>
          </form>
        ) : (
          <form
            onSubmit={signUpForm.handleSubmit(onSignUpSubmit)}
            className="flex flex-col gap-2"
          >
            <div>
              <input
                {...signUpForm.register("name")}
                type="text"
                placeholder="Name"
                className="w-full h-10 px-3 bg-white/90 rounded-md text-sm outline-none focus:ring-2 focus:ring-indigo-500"
              />
              {signUpForm.formState.errors.name && (
                <span className="text-xs text-red-400 block mt-1">
                  {signUpForm.formState.errors.name.message}
                </span>
              )}
            </div>

            <div>
              <input
                {...signUpForm.register("lastname")}
                type="text"
                placeholder="Lastname"
                className="w-full h-10 px-3 bg-white/90 rounded-md text-sm outline-none focus:ring-2 focus:ring-indigo-500"
              />
              {signUpForm.formState.errors.lastname && (
                <span className="text-xs text-red-400 block mt-1">
                  {signUpForm.formState.errors.lastname.message}
                </span>
              )}
            </div>

            <div>
              <input
                {...signUpForm.register("email")}
                type="email"
                placeholder="Email"
                className="w-full h-10 px-3 bg-white/90 rounded-md text-sm outline-none focus:ring-2 focus:ring-indigo-500"
              />
              {signUpForm.formState.errors.email && (
                <span className="text-xs text-red-400 block mt-1">
                  {signUpForm.formState.errors.email.message}
                </span>
              )}
            </div>

            <div>
              <input
                {...signUpForm.register("password")}
                type="password"
                placeholder="Password"
                className="w-full h-10 px-3 bg-white/90 rounded-md text-sm outline-none focus:ring-2 focus:ring-indigo-500"
              />
              {signUpForm.formState.errors.password && (
                <span className="text-xs text-red-400 block mt-1">
                  {signUpForm.formState.errors.password.message}
                </span>
              )}
            </div>

            <div className="text-center my-1">
              {loadingCaptcha ? (
                <div className="h-17.5 w-50 mx-auto bg-gray-300 animate-pulse rounded border border-gray-400" />
              ) : (
                <img
                  src={captcha?.base64_image}
                  alt="Captcha"
                  onClick={() => refetchCaptcha()}
                  className="h-17.5 w-50 mx-auto mb-2 border border-gray-300 bg-gray-100 cursor-pointer object-cover rounded"
                  title="Clique para trocar o Captcha"
                />
              )}
              <input
                {...signUpForm.register("captchaAnswer")}
                type="text"
                placeholder="Code"
                className="w-30.5 h-10 mx-auto px-3 bg-white/90 rounded-md text-sm outline-none block text-center"
              />
              {signUpForm.formState.errors.captchaAnswer && (
                <span className="text-xs text-red-400 block mt-1">
                  {signUpForm.formState.errors.captchaAnswer.message}
                </span>
              )}
            </div>

            <button
              type="submit"
              disabled={signUpMutation.isPending}
              className="w-full h-10 bg-[#362982] hover:bg-[#2a2068] text-white font-bold rounded-md transition-colors disabled:opacity-50"
            >
              {signUpMutation.isPending ? "Wait..." : "Sign Up"}
            </button>

            <button
              type="button"
              onClick={() => {
                setIsSignUp(false);
                loginForm.reset();
              }}
              className="w-full h-10 bg-gray-600 hover:bg-gray-700 text-white font-bold rounded-md text-xs transition-colors mt-2"
            >
              Back to Login
            </button>
          </form>
        )}
      </div>
    </div>
  );
}
