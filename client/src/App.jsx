import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, useNavigate } from "react-router-dom";
import "./App.css";
import { Container } from "semantic-ui-react";
import ToDoList from "./To-do-lists";
import Login from "./Authentication/Login";
import Register from "./Authentication/Register";
import GetTogether from "./Pages/GetTogether";
import HeaderComp from "./Components/Header";
import LeftBar from "./Components/LeftBar";
import RightBar from "./Components/RightBar";
import Aurora from "./Aurora"; // Import Aurora component

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
      <div className="app-container">
        <Aurora colorStops={["#3a3a66", "#7a7aed", "#9a48ad"]} speed={1} />
        <HeaderComp />
        <LeftBar setFilter={setFilter} />
        <RightBar isAuthenticated={isAuthenticated} setIsAuthenticated={setIsAuthenticated} />
        <AppContent filter={filter} isAuthenticated={isAuthenticated} setIsAuthenticated={setIsAuthenticated} />
      </div>
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
    <Container style={{ marginLeft: '220px', marginRight: '220px' }}>
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