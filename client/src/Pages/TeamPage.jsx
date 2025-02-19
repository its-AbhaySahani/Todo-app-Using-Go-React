import React, { useEffect, useState } from "react";
import axios from "axios";
import { useParams, useNavigate } from "react-router-dom";
import { Card, Header, Button, Form, Input, Modal, Checkbox, Icon } from "semantic-ui-react";
import { LocalizationProvider } from "@mui/x-date-pickers";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";
import { DateTimePicker } from "@mui/x-date-pickers/DateTimePicker";
import TextField from "@mui/material/TextField";
import moment from "moment";
import Box from "../Box";
import "./TeamPage.css";

const TeamPage = ({ filter }) => {
  const { teamId } = useParams();
  const [team, setTeam] = useState({});
  const [tasks, setTasks] = useState([]); // Initialize tasks to an empty array
  const [filteredTasks, setFilteredTasks] = useState([]); // State for filtered tasks
  const [modalOpen, setModalOpen] = useState(false);
  const [task, setTask] = useState("");
  const [description, setDescription] = useState("");
  const [important, setImportant] = useState(false);
  const [dateTime, setDateTime] = useState(new Date());
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("token");
    axios
      .get(`http://localhost:9000/api/team/${teamId}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then((res) => {
        setTeam(res.data.team);
        setTasks(res.data.tasks || []); // Ensure tasks is always an array
      })
      .catch((error) => {
        console.error("Error fetching team details:", error);
      });
  }, [teamId]);

  useEffect(() => {
    applyFilter();
  }, [filter, tasks]);

  const applyFilter = () => {
    const today = moment().format("YYYY-MM-DD");
    let filteredTasks = tasks;

    if (filter === 'today') {
      filteredTasks = tasks.filter(task => task.date === today);
    } else if (filter === 'important') {
      filteredTasks = tasks.filter(task => task.important);
    } else if (filter === 'completed') {
      filteredTasks = tasks.filter(task => task.done);
    } else if (filter === 'incomplete') {
      filteredTasks = tasks.filter(task => !task.done);
    }

    setFilteredTasks(filteredTasks);
  };

  const handleAddTask = () => {
    const token = localStorage.getItem("token");
    axios
      .post(
        `http://localhost:9000/api/team/${teamId}/todo`,
        {
          task,
          description,
          important,
          dateTime,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      )
      .then((res) => {
        setTasks([...tasks, res.data]);
        setModalOpen(false);
        setTask("");
        setDescription("");
        setImportant(false);
        setDateTime(new Date());
      })
      .catch((error) => {
        console.error("Error adding task:", error);
      });
  };

  return (
    <div className="team-page">
      <div className="team-header">
        <Header as="h2" className="team-name">
          {team.name}
        </Header>
        <Button className="people-button" onClick={() => navigate(`/team/${teamId}/people`)}>
          <Icon name="users" /> Members
        </Button>
      </div>
      <div className="task-section">
        <Card.Group>
          <Card className="box-card add-task-card" onClick={() => setModalOpen(true)}>
            <Card.Content>
              <Card.Header textAlign="center">
                <Icon name="plus" size="huge" />
              </Card.Header>
            </Card.Content>
          </Card>
          {filteredTasks.map((task) => (
            <Box key={task.id} item={task} />
          ))}
        </Card.Group>
      </div>
      <Modal open={modalOpen} onClose={() => setModalOpen(false)}>
        <Modal.Header>Add Task</Modal.Header>
        <Modal.Content>
          <Form>
            <Form.Field>
              <label>Task Name</label>
              <Input
                type="text"
                value={task}
                onChange={(e) => setTask(e.target.value)}
                placeholder="Enter task name"
              />
            </Form.Field>
            <Form.Field>
              <label>Description</label>
              <Input
                type="text"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Enter description"
              />
            </Form.Field>
            <Form.Field>
              <label>Mark as Important</label>
              <Checkbox
                toggle
                checked={important}
                onChange={(e, { checked }) => setImportant(checked)}
              />
            </Form.Field>
            <Form.Field>
              <LocalizationProvider dateAdapter={AdapterDateFns}>
                <DateTimePicker
                  label="Date & Time"
                  value={dateTime}
                  onChange={(date) => setDateTime(date)}
                  renderInput={(params) => <TextField {...params} />}
                />
              </LocalizationProvider>
            </Form.Field>
            <Button type="button" onClick={handleAddTask}>
              Add Task
            </Button>
          </Form>
        </Modal.Content>
      </Modal>
    </div>
  );
};

export default TeamPage;