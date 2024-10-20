import React, { useState, useContext } from "react";
import axios from "axios";
import AuthContext from "../context/AuthContext";
import { useNavigate } from "react-router-dom";
import "./Auth.css";

const Login = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const { login } = useContext(AuthContext);

  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    setError("");
    e.preventDefault();
    try {
      const response = await axios.post(
        process.env.REACT_APP_TODO_SERVER_PATH + "/user/login",
        {
          username,
          password,
        }
      );
      login(response.data.token);
      navigate("/todos");
    } catch (error) {
      setError(error.response.data?.error ?? "Something went wrong");
      console.error("Login failed:", error.response.data);
    }
  };

  return (
    <div className="auth-container">
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        {error ? <div style={{ color: "red" }}>{error}</div> : <></>}
        <input
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder="Username"
          required
        />
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Password"
          required
        />
        <button type="submit">Login</button>
      </form>
    </div>
  );
};

export default Login;
