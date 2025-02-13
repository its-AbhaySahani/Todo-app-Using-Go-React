import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, Link, useNavigate } from "react-router-dom";
import "./App.css";
import { Container } from "semantic-ui-react";
import ToDoList from "./To-do-lists";
import Login from "./Authentication/Login";
import Register from "./Authentication/Register";
import GetTogether from "./Pages/GetTogether";
import HeaderComp from "./Components/Header"; // Import HeaderComp component
import LeftBar from "./Components/LeftBar"; // Import LeftBar component
import RightBar from "./Components/RightBar"; // Import RightBar component

function App() {
  const [filter, setFilter] = useState('all');
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      setIsAuthenticated(true);
    }
  }, []);

  return (
    <Router>
      <HeaderComp />
      <LeftBar setFilter={setFilter} /> {/* Add LeftBar component */}
      <RightBar isAuthenticated={isAuthenticated} setIsAuthenticated={setIsAuthenticated} /> {/* Add RightBar component */}
      <AppContent filter={filter} isAuthenticated={isAuthenticated} setIsAuthenticated={setIsAuthenticated} />
    </Router>
  );
}

function AppContent({ filter, isAuthenticated, setIsAuthenticated }) {
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem("token");
    setIsAuthenticated(false);
    navigate("/login");
  };

  return (
    <Container style={{ marginLeft: '220px', marginRight: '220px' }}> {/* Add margin to account for the sidebars */}
      <Routes>
        <Route path="/login" element={<Login setIsAuthenticated={setIsAuthenticated} />} />
        <Route path="/register" element={<Register />} />
        <Route path="/" element={<ToDoList filter={filter} />} />
        <Route path="/get-together" element={<GetTogether />} />
      </Routes>
    </Container>
  );
}

export default App;