import {
  createContext,
  type ReactNode,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { loginUser } from "../api/api";
import type { loginResponseType } from "../api/types";

// Tipos para el usuario
export interface User {
  id: string;
  email: string;
}

// Contexto
interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  isInitialized: boolean;
  error: string | null;
  login: (
    email: string,
    password: string,
  ) => Promise<loginResponseType | undefined>;
  logout: () => void;
  clearError: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Provider
interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [user, setUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isInitialized, setIsInitialized] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Función de login centralizada con useCallback para evitar re-renderizados
  const login = useCallback(async (email: string, password: string) => {
    try {
      setIsLoading(true);
      setError(null);

      const response: loginResponseType = await loginUser({ email, password });

      // Guardar token y datos del usuario
      if (response.data.token) {
        localStorage.setItem("auth_token", response.data.token);
      }

      if (response.data.user) {
        localStorage.setItem("user_data", JSON.stringify(response.data.user));
      }

      // Establecer usuario
      setUser(response.data.user);
      setIsAuthenticated(true);

      return response;
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Error de autenticación";
      setError(errorMessage);
      throw error; // Re-lanzar para que el componente pueda manejarlo
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Función de logout centralizada con useCallback
  const logout = useCallback(() => {
    localStorage.removeItem("auth_token");
    localStorage.removeItem("user_data");
    setUser(null);
    setError(null);
    setIsAuthenticated(false);
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // Verificar estado de autenticación mejorado
  const checkAuthStatus = useCallback(async () => {
    const token = localStorage.getItem("auth_token");

    if (!token) {
      setIsInitialized(true);
      return;
    }

    try {
      setIsLoading(true);

      // Por ahora, si hay token asumimos que el usuario está autenticado
      // TODO: Implementar endpoint para verificar token válido
      // En una implementación real, harías algo como:
      // const response = await authClient.get('/auth/verify');
      // setUser(response.user);

      // Simulación temporal - puedes implementar la verificación real del token aquí
      const storedUser = localStorage.getItem("user_data");
      if (storedUser) {
        try {
          const userData = JSON.parse(storedUser);
          setUser(userData);
          setIsAuthenticated(true);
        } catch {
          // Si no se puede parsear, limpiar localStorage
          localStorage.removeItem("auth_token");
          localStorage.removeItem("user_data");
        }
      }
    } catch (error) {
      console.log("auth error: ", error);
      localStorage.removeItem("auth_token");
      localStorage.removeItem("user_data");
      setUser(null);
      setIsAuthenticated(false);
    } finally {
      setIsLoading(false);
      setIsInitialized(true);
    }
  }, []);

  // Verificar estado de autenticación al cargar la app (solo una vez)
  useEffect(() => {
    checkAuthStatus();
  }, [checkAuthStatus]);

  // Memoizar el valor del contexto para evitar re-renderizados innecesarios
  const value = useMemo(
    (): AuthContextType => ({
      user,
      isAuthenticated,
      isLoading,
      isInitialized,
      error,
      login,
      logout,
      clearError,
    }),
    [
      user,
      isAuthenticated,
      isLoading,
      isInitialized,
      error,
      login,
      logout,
      clearError,
    ],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const UseAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
