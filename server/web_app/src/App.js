import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';
import { Watcher } from "./fetcher.js"; 

class App extends Component {
  constructor() {
    super()
    this.state = { watcher: new Watcher() }
    this.bindAll()
  }

  login(e) {
    window.location = "/login"
    //e.preventDefault()
    //this.state.watcher.login()
  }

  bindAll() {
    this.login = this.login.bind(this)
  }

  render() {
    return (
      <div className="App">
        <div className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h2>Welcome to React</h2>
        </div>
        <p className="App-intro">
          To get started, edit <code>src/App.js</code> and save to reload.
        </p>
        <div>
          <button onClick={this.login}>
            Login
          </button>
        </div>
      </div>
    );
  }
}

export default App;
