import { components } from "@/generated/types";

export type TransportResponse = components["schemas"]["TransportResponse"];
export type CreateTransportRequest = components["schemas"]["CreateTransportRequest"];
export type UpdateTransportRequest = components["schemas"]["UpdateTransportRequest"];

export type Transport = {
    fromLocationId: string;
    toLocationId: string;
    departureTime?: string;
    arrivalTime?: string;
    type: "flight" | "train" | "car" | "bus";
    notes?: string;
}