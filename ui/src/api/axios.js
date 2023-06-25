import axios from "axios";
import { SERVER_ENDPOINT } from "../util/env.js";
const BASE_URL = SERVER_ENDPOINT;

export default axios.create({
  baseURL: BASE_URL,
});

export const axiosPrivate = axios.create({
  baseURL: BASE_URL,
  headers: { "Content-Type": "application/json" },
  withCredentials: true,
});
