import { Navigate } from "react-router";
import { UseAuth } from "../auth/context/auth-context";

interface ProtectedRouteProps {
  children: React.ReactNode;
}

export default function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, isInitialized } = UseAuth();

  // Mostrar loading mientras se inicializa
  if (!isInitialized || isLoading) {
    return (
      <div className="flex min-h-svh items-center justify-center">
        <div className="text-muted-foreground">Loading...</div>
      </div>
    );
  }

  // Redirigir a login si no está autenticado
  if (!isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  // Renderizar el contenido protegido
  return <>{children}</>;
}
