/// <reference types="vite/client" />

declare module 'lodash.random' {
  const random: (min: number, max: number) => number
  export default random
}

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string
  readonly VITE_API_TIMEOUT_MS: string
  readonly VITE_IN_MAINTENANCE_MODE: string
  readonly VITE_IN_MAINTENANCE_MODE_CODE: string
  readonly BASE_URL: string
  readonly DEV: boolean
  readonly PROD: boolean
}

declare module 'crypto-js/aes' {
  const AES: {
    encrypt(message: string, key: string): { toString(): string }
    decrypt(
      ciphertext: string,
      key: string,
    ): {
      toString(enc: { parse(str: string): any; stringify(wordArray: any): string }): string
    }
  }
  export default AES
}

declare module 'crypto-js/enc-utf8' {
  const Utf8: { parse(str: string): any; stringify(wordArray: any): string }
  export default Utf8
}

declare module 'crypto-js' {
  const CryptoJS: {
    AES: {
      encrypt(message: string, key: string): { toString(): string }
      decrypt(ciphertext: string, key: string): { toString(enc: typeof CryptoJS.enc.Utf8): string }
    }
    enc: {
      Utf8: { parse(str: string): any; stringify(wordArray: any): string }
    }
  }
  export default CryptoJS
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
