import React from "react";
import PeopleList from "../Components/PeopleList";
import { useParams, useNavigate } from "react-router-dom";
import "./PeoplePage.css";

const PeoplePage = () => {
  const { teamId } = useParams();
  const navigate = useNavigate();

  return (
    <div className="people-page">
      <div className="header">
        <h2>Team Members</h2>
      </div>
      <PeopleList open={true} onClose={() => navigate(`/team/${teamId}`)} teamId={teamId} />
    </div>
  );
};

export default PeoplePage;