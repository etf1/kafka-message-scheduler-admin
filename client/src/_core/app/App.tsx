import React from 'react';
import './App.css';
import Navbar from '_common/component/navbar/Navbar';

function App() {
  return (
    <div className="App">
      <Navbar items={[{label:"Home",href:"/"}]}/>
    </div>
  );
}

export default App;
