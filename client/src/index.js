import React from "react";
import ReactDOM from "react-dom";
import { AuthProvider } from "./context/AuthContext";
import App from "./App";
import { BrowserRouter } from "react-router-dom";
import "./index.css";

ReactDOM.render(
  <AuthProvider>
    <BrowserRouter basename="/todo-app">
      <App />
    </BrowserRouter>
  </AuthProvider>,
  document.getElementById("root")
);
