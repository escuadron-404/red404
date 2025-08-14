import { authClient, type HttpClient } from "../../libs/authClient";
import type { LoginType, RegisterType } from "./types";
import { AuthEndpoints } from "./endpoints";

export const registerUser = async (data: RegisterType, httpClient: HttpClient = authClient) => {
    const response = await httpClient.post(AuthEndpoints.register, data);
    return response;
};

export const loginUser = async (data: LoginType, httpClient: HttpClient = authClient) => {
    const response = await httpClient.post(AuthEndpoints.login, data);
    return response;
};