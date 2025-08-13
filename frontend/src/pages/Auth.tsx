import { useState } from "react";
import { useNavigate } from "react-router";
import GoogleIcon from "../assets/icons/GoogleIcon";
import MetaIcon from "../assets/icons/MetaIcon";
import { loginUser } from "../auth/api/login";
import { registerUser } from "../auth/api/register";
import googleAuth from "../auth/providers/googleProvider";
import metaAuth from "../auth/providers/metaProvider";
import Button from "../components/AuthComponents/Button";

function GoogleAuthButton() {
  return (
    <Button type="button" provider="Google" onClick={googleAuth}>
      <GoogleIcon />
    </Button>
  );
}

function MetaAuthButton() {
  return (
    <Button type="button" provider="Meta" onClick={metaAuth}>
      <MetaIcon />
    </Button>
  );
}
function Register() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    try {
      await registerUser(email, password);
    } catch (err) {
      setError((err as Error).message);
    }
  };
  return (
    <>
      <form
        className="flex flex-col items-center pb-5"
        onSubmit={handleRegister}
      >
        <div className="container flex flex-col gap-3">
          <label className="self-start" htmlFor="email">
            Email
          </label>
          <input
            type="text"
            id="email"
            placeholder="Enter your email"
            className="rounded-xl border border-accent-secondary py-2 px-2 focus:outline-accent-secondary "
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
            className="rounded-xl border border-accent-secondary py-2 px-2 focus:outline-accent-secondary"
            value={password}
            required
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setPassword(e.target.value)
            }
          />
          <Button type="submit" text="Continue" />
        </div>
      </form>
      {error && (
        <span className="absolute bottom-10 error text-red-800">{error}</span>
      )}
    </>
  );
}
function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [token, setToken] = useState<string | null>(null);
  const navigate = useNavigate();
  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    try {
      const token = await loginUser(email, password);
      setToken(token);
      localStorage.setItem("authToken", token);
      navigate("/home");
    } catch (err) {
      setError((err as Error).message);
    }
  };
  if (token) {
    return <div>Login successful!</div>;
  }
  return (
    <>
      <form className="flex flex-col items-center pb-5" onSubmit={handleLogin}>
        <div className="container flex flex-col gap-3">
          <label className="self-start" htmlFor="email">
            Email
          </label>
          <input
            type="text"
            id="email"
            placeholder="Enter your email"
            className="rounded-xl border border-accent-secondary py-2 px-2 focus:outline-accent-secondary "
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
            className="rounded-xl border border-accent-secondary py-2 px-2 focus:outline-accent-secondary"
            value={password}
            minLength={8}
            required
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setPassword(e.target.value)
            }
          />
          <Button type="submit" text="Continue" />
        </div>
      </form>
      {error && (
        <span className="absolute bottom-10 error text-red-800">{error}</span>
      )}
    </>
  );
}
function AuthPage() {
  const [showRegister, setShowRegister] = useState(false);
  return (
    <main className="w-full h-screen flex flex-col items-center justify-center bg-base">
      <h1 className="text-3xl text-primary mb-5">red404</h1>
      <h3 className="text-lg mb-4">
        {showRegister ? "nice to meet u" : "welcome back"}
      </h3>
      <div className="auth-providers flex flex-col gap-5 mb-5 ">
        <GoogleAuthButton />
        <MetaAuthButton />
      </div>
      {showRegister ? <Register /> : <Login />}
      <button
        onClick={() => setShowRegister((prev) => !prev)}
        className=""
        type="button"
      >
        {showRegister
          ? "already registered? login"
          : "don't have an account? sign up"}
      </button>
    </main>
  );
}
export default AuthPage;
