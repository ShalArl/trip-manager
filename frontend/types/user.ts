import {components} from "@/generated/types";


export type UserResponse = components["schemas"]["UserResponse"];
export type UpdateUserRequest = components["schemas"]["UpdateUserRequest"];
export type ProvisionUserRequest = components["schemas"]["ProvisionUserRequest"]

// types/user.ts (oder types/middleware.ts, wo auch immer es passt)

export interface CreateUserRequest {
    email: string;
    password: string;
    name: string;
    tenantName?: string;
    tier?: string;
}

export interface LoginRequest {
    email: string;
    password: string;
}

export interface ChangePasswordRequest {
    currentPassword: string;
    newPassword: string;
}