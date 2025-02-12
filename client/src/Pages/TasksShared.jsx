import React, { useEffect, useState } from "react";
import axios from "axios";
import { Card, Header } from "semantic-ui-react";
import Box from "../Box";
import "./GetTogether.css";

const TasksShared = () => {
  const [sharedItems, setSharedItems] = useState([]);

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
          const shared = res.data.filter(item => item.shared_by === token);
          console.log("Filtered shared tasks:", shared);
          setSharedItems(shared);
        }
      })
      .catch((error) => {
        console.error("Error fetching shared tasks:", error);
      });
  }, []);

  return (
    <div className="get-together">
      <Header as="h2" color="yellow">
        Tasks I Shared
      </Header>
      <div className="columns">
        <div className="column">
          <Card.Group>
            {sharedItems.map((item) => (
              <Box key={item.id} item={item} />
            ))}
          </Card.Group>
        </div>
      </div>
    </div>
  );
};

export default TasksShared;