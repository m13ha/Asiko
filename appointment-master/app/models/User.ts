export interface User {
  id: string;
  name: string;
  email: string;
  phone: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface UserRequest {
  name: string;
  email: string;
  phoneNumber: string;
  password: string;
}