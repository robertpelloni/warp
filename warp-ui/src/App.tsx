import { useState, useEffect } from 'react';
import reactLogo from './assets/react.svg';
import viteLogo from './assets/vite.svg';
import './App.css';

interface ServerStatus {
  status: string;
  version: string;
  message: string;
}

function App() {
  const [serverStatus, setServerStatus] = useState<ServerStatus | null>(null);

  useEffect(() => {
    fetch('http://localhost:8080/status')
      .then((res) => res.json())
      .then((data) => setServerStatus(data))
      .catch((err) => console.error("Error fetching server status:", err));
  }, []);

  return (
    <>
      <div>
        <a href="https://vite.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Warp UI</h1>
      <div className="card">
        {serverStatus ? (
          <div>
            <p><strong>Status:</strong> {serverStatus.status}</p>
            <p><strong>Version:</strong> {serverStatus.version}</p>
            <p><strong>Message:</strong> {serverStatus.message}</p>
          </div>
        ) : (
          <p>Connecting to backend server...</p>
        )}
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  )
}

export default App;
