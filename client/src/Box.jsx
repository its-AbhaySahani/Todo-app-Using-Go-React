import React, { useState, useEffect } from "react";
import { Card, Icon } from "semantic-ui-react";
import moment from "moment";
import ShareBox from "./ShareBox";
import "./Box.css";

const Box = ({ item, editTask, updateTask, undoTask, deleteTask }) => {
  const [shareOpen, setShareOpen] = useState(false);

  useEffect(() => {
    console.log("Rendering Box component with item:", item);
  }, [item]);

  return (
    <Card key={item.id} className="box-card" color={item.done ? "green" : "red"}>
      <Card.Content>
        <Card.Header textAlign="left" className="box-card-header">
          <div style={{ wordWrap: "break-word" }}>{item.task}</div>
          {item.important && <Icon name="star" color="yellow" />}
        </Card.Header>
        <Card.Meta textAlign="left" className="box-card-meta">
          <span>{moment.utc(`${item.date} ${item.time}`).format("YYYY-MM-DD HH:mm:ss")}</span>
        </Card.Meta>
        <Card.Meta textAlign="right" className="box-card-actions">
          <div>
            <Icon
              name="edit"
              color="blue"
              onClick={() => editTask(item.id, item.task, item.description, item.date, item.time, item.important)}
              className="icon"
            />
            <span>Edit</span>
            <Icon
              name="check circle"
              color="green"
              onClick={() => updateTask(item.id)}
              className="icon"
            />
            <span>Done</span>
            <Icon
              name="undo"
              color="yellow"
              onClick={() => undoTask(item.id)}
              className="icon"
            />
            <span>Undo</span>
            <Icon
              name="delete"
              color="red"
              onClick={() => deleteTask(item.id)}
              className="icon"
            />
            <span>Delete</span>
            <Icon
              name="share alternate"
              color="blue"
              onClick={() => setShareOpen(true)}
              className="icon"
            />
            <span>Share</span>
          </div>
        </Card.Meta>
      </Card.Content>
      <ShareBox
        open={shareOpen}
        onClose={() => setShareOpen(false)}
        taskId={item.id}
      />
    </Card>
  );
};

export default Box;