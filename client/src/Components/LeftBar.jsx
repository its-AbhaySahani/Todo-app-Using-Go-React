import React from 'react';
import './LeftBar.css';

const LeftBar = ({ setFilter, isTeamPage }) => {
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
    </div>
  );
};

export default LeftBar;