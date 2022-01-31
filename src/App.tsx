import React, { useEffect, useState } from 'react';
import Peer from 'simple-peer';
import wasm from './go/main.go'
import './App.css';

let secret: string;
let publicKey: string;

const peer = new Peer({
  initiator: !!!window.location.hash,
  trickle: false
});

async function encrypt(signalData: string) {
  const uploadKey = await wasm.GenerateRandomString(24);
  const buffer = new Buffer(signalData);

  publicKey = await wasm.GenerateKey();

  const encryptedSignalData = await wasm.EncryptFile('signal', buffer, buffer.length, uploadKey);
  //const response = await wasm.UploadFile(data);

  console.log('uploadKey', uploadKey)
  console.log('uploadData', encryptedSignalData);
  console.log('publicKey', publicKey)
}

function App() {
  const [incomingText, setIncomingtext] = useState("")
  const [outgoingText, setOutgoingtext] = useState("")

  useEffect(() => {
    peer.on('error', err => console.log('error', err))

    peer.on('signal', data => {
      encrypt(JSON.stringify(data));
      setOutgoingtext(JSON.stringify(data));
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
  }, []);

  const handleSubmit = (event: React.SyntheticEvent) => {
    event.preventDefault();
    peer.signal(JSON.parse(incomingText))
  }

  return (
    <div className="container">
      <div className="body">
        <div className="App">
          <div className="row logo">
            <div className="col-md-12 text-center">
              <h1 className="h1"><a href="/">chat<span>sh.it</span></a></h1>
            </div>
          </div>
          <div className="text-center col-md-12">
            <p></p>
            <form onSubmit={ handleSubmit }>
              <textarea
                id="incoming"
                value = {incomingText}
                onChange={e => setIncomingtext(e.target.value)}>
              </textarea>
              <button type="submit" className="btn btn-default">Go</button>
            </form>
          </div>

          <pre id="outgoing">{ outgoingText }</pre>
        </div>
      </div>
    </div>
  );
}

export default App;
