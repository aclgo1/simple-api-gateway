import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useQuery, useMutation } from "@tanstack/react-query";

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
    window.location.href = "/home";
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

      return response.data;
    },
    onSuccess: (data) => {
      // Corrigido: parênteses do alert
      alert(data.message || "Conta criada com sucesso");
      handleAuthSuccess(
        data.created.tokens.access_token,
        data.created.tokens.refresh_token,
      );
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
      className="min-h-screen w-full flex items-center justify-center bg-cover bg-center bg-fixed"
      style={{
        backgroundImage: 'url("/background.jpg")',
      }}
    >
      <h1 className="text-white text-3xl font-bold">Login</h1>
    </div>
  );
}
