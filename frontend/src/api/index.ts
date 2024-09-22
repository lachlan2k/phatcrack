import axios, { type AxiosInstance } from "axios";

export let client = axios.create();

export function setClient(newClient: AxiosInstance) {
  client = newClient;
}

export function ping(): Promise<string> {
  return client.get("/api/v1/ping").then((res) => res.data);
}

export * from "./account";
export * from "./admin";
export * from "./agent";
export * from "./attackTemplate";
export * from "./auth";
export * from "./config";
export * from "./hashcat";
export * from "./listfiles";
export * from "./potfile";
export * from "./project";
export * from "./types";
export * from "./users";
