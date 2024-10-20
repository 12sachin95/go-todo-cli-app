import React, { useState, useEffect, useContext } from "react";
import axios from "axios";
import AuthContext from "../context/AuthContext";
import { Navigate } from "react-router-dom";
import "./TodoList.css";

const TodoList = () => {
  const [todos, setTodos] = useState([]);
  const { token } = useContext(AuthContext);
  const [newTodo, setNewTodo] = useState({ title: "", completed: false });
  const [isEditing, setIsEditing] = useState(null);

  const [updatedTodo, setUpdatedTodo] = useState({
    title: "",
    completed: false,
  });

  useEffect(() => {
    const fetchTodos = async () => {
      try {
        const response = await axios.get(
          process.env.REACT_APP_TODO_SERVER_PATH + "/todos",
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        setTodos(response.data);
      } catch (error) {
        console.error("Error fetching todos:", error);
      }
    };

    if (token) {
      fetchTodos();
    }
  }, [token]);

  const handleAddTodo = async (e) => {
    e.preventDefault();
    try {
      await axios.post(
        process.env.REACT_APP_TODO_SERVER_PATH + "/todos",
        newTodo,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      setTodos([newTodo, ...(todos ?? [])]);
      setNewTodo({ title: "", completed: false });
    } catch (error) {
      console.error("Error adding todo:", error);
    }
  };

  const handleDeleteTodo = async (id) => {
    try {
      await axios.delete(
        process.env.REACT_APP_TODO_SERVER_PATH + `/todos/${id}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      setTodos(todos.filter((todo) => todo.id !== id));
    } catch (error) {
      console.error("Error deleting todo:", error);
    }
  };

  const startEditing = (todo) => {
    setIsEditing(todo.id);
    setUpdatedTodo({ title: todo.title, completed: todo.completed });
  };

  const handleUpdateTodo = async (id) => {
    try {
      await axios.put(
        process.env.REACT_APP_TODO_SERVER_PATH + `/todos/${id}`,
        updatedTodo,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      setTodos((previews) =>
        previews.map((todo) =>
          todo.id === id ? { ...todo, updatedTodo } : todo
        )
      );
      setIsEditing(null);
    } catch (error) {
      console.error("Error updating todo:", error);
    }
  };

  if (!token) {
    return <Navigate to="/login" />;
  }

  return (
    <div className="todo-container">
      <h2>Your Todos</h2>
      <form onSubmit={handleAddTodo}>
        <input
          type="text"
          value={newTodo.title}
          onChange={(e) =>
            setNewTodo({ title: e.target.value, completed: false })
          }
          placeholder="Add a new todo"
          required
        />
        <button type="submit">Add Todo</button>
      </form>
      <ul>
        {todos?.map((todo) => (
          <li key={todo.id}>
            {isEditing === todo.id ? (
              <div
                style={{
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "space-between",
                  width: "100%",
                  gap: "8px",
                }}
              >
                <input
                  type="text"
                  value={updatedTodo.title}
                  onChange={(e) =>
                    setUpdatedTodo((previews) => ({
                      ...previews,
                      title: e.target.value,
                    }))
                  }
                />
                <div class="todo-item">
                  <label class="custom-checkbox">
                    <input
                      type="checkbox"
                      id="todo-completed"
                      checked={updatedTodo.completed}
                      onChange={(e) => {
                        console.log(e);
                        setUpdatedTodo((previews) => ({
                          ...previews,
                          completed: e.target.checked,
                        }));
                      }}
                    />
                    <span class="checkmark"></span>
                  </label>
                </div>
                <button
                  onClick={() => handleUpdateTodo(todo.id)}
                  className="edit-button"
                >
                  Save
                </button>
                <button onClick={() => setIsEditing(null)}>Cancel</button>
              </div>
            ) : (
              <div
                style={{
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "space-between",
                  width: "100%",
                }}
              >
                <span
                  style={{
                    textDecoration: todo.completed ? "line-through" : "none",
                  }}
                >
                  {todo.title}{" "}
                </span>
                <div
                  style={{
                    display: "flex",
                    gap: "8px",
                  }}
                >
                  <button
                    className="edit-button"
                    onClick={() => startEditing(todo)}
                  >
                    Edit
                  </button>
                  <button onClick={() => handleDeleteTodo(todo.id)}>
                    Delete
                  </button>
                </div>
              </div>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default TodoList;
