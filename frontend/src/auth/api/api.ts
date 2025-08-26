import { AuthEndpoints } from "@/auth/api/endpoints";
import type { LoginType, RegisterType } from "@/auth/api/types";
import { authClient, type HttpClient } from "@/libs/authClient";

export const registerUser = async (
  data: RegisterType,
  httpClient: HttpClient = authClient,
) => {
  const response = await httpClient.post(AuthEndpoints.register, data);
  return response;
};

export const loginUser = async (
  data: LoginType,
  httpClient: HttpClient = authClient,
) => {
  const response = await httpClient.post(AuthEndpoints.login, data);
  return response;
};
