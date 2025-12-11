export interface LoginRequest {
  username: string;
  password: string;
  remember?: boolean;
}

export interface RegisterRequest {
  username: string;
  password: string;
}

export interface AuthTokens {
  accessToken: string;
  expiresIn: number;
}

export interface Profile {
  id: number;
  username: string;
}
