import { components } from "@/generated/types";

export type LocationResponse = components["schemas"]["LocationResponse"];
export type CreateLocationRequest = components["schemas"]["CreateLocationRequest"];
export type UpdateLocationRequest = components["schemas"]["UpdateLocationRequest"];

export type Location = {
    name: string;
    city: string;
    country: string;
    countryCode: string;
    shortDescription: string;
    dateFrom: string;
    dateTo: string;
    latitude?: number;
    longitude?: number;
    notes?: string;
}