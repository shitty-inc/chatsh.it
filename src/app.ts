import Peer from 'simple-peer';
import wasm from './main.go'

const peer = new Peer({
  initiator: !!!window.location.hash,
  trickle: false
});

peer.on('error', err => console.log('error', err))

peer.on('signal', data => {
  encrypt(JSON.stringify(data));
})

peer.on('connect', () => {
  console.log('CONNECT')
  peer.send('whatever' + Math.random())
})

peer.on('data', data => {
  console.log('DATA', data)
})

async function encrypt(json: string) {
  const buffer = new Buffer(json);
  const key = await wasm.GenerateRandomString(24);
  const data = await wasm.EncryptFile('signal', buffer, buffer.length, key);
  const response = await wasm.UploadFile(data);

  console.log('KEY', key)
  console.log('DONE', response);
}
