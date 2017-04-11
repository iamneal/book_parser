import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';
import { Watcher } from "./watcher.js"; 

class App extends Component {
  constructor() {
    super()
    this.bindAll()
    this.state = { watcher: new Watcher() }
  }

  login(e) {
    window.location = "/login"
  }

  bindAll() {
    this.login = this.login.bind(this)
  }

  render() {
    return (
      <div className="App">
        <div className="App-header">
          <h2>Welcome to React</h2>
        </div>
        <p className="App-intro">
          To get started, edit <code>src/App.js</code> and save to reload.
        </p>
        {(this.state.watcher.token === "") ? (
          <div>
            <button onClick={this.login}>
              Login
            </button>
          </div>
        ) : (
          <div> Nothing </div>
        )}
      </div>
    );
  }
}

export default App;
