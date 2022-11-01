
import axios from 'axios';

export const Api = axios.create({
  baseURL: `${import.meta.env.VITE_NMAD_API_BASE_URL}/api/v1`
})
