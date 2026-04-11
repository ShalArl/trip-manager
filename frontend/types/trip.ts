import { components } from "@/generated/types";


export type CreateTripRequest = components["schemas"]["CreateTripRequest"];
export type TripResponse = components["schemas"]["TripResponse"];


export type Trip = {
    title: string;
    startDate: string;
    endDate: string;
    shortDescription: string;
    description: string;
}