import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { loginUser } from "@/auth/api/api";
import type { ResponseType } from "@/auth/api/types";
import Button from "@/components/AuthComponents/Button";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [token, setToken] = useState<string | null>(null);
  const navigate = useNavigate();
  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    try {
      const response = (await loginUser({ email, password })) as ResponseType;
      if (response.success) {
        const data = response.data as { token: string };
        setToken(data.token);
        navigate("/home");
      } else {
        setError(String(response.message));
      }
    } catch (err) {
      setError((err as Error).message);
    }
  };
  if (token) {
    return <div>Login successful!</div>;
  }
  return (
    <main className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10 relative card">
      <div className="w-full max-w-md sm:max-w-lg border p-6 sm:p-8 rounded-2xl bg-accent drop-shadow-red-500 drop-shadow-sm flex flex-col gap-6">
        <div className="flex flex-col items-center justify-center gap-2">
          <h1 className="text-3xl sm:text-4xl p-2 uppercase drop-shadow-red-500 drop-shadow-sm">
            red404
          </h1>
          <h3 className="text-muted-foreground text-center">Welcome back</h3>
        </div>

        <form
          className="w-full flex flex-col items-center pb-4 sm:pb-6"
          onSubmit={handleLogin}
        >
          <div className="w-full flex flex-col gap-3 sm:gap-4 p-2 sm:p-0">
            <label className="self-start" htmlFor="email">
              Email
            </label>
            <input
              type="text"
              id="email"
              placeholder="Enter your email"
              className="w-full rounded-2xl border border-accent-secondary py-2 px-3 focus:outline-accent-secondary hover:bg-accent-secondary transition-all duration-300"
              value={email}
              required
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setEmail(e.target.value)
              }
            />
            <label className="self-start" htmlFor="password">
              Password
            </label>
            <input
              type="password"
              id="password"
              placeholder="Enter your password"
              className="w-full rounded-2xl border border-accent-secondary py-2 px-3 focus:outline-accent-secondary hover:bg-accent-secondary transition-all duration-300"
              value={password}
              minLength={8}
              required
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setPassword(e.target.value)
              }
            />
            <Button
              type="submit"
              text="Login"
              className="w-full hover:bg-accent-secondary my-2 "
            />
          </div>
          <p className="text-muted-foreground text-center">
            Don't have an account?{" "}
            <Link
              to="/register"
              className="hover:underline hover:text-red-500 transition-all duration-300"
            >
              Register
            </Link>
          </p>
        </form>
        {error && (
          <span className="mt-2 block text-center error text-red-800">
            {error}
          </span>
        )}
      </div>
    </main>
  );
}
