import React, { useEffect, useState } from "react";
import axios from "axios";
import { Card, Header } from "semantic-ui-react";
import Box from "../Box";
import "./GetTogether.css";

const GetTogether = () => {
  const [sharedItems, setSharedItems] = useState([]);
  const [receivedItems, setReceivedItems] = useState([]);

  useEffect(() => {
    const token = localStorage.getItem("token");
    axios
      .get(`http://localhost:9000/api/shared`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then((res) => {
        if (res.data) {
          console.log("Shared tasks fetched:", res.data);
          setSharedItems(res.data.filter(item => item.shared_by === token));
          setReceivedItems(res.data.filter(item => item.user_id === token));
        }
      })
      .catch((error) => {
        console.error("Error fetching shared tasks:", error);
      });
  }, []);

  return (
    <div className="get-together">
      <Header as="h2" color="yellow">
        Shared Tasks
      </Header>
      <div className="columns">
        <div className="column">
          <Header as="h3">Tasks I Shared</Header>
          <Card.Group>
            {sharedItems.map((item) => (
              <Box key={item.id} item={item} />
            ))}
          </Card.Group>
        </div>
        <div className="column">
          <Header as="h3">Tasks Shared With Me</Header>
          <Card.Group>
            {receivedItems.map((item) => (
              <Box key={item.id} item={item} />
            ))}
          </Card.Group>
        </div>
      </div>
    </div>
  );
};

export default GetTogether;