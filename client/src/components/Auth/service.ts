const AUTH_PATH = "/auth";

type LoginProps = {
  name: string;
  email?: string;
  channel: string;
}

export type UserType = {
  id: string;
  name: string;
  avatar: string;
}

export type AuthType = {
  user: UserType;
  channel: string;
}

export class AuthService {
  static get token(): string | null {
    return localStorage.getItem('token');
  }

  static set token(value: string | null) {
    if (value) {
      localStorage.setItem('token', value);
    } else {
      localStorage.removeItem('token');
    }

  }

  static serverAuthUrl() {
    return process.env.REACT_APP_SERVER_URL + AUTH_PATH;
  }

  static async auth(): Promise<AuthType | null> {
    const token = this.token;

    if (!token) {
      return null;
    }

    const url = new URL(this.serverAuthUrl());

    url.searchParams.append("token", token);

    return fetch(url).then((response) => response.ok ? response.json() : null);
  }

  token(): string | null {
    return localStorage.getItem('token');
  }

  static async login(auth: LoginProps): Promise<AuthType | null> {
    const body = new FormData();

    for (let [key, value] of Object.entries(auth)) {
      body.append(key, value);
    }

    const token = await fetch(this.serverAuthUrl(), {
      method: "POST",
      body
    }).then((response) => {
      if (!response.ok) {
        throw new Error(response.statusText);
      }

      return response.json();
    }).then((body) => body['token']);

    this.token = token;

    return this.auth();
  }

  static reset() {
    this.token = null;
  }
}

