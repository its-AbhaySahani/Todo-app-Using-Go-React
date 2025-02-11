import React, { useState } from "react";
import { Card, Icon } from "semantic-ui-react";
import moment from "moment";
import ShareBox from "./ShareBox";

const Box = ({ item, editTask, updateTask, undoTask, deleteTask }) => {
  const [shareOpen, setShareOpen] = useState(false);

  return (
    <Card key={item.id} color={item.done ? "green" : "red"}>
      <Card.Content>
        <Card.Header textAlign="left">
          <div style={{ wordWrap: "break-word" }}>{item.task}</div>
        </Card.Header>
        <Card.Meta textAlign="left">
          <span>{moment.utc(`${item.date} ${item.time}`).format("YYYY-MM-DD HH:mm:ss")}</span>
        </Card.Meta>
        <Card.Meta textAlign="right">
          <Icon
            name="edit"
            color="blue"
            onClick={() => editTask(item.id, item.task, item.date, item.time)}
          />
          <span style={{ paddingRight: 10 }}>Edit</span>
          <Icon
            name="check circle"
            color="green"
            onClick={() => updateTask(item.id)}
          />
          <span style={{ paddingRight: 10 }}>Done</span>
          <Icon
            name="undo"
            color="yellow"
            onClick={() => undoTask(item.id)}
          />
          <span style={{ paddingRight: 10 }}>Undo</span>
          <Icon
            name="delete"
            color="red"
            onClick={() => deleteTask(item.id)}
          />
          <span style={{ paddingRight: 10 }}>Delete</span>
          <Icon
            name="share alternate"
            color="blue"
            onClick={() => setShareOpen(true)}
          />
          <span>Share</span>
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