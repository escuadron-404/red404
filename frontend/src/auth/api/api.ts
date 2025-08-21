import { authClient, type HttpClient } from "../../libs/authClient";
import { AuthEndpoints } from "./endpoints";
import type { LoginType, RegisterType, loginResponseType } from "./types";

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
): Promise<loginResponseType> => {
  const response = await httpClient.post(AuthEndpoints.login, data);
  return response as loginResponseType;
};
