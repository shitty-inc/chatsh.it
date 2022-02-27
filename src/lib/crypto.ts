import wasm from './main.go';

export async function encrypt(data: string, secret: string): Promise<string> {
  const buffer = new Buffer(data);
  const encryptedData = await wasm.Encrypt(buffer, buffer.length, secret);

  return Buffer.from(encryptedData).toString('hex');
}

export async function decrypt(data: string, secret: string): Promise<string> {
  const buffer = Buffer.from(data, 'hex');
  const decryptedData = await wasm.Decrypt(buffer, buffer.length, secret);

  return Buffer.from(decryptedData).toString();
}

export async function ComputeSecret(key: string): Promise<string> {
  return await wasm.ComputeSecret(key);
}

export async function GenerateRandomString(size: number): Promise<string> {
  return await wasm.GenerateRandomString(size);
}

export async function GenerateKeyPair(): Promise<string> {
  return await wasm.GenerateKeyPair();
}
