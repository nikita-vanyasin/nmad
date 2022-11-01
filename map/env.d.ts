/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_NMAD_API_BASE_URL: string
  readonly VITE_GOOGLE_API_KEY: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}