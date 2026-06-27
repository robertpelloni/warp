import { useState, useEffect, useRef } from 'react';
import reactLogo from './assets/react.svg';
import './App.css';

interface ServerStatus {
  status: string;
  version: string;
  message: string;
}

function App() {
  const [serverStatus, setServerStatus] = useState<ServerStatus | null>(null);
  const [messages, setMessages] = useState<string[]>([]);
  const [input, setInput] = useState('');
  const ws = useRef<WebSocket | null>(null);

  useEffect(() => {
    fetch('http://localhost:8080/status')
      .then((res) => res.json())
      .then((data) => setServerStatus(data))
      .catch((err) => console.error("Error fetching server status:", err));

    ws.current = new WebSocket('ws://localhost:8080/ws');

    ws.current.onopen = () => {
      setMessages(prev => [...prev, "[System] WebSocket connected to Warp Orchestrator"]);
    };

    ws.current.onmessage = (event) => {
      setMessages(prev => [...prev, event.data]);
    };

    ws.current.onclose = () => {
      setMessages(prev => [...prev, "[System] WebSocket disconnected"]);
    };

    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, []);

  const sendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (input.trim() && ws.current?.readyState === WebSocket.OPEN) {
      setMessages(prev => [...prev, `[User] ${input}`]);
      ws.current.send(input);
      setInput('');
    }
  };

  return (
    <>
      <div>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Warp UI (Alpha)</h1>
      <div className="card">
        {serverStatus ? (
          <div style={{ marginBottom: "2rem", borderBottom: "1px solid #ccc", paddingBottom: "1rem" }}>
            <p><strong>Status:</strong> {serverStatus.status}</p>
            <p><strong>Version:</strong> {serverStatus.version}</p>
            <p><strong>Message:</strong> {serverStatus.message}</p>
          </div>
        ) : (
          <p>Connecting to backend server...</p>
        )}

        <div style={{
          height: "300px",
          overflowY: "auto",
          background: "#1a1a1a",
          padding: "1rem",
          borderRadius: "8px",
          textAlign: "left",
          fontFamily: "monospace"
        }}>
          {messages.map((msg, i) => (
            <div key={i} style={{ marginBottom: "0.5rem" }}>{msg}</div>
          ))}
        </div>

        <form onSubmit={sendMessage} style={{ marginTop: "1rem", display: "flex", gap: "10px" }}>
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Type a command (/prompt, /browser...)"
            style={{ flexGrow: 1, padding: "0.5rem", borderRadius: "4px" }}
          />
          <button type="submit" disabled={!ws.current || ws.current.readyState !== WebSocket.OPEN}>
            Send
          </button>
        </form>
      </div>
    </>
  )
}

export default App;
