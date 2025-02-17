import React, { useState } from "react";
import { Button, Form, Input, Message } from "semantic-ui-react";
import axios from "axios";
import { useParams, useNavigate } from "react-router-dom";
import "./AddMemberPage.css";

const AddMemberPage = () => {
  const { teamId } = useParams();
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");

  const handleAddPerson = () => {
    const token = localStorage.getItem("token");
    axios
      .post(
        `http://localhost:9000/api/team/${teamId}/member`,
        { username },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      )
      .then(() => {
        setSuccessMessage("Person added successfully!");
        setErrorMessage("");
        setUsername("");
      })
      .catch((error) => {
        setErrorMessage("Error adding person. Please try again.");
        setSuccessMessage("");
      });
  };

  return (
    <div className="add-member-page">
      <h2>Add Member</h2>
      <Form>
        <Form.Field>
          <label>Username</label>
          <Input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="Enter username"
          />
        </Form.Field>
        <Button type="button" onClick={handleAddPerson}>
          Add
        </Button>
      </Form>
      {successMessage && <Message positive>{successMessage}</Message>}
      {errorMessage && <Message negative>{errorMessage}</Message>}
      <Button onClick={() => navigate(`/team/${teamId}/people`)}>Back to Members</Button>
    </div>
  );
};

export default AddMemberPage;