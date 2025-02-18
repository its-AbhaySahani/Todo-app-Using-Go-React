import React, { useEffect, useState } from "react";
import axios from "axios";
import { Card, Button } from "semantic-ui-react";
import { useNavigate } from "react-router-dom";
import JoinTeam from "./JoinTeam"; // Import JoinTeam component

const TeamDetailsBox = () => {
  const [teams, setTeams] = useState([]);
  const [joinTeamOpen, setJoinTeamOpen] = useState(false); // State to control JoinTeam modal
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("token");
    axios
      .get(`http://localhost:9000/api/teams`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then((res) => {
        setTeams(res.data || []);
      })
      .catch((error) => {
        console.error("Error fetching teams:", error);
      });
  }, []);

  return (
    <div className="team-details-box">
      <h2>My Teams</h2>
      <Button onClick={() => setJoinTeamOpen(true)}>Join Team</Button> {/* Button to open JoinTeam modal */}
      <Card.Group>
        {teams.length > 0 ? (
          teams.map((team) => (
            <Card key={team.id}>
              <Card.Content>
                <Card.Header>{team.name}</Card.Header>
                <Button onClick={() => navigate(`/team/${team.id}`)}>View Team</Button>
              </Card.Content>
            </Card>
          ))
        ) : (
          <p>Bhai tu koi team me nahi hai.</p>
        )}
      </Card.Group>
      <JoinTeam open={joinTeamOpen} onClose={() => setJoinTeamOpen(false)} /> {/* JoinTeam modal */}
    </div>
  );
};

export default TeamDetailsBox;