import React, { useEffect, useState } from "react";
import axios from "axios";
import { Card, Header, Button } from "semantic-ui-react";
import Box from "../Box";
import "./GetTogether.css";

const GetTogether = () => {
  const [sharedItems, setSharedItems] = useState([]);
  const [receivedItems, setReceivedItems] = useState([]);
  const [activeSection, setActiveSection] = useState("shared");

  useEffect(() => {
    const token = localStorage.getItem("token");
    axios
      .get(`http://localhost:9000/api/shared`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then((res) => {
        console.log("API Response:", res.data); // Log the API response
        if (res.data) {
          console.log("Shared tasks fetched:", res.data);
          const shared = res.data.shared || [];
          const received = res.data.received || [];
          console.log("Filtered shared tasks:", shared);
          console.log("Filtered received tasks:", received);
          setSharedItems(shared);
          setReceivedItems(received);
        }
      })
      .catch((error) => {
        console.error("Error fetching shared tasks:", error);
      });
  }, []);

  useEffect(() => {
    console.log("Active section:", activeSection);
    console.log("Shared items:", sharedItems);
    console.log("Received items:", receivedItems);
  }, [activeSection, sharedItems, receivedItems]);

  const renderSharedItems = () => (
    <div className="column">
      <Header as="h3">Tasks I Shared</Header>
      {sharedItems.length > 0 ? (
        <Card.Group>
          {sharedItems.map((item) => (
            <Box key={item.id} item={item} />
          ))}
        </Card.Group>
      ) : (
        <p>No tasks shared by you.</p>
      )}
    </div>
  );

  const renderReceivedItems = () => (
    <div className="column">
      <Header as="h3">Tasks Shared With Me</Header> 
      {receivedItems.length > 0 ? (
        <Card.Group>
          {receivedItems.map((item) => (
            <Box key={item.id} item={item} />
          ))}
        </Card.Group>
      ) : (
        <p>No tasks shared with you.</p>
      )}
    </div>
  );

  return (
    <div className="get-together">
      <Header as="h2" color="yellow">
        Shared Tasks
      </Header>
      <div className="button-group">
        <Button
          color={activeSection === "shared" ? "blue" : "grey"}
          onClick={() => setActiveSection("shared")}
        >
          Tasks I Shared
        </Button>
        <Button
          color={activeSection === "received" ? "blue" : "grey"}
          onClick={() => setActiveSection("received")}
        >
          Tasks Shared With Me
        </Button>
      </div>
      <div className="columns">
        {activeSection === "shared" && renderSharedItems()}
        {activeSection === "received" && renderReceivedItems()}
      </div>
    </div>
  );
};

export default GetTogether;