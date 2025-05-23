import React, { useEffect, useState } from "react";
import axios from "axios";
import { Modal, List, Button, Icon } from "semantic-ui-react";
import { useNavigate } from "react-router-dom";
import "./PeopleList.css"; // Ensure this import is present

const PeopleList = ({ open, onClose, teamId }) => {
  const [people, setPeople] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("token");
    axios
      .get(`http://localhost:9000/api/team/${teamId}/members`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then((res) => {
        setPeople(res.data);
      })
      .catch((error) => {
        console.error("Error fetching team members:", error);
      });
  }, [teamId]);

  return (
    <Modal open={open} onClose={onClose}>
      <Modal.Header>Team Members</Modal.Header>
      <Modal.Content>
        <List>
          {people.map((person) => (
            <List.Item key={person.id}>
              <List.Content>
                <List.Header>
                  {person.username} {person.is_admin && <span className="admin-label">Admin</span>}
                </List.Header>
              </List.Content>
            </List.Item>
          ))}
        </List>
      </Modal.Content>
      <Modal.Actions>
        <Button onClick={() => navigate(`/team/${teamId}`)}>Close</Button>
        <Button color="blue" onClick={() => navigate(`/team/${teamId}/add-member`)}>
          <Icon name="user plus" /> Add Member
        </Button>
      </Modal.Actions>
    </Modal>
  );
};

export default PeopleList;