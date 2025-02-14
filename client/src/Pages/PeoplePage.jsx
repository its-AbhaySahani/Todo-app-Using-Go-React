import React, { useState } from "react";
import { Button, Modal, Form, Input } from "semantic-ui-react";
import PeopleList from "../Components/PeopleList";
import axios from "axios";
import { useParams, useNavigate } from "react-router-dom";
import "./PeoplePage.css";

const PeoplePage = () => {
  const { teamId } = useParams();
  const navigate = useNavigate();
  const [modalOpen, setModalOpen] = useState(false);
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
        setModalOpen(false);
      })
      .catch((error) => {
        setErrorMessage("Error adding person. Please try again.");
        setSuccessMessage("");
      });
  };

  return (
    <div className="people-page">
      <div className="header">
        <h2>Team Members</h2>
        <Button className="add-member-button" onClick={() => setModalOpen(true)}>
          Add Member
        </Button>
      </div>
      <PeopleList open={true} onClose={() => navigate(`/team/${teamId}`)} teamId={teamId} />
      <Modal open={modalOpen} onClose={() => setModalOpen(false)}>
        <Modal.Header>Add Member</Modal.Header>
        <Modal.Content>
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
          {successMessage && <p className="success-message">{successMessage}</p>}
          {errorMessage && <p className="error-message">{errorMessage}</p>}
        </Modal.Content>
      </Modal>
    </div>
  );
};

export default PeoplePage;