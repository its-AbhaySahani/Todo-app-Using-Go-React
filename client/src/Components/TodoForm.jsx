// client/src/Components/TodoForm.jsx
import React, { Component } from "react";
import { Form, Input, Button, Checkbox, Dropdown } from "semantic-ui-react";
import { LocalizationProvider } from "@mui/lab";
import AdapterDateFns from "@mui/lab/AdapterDateFns";
import DateTimePicker from "@mui/lab/DateTimePicker";
import TextField from "@mui/material/TextField";
import moment from "moment";
import axios from "axios";
import "./TodoForm.css"; // Import the new CSS file

let endpoint = "http://localhost:9000";

class TodoForm extends Component {
  constructor(props) {
    super(props);
    
    // Get current day of the week in lowercase (e.g., "monday")
    const currentDay = moment().format('dddd').toLowerCase();
    
    this.state = {
      task: props.initialTask || "",
      description: props.initialDescription || "",
      important: props.initialImportant || false,
      dateTime: props.initialDateTime || new Date(),
      editTaskId: props.editTaskId || null,
      morning: props.initialMorning || false,
      noon: props.initialNoon || false,
      evening: props.initialEvening || false,
      night: props.initialNight || false,
      selectedDay: props.initialDay || currentDay
    };
    
    this.dayOptions = [
      { key: 'sunday', text: 'Sunday', value: 'sunday' },
      { key: 'monday', text: 'Monday', value: 'monday' },
      { key: 'tuesday', text: 'Tuesday', value: 'tuesday' },
      { key: 'wednesday', text: 'Wednesday', value: 'wednesday' },
      { key: 'thursday', text: 'Thursday', value: 'thursday' },
      { key: 'friday', text: 'Friday', value: 'friday' },
      { key: 'saturday', text: 'Saturday', value: 'saturday' }
    ];
  }

  componentDidUpdate(prevProps) {
    if (
      prevProps.initialTask !== this.props.initialTask ||
      prevProps.editTaskId !== this.props.editTaskId
    ) {
      // Get current day of the week in lowercase
      const currentDay = moment().format('dddd').toLowerCase();
      
      this.setState({
        task: this.props.initialTask || "",
        description: this.props.initialDescription || "",
        important: this.props.initialImportant || false,
        dateTime: this.props.initialDateTime || new Date(),
        editTaskId: this.props.editTaskId || null,
        morning: this.props.initialMorning || false,
        noon: this.props.initialNoon || false,
        evening: this.props.initialEvening || false,
        night: this.props.initialNight || false,
        selectedDay: this.props.initialDay || currentDay
      });
    }
  }

  onChange = (e) => {
    this.setState({ [e.target.name]: e.target.value });
  };

  onCheckboxChange = (e, { checked }) => {
    this.setState({ important: checked });
  };

  onDateTimeChange = (dateTime) => {
    this.setState({ dateTime });
  };

  onDayChange = (e, { value }) => {
    this.setState({ selectedDay: value });
  };

