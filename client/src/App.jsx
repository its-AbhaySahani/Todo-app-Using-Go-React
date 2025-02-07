import React from "react";
import "./App.css";
import { Container } from "semantic-ui-react";
import ToDoList from "./To-do-lists";

function App() {
  return (
    <Container>
      <ToDoList />
    </Container>
  );
}

export default App;