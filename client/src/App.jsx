import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, Link, useNavigate } from "react-router-dom";
import "./App.css";
import { Container, Menu, Icon } from "semantic-ui-react";
import ToDoList from "./To-do-lists";
import Login from "./Authentication/Login";
import Register from "./Authentication/Register";
import GetTogether from "./Pages/GetTogether";
import HeaderComp from "./Components/Header"; // Import HeaderComp component
import LeftBar from "./Components/LeftBar"; // Import LeftBar component

function App() {
  const [filter, setFilter] = useState('all');

  return (
    <Router>
      <HeaderComp />
      <LeftBar setFilter={setFilter} /> {/* Add LeftBar component */}
      <AppContent filter={filter} />
    </Router>
  );
}

function AppContent({ filter }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      setIsAuthenticated(true);
    }
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("token");
    setIsAuthenticated(false);
    navigate("/login");
  };

  return (
    <Container style={{ marginLeft: '220px' }}> {/* Add margin to account for the sidebar */}
      <Menu>
        <Menu.Item as={Link} to="/">
          Home
        </Menu.Item>
        {isAuthenticated ? (
          <>
            <Menu.Item position="right">
              <Icon name="user" />
            </Menu.Item>
            <Menu.Item onClick={handleLogout}>
              Logout
            </Menu.Item>
          </>
        ) : (
          <>
            <Menu.Item as={Link} to="/login">
              Login
            </Menu.Item>
            <Menu.Item as={Link} to="/register">
              Register
            </Menu.Item>
          </>
        )}
      </Menu>
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