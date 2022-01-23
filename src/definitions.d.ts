declare module "*.go" {
  export declare const GenerateRandomString: (size: number) => string;
  export declare const EncryptFile: (name: string, bytes: Uint8Array, length: number, key: string) => string;
  export declare const UploadFile: (data: string) => string;
  export declare const GenerateKey: () => string;
  export declare const ComputeSecret: (public: string) => string;
}
