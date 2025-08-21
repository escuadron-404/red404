export type LoginType = {
  email: string;
  password: string;
};

// this might change in the future
export type RegisterType = {
  email: string;
  password: string;
  [key: string]: string;
};

export type loginResponseType = {
  success: boolean;
  message: string;
  data: {
    token: string;
    user: {
      id: string;
      email: string;
    };
  };
};

export type ResponseType = Record<string, unknown>;
