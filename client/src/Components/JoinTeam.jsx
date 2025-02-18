import React, { useState } from "react";
import axios from "axios";
import { Modal, Form, Input, Button, Message } from "semantic-ui-react";
import "./JoinTeam.css";

const JoinTeam = ({ open, onClose }) => {
  const [teamName, setTeamName] = useState("");
  const [password, setPassword] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [errorMessage, setErrorMessage] = useState("");

  const handleJoinTeam = () => {
    const token = localStorage.getItem("token");
    axios
      .post(
        `http://localhost:9000/api/team/join`,
        {
          teamName,
          password,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      )
      .then(() => {
        setSuccessMessage("Joined team successfully!");
        setErrorMessage("");
        setTeamName("");
        setPassword("");
        onClose(); // Close the modal after successful join
      })
      .catch((error) => {
        setErrorMessage("Error joining team. Please try again.");
        setSuccessMessage("");
      });
  };

  return (
    <Modal open={open} onClose={onClose}>
      <Modal.Header>Join Team</Modal.Header>
      <Modal.Content>
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
          <Button type="button" onClick={handleJoinTeam}>
            Join Team
          </Button>
        </Form>
        {successMessage && <Message positive>{successMessage}</Message>}
        {errorMessage && <Message negative>{errorMessage}</Message>}
      </Modal.Content>
    </Modal>
  );
};

export default JoinTeam;