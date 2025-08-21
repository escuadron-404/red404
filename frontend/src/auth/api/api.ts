import { authClient, type HttpClient } from "../../libs/authClient";
import { AuthEndpoints } from "./endpoints";
import type { LoginType, RegisterType, loginResponseType, registerResponseType } from "./types";

export const registerUser = async (
  data: RegisterType,
  httpClient: HttpClient = authClient,
): Promise<registerResponseType> => {
  const response: registerResponseType = await httpClient.post(
    AuthEndpoints.register,
    data,
  );
  return response;
};

export const loginUser = async (
  data: LoginType,
  httpClient: HttpClient = authClient,
): Promise<loginResponseType> => {
  const response: loginResponseType = await httpClient.post(
    AuthEndpoints.login,
    data,
  );
  return response;
};
