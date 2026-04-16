import {components} from "@/generated/types";

export type PresignedURLRequest = components["schemas"]["PresignedURLRequest"]
export type PresignedURLResponse = components["schemas"]["PresignedURLResponse"]

// This is used to directly upload to minio/s3 and is not intended for use with the backend api
// therefore it's not defined in the openapi spec and generated types
export type FileUploadRequest = {
    url: string;
    file: File;
}