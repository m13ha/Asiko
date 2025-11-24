const TOKEN_KEY = 'token';
const REFRESH_TOKEN_KEY = 'refreshToken';
const EXPIRES_AT_KEY = 'tokenExpiresAt';

export type AuthTokens = {
  accessToken: string;
  refreshToken?: string;
  /**
   * Epoch milliseconds when the access token expires. Optional; will be derived from JWT exp if missing.
   */
  expiresAt?: number;
};

export function getToken() {
  return localStorage.getItem(TOKEN_KEY) || '';
}

export function getRefreshToken() {
  return localStorage.getItem(REFRESH_TOKEN_KEY) || '';
}

export function getTokens(): AuthTokens {
  const accessToken = getToken();
  const refreshToken = getRefreshToken();
  const expiresAtStr = localStorage.getItem(EXPIRES_AT_KEY);
  const expiresAt = expiresAtStr ? Number(expiresAtStr) : undefined;
  return { accessToken, refreshToken, expiresAt };
}

export function setTokens(tokens: AuthTokens) {
  localStorage.setItem(TOKEN_KEY, tokens.accessToken);
  if (tokens.refreshToken) {
    localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refreshToken);
  }
  if (tokens.expiresAt) {
    localStorage.setItem(EXPIRES_AT_KEY, tokens.expiresAt.toString());
  } else {
    localStorage.removeItem(EXPIRES_AT_KEY);
  }
}

export function setToken(token: string) {
  setTokens({ accessToken: token });
}

export function clearToken() {
  clearTokens();
}

export function clearTokens() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
  localStorage.removeItem(EXPIRES_AT_KEY);
}

export function isAccessTokenFresh(thresholdMs = 60_000) {
  const { accessToken, expiresAt } = getTokens();
  if (!accessToken) return false;
  if (!expiresAt) return true;
  return expiresAt - Date.now() > thresholdMs;
}
