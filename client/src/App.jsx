import React from "react";
import { BrowserRouter as Router, Routes, Route, Link } from "react-router-dom";
import "./App.css";
import { Container, Menu } from "semantic-ui-react";
import ToDoList from "./To-do-lists";
import Login from "./Authentication/Login";
import Register from "./Authentication/Register";

function App() {
  return (
    <Router>
      <Container>
        <Menu>
          <Menu.Item as={Link} to="/">
            Home
          </Menu.Item>
          <Menu.Item as={Link} to="/login">
            Login
          </Menu.Item>
          <Menu.Item as={Link} to="/register">
            Register
          </Menu.Item>
        </Menu>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/" element={<ToDoList />} />
        </Routes>
      </Container>
    </Router>
  );
}

export default App;