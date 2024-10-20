import React, { useContext } from "react";
import AuthContext from "../context/AuthContext";
import axios from "axios";
import { useNavigate } from "react-router-dom";

const Header = () => {
  const { logout, token } = useContext(AuthContext);
  const navigate = useNavigate();

  const handleLogout = async (e) => {
    e.preventDefault();
    try {
      await axios.post(
        process.env.REACT_APP_TODO_SERVER_PATH + "/user/logout",
        {},
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      logout();
      navigate("/login");
    } catch (error) {
      console.error("Logout failed:", error.response.data);
    }
  };

  return (
    <header>
      <h1>Todo App</h1>
      <nav>
        {token ? (
          <button className="logout-button" onClick={handleLogout}>
            Logout
          </button>
        ) : (
          <></>
        )}
      </nav>
    </header>
  );
};

export default Header;
