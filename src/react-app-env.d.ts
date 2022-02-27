/// <reference types="react-scripts" />

declare module "*.go" {
  export declare const GenerateRandomString: (size: number) => string;
  export declare const Encrypt: (bytes: Uint8Array, length: number, key: string) => Uint8Array;
  export declare const Decrypt: (bytes: Uint8Array, length: number, key: string) => string;
  export declare const GenerateKeyPair: () => string;
  export declare const ComputeSecret: (public: string) => string;
}