  onSubmit = (e) => {
    e.preventDefault();
    const token = localStorage.getItem("token");
    const formattedDate = moment.utc(this.state.dateTime).format("YYYY-MM-DD");
    const formattedTime = moment.utc(this.state.dateTime).format("HH:mm:ss");
    
    // Create a list of selected scheduleTypes
    const schedules = [];
    if (this.state.morning) schedules.push("morning");
    if (this.state.noon) schedules.push("noon");
    if (this.state.evening) schedules.push("evening");
    if (this.state.night) schedules.push("night");
    
    console.log("Selected day:", this.state.selectedDay);

    if (this.state.task) {
      if (this.state.editTaskId) {
        // Update task
        axios
          .put(
            `${endpoint}/api/todo/${this.state.editTaskId}`,
            {
              task: this.state.task,
              description: this.state.description,
              important: this.state.important,
              done: false,
              date: formattedDate,
              time: formattedTime,
            },
            {
              headers: {
                Authorization: `Bearer ${token}`,
              },
            }
          )
          .then((res) => {
            // After updating the task, update the routines if any schedules were selected
            if (schedules.length > 0) {
              axios.post(
                `${endpoint}/api/routine`,
                {
                  taskId: this.state.editTaskId,
                  schedules: schedules,
                  day: this.state.selectedDay
                },
                {
                  headers: {
                    Authorization: `Bearer ${token}`,
                  },
                }
              ).then(response => {
                console.log("Routine response:", response.data);
              }).catch(err => {
                console.error("Error creating/updating routine:", err);
              });
            }
            // Reset form and notify parent component
            this.props.onFormSubmit();
          });
      } else {
        // Create new task
        axios
          .post(
            `${endpoint}/api/todo`,
            {
              task: this.state.task,
              description: this.state.description,
              important: this.state.important,
              done: false,
              date: formattedDate,
              time: formattedTime,
            },
            {
              headers: {
                Authorization: `Bearer ${token}`,
              },
            }
          )
          .then((res) => {
            // After creating the task, create routines if any schedules were selected
            if (schedules.length > 0) {
              axios.post(
                `${endpoint}/api/routine`,
                {
                  taskId: res.data.id,
                  schedules: schedules,
                  day: this.state.selectedDay
                },
                {
                  headers: {
                    Authorization: `Bearer ${token}`,
                  },
                }
              ).then(response => {
                console.log("Routine response:", response.data);
              }).catch(err => {
                console.error("Error creating/updating routine:", err);
              });
            }
            // Reset form and notify parent component
            this.props.onFormSubmit();
          });
      }
    }
  };

  render() {
    return (
      <Form onSubmit={this.onSubmit} className="todo-form">
        <Form.Field>
          <label>Task Name</label>
          <Input
            type="text"
            name="task"
            onChange={this.onChange}
            value={this.state.task}
            fluid
            placeholder="Task Name"
          />
        </Form.Field>
        <Form.Field>
          <label>Description</label>
          <Input
            type="text"
            name="description"
            onChange={this.onChange}
            value={this.state.description}
            fluid
            placeholder="Description"
          />
        </Form.Field>
        <Form.Field>
          <label>Mark as Important</label>
          <Checkbox
            toggle
            checked={this.state.important}
            onChange={this.onCheckboxChange}
          />
        </Form.Field>
        <Form.Field>
          <LocalizationProvider dateAdapter={AdapterDateFns}>
            <DateTimePicker
              label="Date & Time"
              value={this.state.dateTime}
              onChange={this.onDateTimeChange}
              renderInput={(params) => <TextField {...params} />}
            />
          </LocalizationProvider>
        </Form.Field>
        
        {/* Routine Section Heading */}
        <h4>Daily Routine Settings</h4>
        
        {/* Day Selection Dropdown */}
        <Form.Field>
          <label>Day of Week</label>
          <Dropdown
            selection
            options={this.dayOptions}
            value={this.state.selectedDay}
            onChange={this.onDayChange}
            fluid
          />
        </Form.Field>
        
        {/* Routine Checkboxes */}
        <Form.Field>
          <label>Time of Day</label>
          <div className="routine-checkboxes">
            <div className="routine-checkbox">
              <Checkbox 
                label="Morning" 
                checked={this.state.morning}
                onChange={(e, { checked }) => this.setState({ morning: checked })}
              />
            </div>
            <div className="routine-checkbox">
              <Checkbox 
                label="Noon" 
                checked={this.state.noon}
                onChange={(e, { checked }) => this.setState({ noon: checked })}
              />
            </div>
            <div className="routine-checkbox">
              <Checkbox 
                label="Evening" 
                checked={this.state.evening}
                onChange={(e, { checked }) => this.setState({ evening: checked })}
              />
            </div>
            <div className="routine-checkbox">
              <Checkbox 
                label="Night" 
                checked={this.state.night}
                onChange={(e, { checked }) => this.setState({ night: checked })}
              />
            </div>
          </div>
        </Form.Field>
        
        <Button type="submit" className="submit-button">
          {this.state.editTaskId ? "Update Task" : "Add Task"}
        </Button>
      </Form>
    );
  }
}

export default TodoForm;