import Peer from 'simple-peer';
import wasm from './main.go'

const peer = new Peer({
  initiator: !!!window.location.hash,
  trickle: false
});

let secret: string;
let publicKey: string;

peer.on('error', err => console.log('error', err))

peer.on('signal', data => {
  encrypt(JSON.stringify(data));
  document.querySelector('#outgoing')!.textContent = JSON.stringify(data)
})

document.querySelector('form')!.addEventListener('submit', ev => {
  ev.preventDefault()
  peer.signal(JSON.parse((document.querySelector('#incoming') as HTMLInputElement)!.value!))
})

peer.on('connect', async () => {
  peer.send(`keyData:${publicKey}`);
})

peer.on('data', async (data) => {
  const string = new TextDecoder("utf-8").decode(data);

  if (string.startsWith("keyData:")) {
    secret = await wasm.ComputeSecret(string.replace("keyData:", ""));
    console.log('secret', secret)
  }

  console.log('data', string);
})

async function encrypt(signalData: string) {
  const uploadKey = await wasm.GenerateRandomString(24);
  const buffer = new Buffer(signalData);

  publicKey = await wasm.GenerateKey();

  const encryptedSignalData = await wasm.EncryptFile('signal', buffer, buffer.length, uploadKey);
  //const response = await wasm.UploadFile(data);

  console.log('uploadKey', uploadKey)
  console.log('uploadData', signalData);
  console.log('publicKey', publicKey)
}
