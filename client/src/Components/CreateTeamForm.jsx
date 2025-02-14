import React, { useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { Form, Input, Button, Message } from "semantic-ui-react";

const CreateTeamForm = () => {
  const [teamName, setTeamName] = useState("");
  const [password, setPassword] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const navigate = useNavigate();

  const handleCreateTeam = () => {
    const token = localStorage.getItem("token");
    axios
      .post(
        `http://localhost:9000/api/team`,
        {
          name: teamName,
          password,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      )
      .then(() => {
        setSuccessMessage("Team created successfully!");
        setErrorMessage("");
        setTeamName("");
        setPassword("");
        navigate("/my-teams");
      })
      .catch((error) => {
        setErrorMessage("Error creating team. Please try again.");
        setSuccessMessage("");
      });
  };

  return (
    <div className="create-team-form">
      <h2>Create Team</h2>
      <Form>
        <Form.Field>
          <label>Team Name</label>
          <Input
            type="text"
            value={teamName}
            onChange={(e) => setTeamName(e.target.value)}
            placeholder="Enter team name"
          />
        </Form.Field>
        <Form.Field>
          <label>Password</label>
          <Input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Enter password"
          />
        </Form.Field>
        <Button type="button" onClick={handleCreateTeam}>
          Create Team
        </Button>
      </Form>
      {successMessage && <Message positive>{successMessage}</Message>}
      {errorMessage && <Message negative>{errorMessage}</Message>}
    </div>
  );
};

export default CreateTeamForm;