import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, Link, useNavigate } from "react-router-dom";
import "./App.css";
import { Container, Menu, Icon } from "semantic-ui-react";
import ToDoList from "./To-do-lists";
import Login from "./Authentication/Login";
import Register from "./Authentication/Register";
import GetTogether from "./Pages/GetTogether"; // Import GetTogether component
// import TasksShared from "./Pages/TasksShared";
// import TasksReceived from "./Pages/TasksReceived";


function App() {
  return (
    <Router>
      <AppContent />
    </Router>
  );
}

function AppContent() {
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
    <Container>
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
        <Route path="/" element={<ToDoList />} />
        <Route path="/get-together" element={<GetTogether />} /> {/* Add GetTogether route */}
        {/* <Route path="/tasks-shared" element={<TasksShared />} />
        <Route path="/tasks-received" element={<TasksReceived />} /> */}
      </Routes>
    </Container>
  );
}

export default App;