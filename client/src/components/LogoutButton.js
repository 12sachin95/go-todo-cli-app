import React, { useContext } from "react";
import AuthContext from "../context/AuthContext";

const LogoutButton = () => {
  const { logout } = useContext(AuthContext);

  return (
    <button className="logout-btn" onClick={logout}>
      Logout
    </button>
  );
};

export default LogoutButton;
