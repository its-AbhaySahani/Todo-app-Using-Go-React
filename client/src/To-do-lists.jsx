import React, { Component } from "react";
import axios from "axios";
import { Card, Modal, Icon } from "semantic-ui-react";
import moment from "moment";
import Box from "./Box";
import TodoForm from "./Components/TodoForm";
import "./Box.css";
import "./To-do-lists.css";

let endpoint = "http://localhost:9000";

class ToDoList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      items: [],
      editTaskId: null,
      modalOpen: false,
      task: "",
      description: "",
      important: false,
      dateTime: new Date(),
      morning: false,
      noon: false,
      evening: false,
      night: false
    };
  }

  componentDidMount() {
    this.getTasks();
  }

  componentDidUpdate(prevProps) {
    if (prevProps.filter !== this.props.filter) {
      this.getTasks();
    }
  }

  getTasks = () => {
    const token = localStorage.getItem("token");
    axios
      .get(`${endpoint}/api/todos`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then((res) => {
        if (res.data) {
          let items = res.data.map((item) => ({
            ...item,
            dateTime: moment.utc(`${item.date} ${item.time}`).toDate(),
          }));

          // Apply filter
          const today = moment().format("YYYY-MM-DD");
          if (this.props.filter === 'today') {
            items = items.filter(item => item.date === today);
          } else if (this.props.filter === 'important') {
            items = items.filter(item => item.important);
          } else if (this.props.filter === 'completed') {
            items = items.filter(item => item.done);
          } else if (this.props.filter === 'incomplete') {
            items = items.filter(item => !item.done);
          }

          this.setState({ items });
        }
      });
  };

  updateTask = (id) => {
    const token = localStorage.getItem("token");
    axios
      .put(
        `${endpoint}/api/todo/${id}`,
        {
          done: true,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      )
      .then(() => {
        this.getTasks();
      });
  };

  deleteTask = (id) => {
    const token = localStorage.getItem("token");
    axios
      .delete(`${endpoint}/api/todo/${id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then(() => {
        this.getTasks();
      });
  };

  undoTask = (id) => {
    const token = localStorage.getItem("token");
    axios
      .put(
        `${endpoint}/api/todo/undo/${id}`,
        {
          done: false,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      )
      .then(() => {
        this.getTasks();
      });
  };


editTask = (id, task, description, date, time, important) => {
  // Get the current routine settings for this task when editing
  const token = localStorage.getItem("token");
  
  axios.get(`${endpoint}/api/routine/task/${id}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })
  .then((res) => {
    const routines = res.data || [];
    const dayOfWeek = moment().format('dddd').toLowerCase();
    
    // Set the checkboxes based on active routines for today
    let morning = false;
    let noon = false;
    let evening = false;
    let night = false;
    let selectedDay = dayOfWeek; // Default to current day
    
    // If there are active routines, get the day from the first one
    if (routines.length > 0 && routines.some(r => r.isActive)) {
      const activeRoutines = routines.filter(r => r.isActive);
      if (activeRoutines.length > 0) {
        selectedDay = activeRoutines[0].day;
      }
      
      // Check which schedule types are active
      routines.forEach(routine => {
        if (routine.isActive && routine.day === selectedDay) {
          switch(routine.scheduleType) {
            case 'morning':
              morning = true;
              break;
            case 'noon':
              noon = true;
              break;
            case 'evening':
              evening = true;
              break;
            case 'night':
              night = true;
              break;
            default:
              break;
          }
        }
      });
    }
    
    this.setState({ 
      task, 
      description, 
      important, 
      editTaskId: id, 
      dateTime: moment.utc(`${date} ${time}`).toDate(),
      morning,
      noon,
      evening,
      night,
      selectedDay,
      modalOpen: true 
    });
  })
  .catch(err => {
    // If there's an error or no routines, just open the modal without setting routine checkboxes
    this.setState({ 
      task, 
      description, 
      important, 
      editTaskId: id, 
      dateTime: moment.utc(`${date} ${time}`).toDate(),
      morning: false,
      noon: false,
      evening: false,
      night: false,
      selectedDay: moment().format('dddd').toLowerCase(), // Default to current day
      modalOpen: true 
    });
  });
};

  openModal = () => {
    this.setState({ 
      modalOpen: true,
      task: "",
      description: "",
      important: false,
      morning: false,
      noon: false,
      evening: false,
      night: false,
      editTaskId: null,
      dateTime: new Date()
    });
  };

  closeModal = () => {
    this.setState({ modalOpen: false });
  };

  handleFormSubmit = () => {
    this.getTasks();
    this.setState({ 
      modalOpen: false,
      task: "",
      description: "",
      important: false,
      morning: false,
      noon: false,
      evening: false,
      night: false,
      editTaskId: null,
      dateTime: new Date()
    });
  };

  render() {
    return (
      <div className="todo-list-container">
        <div className="row">
          {/* Header removed as in your original code */}
        </div>
        <div className="row">
          <Card.Group>
            <Card className="box-card add-task-card" onClick={this.openModal}>
              <Card.Content>
                <Card.Header textAlign="center">
                  <Icon name="plus" size="huge" />
                </Card.Header>
              </Card.Content>
            </Card>
            {this.state.items.map((item) => (
              <Box
                key={item.id}
                item={item}
                editTask={this.editTask}
                updateTask={this.updateTask}
                undoTask={this.undoTask}
                deleteTask={this.deleteTask}
              />
            ))}
          </Card.Group>
        </div>
        <Modal open={this.state.modalOpen} onClose={this.closeModal}>
          <Modal.Header>{this.state.editTaskId ? "Update Task" : "Add Task"}</Modal.Header>
          <Modal.Content>
            <TodoForm
              initialTask={this.state.task}
              initialDescription={this.state.description}
              initialImportant={this.state.important}
              initialDateTime={this.state.dateTime}
              editTaskId={this.state.editTaskId}
              initialMorning={this.state.morning}
              initialNoon={this.state.noon}
              initialEvening={this.state.evening}
              initialNight={this.state.night}
              onFormSubmit={this.handleFormSubmit}
            />
          </Modal.Content>
        </Modal>
      </div>
    );
  }
}

export default ToDoList;