import React, { useEffect, useState } from "react";
import axios from "axios";
import { Card, Header } from "semantic-ui-react";
import Box from "../Box";
import "./GetTogether.css";

const TasksReceived = () => {
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
          const received = res.data.filter(item => item.user_id === token);
          console.log("Filtered received tasks:", received);
          setReceivedItems(received);
        }
      })
      .catch((error) => {
        console.error("Error fetching shared tasks:", error);
      });
  }, []);

  return (
    <div className="get-together">
      <Header as="h2" color="yellow">
        Tasks Shared With Me
      </Header>
      <div className="columns">
        <div className="column">
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

export default TasksReceived;