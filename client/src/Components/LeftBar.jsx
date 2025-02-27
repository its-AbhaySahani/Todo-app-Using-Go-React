// client/src/Components/LeftBar.jsx
import React from 'react';
import { useNavigate } from 'react-router-dom';
import './LeftBar.css';

const LeftBar = ({ setFilter, isTeamPage }) => {
  const navigate = useNavigate();
  
  return (
    <div className="leftbar-container">
      <div className="leftbar-heading">CHECKMATE</div>
      <div className="leftbar-menu">
        {!isTeamPage ? (
          <>
            <button onClick={() => setFilter('today')}>Today's Tasks</button>
            <button onClick={() => setFilter('all')}>All Tasks</button>
            <button onClick={() => setFilter('important')}>Important Tasks</button>
            <button onClick={() => setFilter('completed')}>Completed Tasks</button>
            <button onClick={() => setFilter('incomplete')}>Incomplete Tasks</button>
          </>
        ) : (
          <>
            <button onClick={() => setFilter('today')}>Today's Team Tasks</button>
            <button onClick={() => setFilter('all')}>All Team Tasks</button>
            <button onClick={() => setFilter('important')}>Important Team Tasks</button>
            <button onClick={() => setFilter('completed')}>Completed Team Tasks</button>
            <button onClick={() => setFilter('incomplete')}>Incomplete Team Tasks</button>
          </>
        )}
      </div>
      <div className="leftbar-bottom">
        <button className="routine-button" onClick={() => navigate('/my-routines')}>My Routines</button>
      </div>
    </div>
  );
};

export default LeftBar;