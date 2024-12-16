import React, { useState } from "react";
import axios from "axios";
import "./App.css";

function App() {
  const [file, setFile] = useState(null);
  const [query, setQuery] = useState("");
  const [question, setQuestion] = useState("");
  const [response, setResponse] = useState("");
  const [activeMenu, setActiveMenu] = useState("chatAI");
  const [isLoading, setIsLoading] = useState(false);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleQuestionChange = (e) => {
    setQuestion(e.target.value);
  };

  const handleUpload = async () => {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("question", question);

    setIsLoading(true);
    try {
      const res = await axios.post("http://localhost:8080/upload", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
      setResponse(res.data.answer);
    } catch (error) {
      console.error("Error uploading file:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleChat = async () => {
    setIsLoading(true);
    try {
      const res = await axios.post("http://localhost:8080/chat", { query });
      setResponse(res.data.answer);
    } catch (error) {
      console.error("Error querying chat:", error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-between",
        height: "100vh",
        fontFamily: "Arial, sans-serif",
        backgroundColor: "#1a2238",
        color: "#e0e6ed",
      }}
    >
      <div
        style={{
          flex: 1,
          overflowY: "auto",
          padding: "20px",
          boxSizing: "border-box",
        }}
      >
        <h1
          style={{
            color: "#e0e6ed",
            fontSize: "1.8rem",
            marginBottom: "20px",
            textAlign: "center",
          }}
        >
          TAKON AI
        </h1>
        <div className="button-menu" style={{ marginBottom: "10px" }}>
          <button className="button" onClick={() => setActiveMenu("chatAI")}>
            Chat AI
          </button>
          <button className="button" onClick={() => setActiveMenu("uploadAnalyze")}>
            Upload and Analyze
          </button>
        </div>
        <div
          style={{
            padding: "20px",
            border: "1px solid #2c3e50",
            borderRadius: "8px",
            backgroundColor: "#273c75",
          }}
        >
          <h2 style={{ fontSize: "1.2rem", marginBottom: "10px" }}>Response</h2>
          {isLoading ? (
            <div
              style={{
                color: "#2e86de",
                fontStyle: "italic",
                fontSize: "1rem",
              }}
            >
              Typing<span className="dots">...</span>
            </div>
          ) : (
            <p style={{ whiteSpace: "pre-wrap", wordWrap: "break-word" }}>{response}</p>
          )}
        </div>
      </div>

      <div
        style={{
          borderTop: "1px solid #2c3e50",
          padding: "10px 20px",
          backgroundColor: "#1a2238",
          boxSizing: "border-box",
        }}
      >
        {activeMenu === "chatAI" && (
          <div style={{ position: "relative", display: "flex", alignItems: "center" }}>
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Ask a question..."
              style={{
                flex: 1,
                padding: "10px 40px 10px 10px",
                border: "1px solid #2c3e50",
                borderRadius: "6px",
                backgroundColor: "#2c3e50",
                color: "#e0e6ed",
              }}
            />
            <button
              onClick={handleChat}
              style={{
                position: "absolute",
                right: "10px",
                top: "50%",
                transform: "translateY(-50%)",
                backgroundColor: "transparent",
                border: "none",
                color: "#2e86de",
                fontSize: "1.2rem",
                cursor: "pointer",
              }}
            >
              ➤
            </button>
          </div>
        )}

        {activeMenu === "uploadAnalyze" && (
          <div style={{ position: "relative", display: "flex", flexDirection: "column", gap: "10px" }}>
            <input
              type="text"
              value={question}
              onChange={handleQuestionChange}
              placeholder="Question related to file"
              style={{
                padding: "10px 40px 10px 40px",
                border: "1px solid #2c3e50",
                borderRadius: "6px",
                backgroundColor: "#2c3e50",
                color: "#e0e6ed",
                position: "relative",
              }}
            />
            <label
              htmlFor="file-upload"
              style={{
                position: "absolute",
                left: "10px",
                top: "50%",
                transform: "translateY(-50%)",
                cursor: "pointer",
                color: "#e0e6ed",
              }}
            >
              📎
            </label>
            <input
              id="file-upload"
              type="file"
              onChange={handleFileChange}
              style={{
                display: "none",
              }}
            />
            <button
              onClick={handleUpload}
              style={{
                position: "absolute",
                right: "10px",
                top: "50%",
                transform: "translateY(-50%)",
                backgroundColor: "transparent",
                border: "none",
                color: "#2e86de",
                fontSize: "1.2rem",
                cursor: "pointer",
              }}
            >
              ➤
            </button>
          </div>
        )}
      </div>
      <footer
        style={{
          padding: "10px 20px",
          backgroundColor: "#1a2238",
          textAlign: "center",
          color: "#e0e6ed",
          fontSize: "0.9rem",
        }}
      >
        <span>
          TAKON AI | Powered by <strong>ADTY • Ruang Guru</strong>
        </span>
      </footer>
    </div>
  );
}

export default App;
