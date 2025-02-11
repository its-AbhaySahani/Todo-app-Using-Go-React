import React, { useState } from "react";
import { Modal, Form, Input, Button, Message } from "semantic-ui-react";
import axios from "axios";

const ShareBox = ({ open, onClose, taskId }) => {
  const [username, setUsername] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [errorMessage, setErrorMessage] = useState("");

  const handleShare = () => {
    const token = localStorage.getItem("token");
    axios
      .post(
        `http://localhost:9000/api/share`,
        {
          taskId,
          username,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      )
      .then(() => {
        setSuccessMessage("Task shared successfully!");
        setErrorMessage("");
        setUsername("");
      })
      .catch((error) => {
        setErrorMessage("Error sharing task. Please try again.");
        setSuccessMessage("");
      });
  };

  return (
    <Modal open={open} onClose={onClose}>
      <Modal.Header>Share Task</Modal.Header>
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
          <Button type="button" onClick={handleShare}>
            Share
          </Button>
        </Form>
        {successMessage && <Message positive>{successMessage}</Message>}
        {errorMessage && <Message negative>{errorMessage}</Message>}
      </Modal.Content>
    </Modal>
  );
};

export default ShareBox;