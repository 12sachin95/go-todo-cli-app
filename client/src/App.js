import React from "react";
import { Route, Routes } from "react-router-dom";
import Login from "./pages/Login";
import Register from "./pages/Register";
import TodoList from "./pages/TodoList";
import Header from "./components/Header";
import "./App.css";

function App() {
  return (
    <>
      <div className="app-container">
        <nav>
          <Header />
        </nav>
        <Routes>
          <Route path="/" element={<>Todo APP</>} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/todos" element={<TodoList />} />
        </Routes>
      </div>
    </>
  );
}

export default App;
