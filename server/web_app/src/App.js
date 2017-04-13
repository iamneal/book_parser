import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';
import { Watcher } from "./watcher.js"; 
import Promise from "bluebird";

class App extends Component {
  constructor() {
    super()
    this.bindAll()
    this.state = {
      watcher: new Watcher(),
      currentDebugInfo: ""
    }
  }

  login(e) {
    window.location = "/login"
  }
  
  debugReq(e) {
    this.state.watcher.makePostRequest("/api/debug").then((xmlhttp) => {
      this.state.currentDebugInfo = xmlhttp.response
      this.setState(this.state)
    }, (xmlhttp) => {
      this.state.currentDebugInfo = "request failed status: " + xmlhttp.status
      this.setState(this.state)
    })
  }

  listBooks(e) {
    this.state.watcher.makePostRequest("/api/book/list").then((xmlhttp) => {
      this.state.currentDebugInfo = xmlhttp.response
    }).catch((xmlhttp) => {
      this.state.currentDebugInfo = "request failed status: " + xmlhttp.status
    }).finally(() => this.setState(this.state))
  }

  bindAll() {
    this.login = this.login.bind(this)
    this.debugReq = this.debugReq.bind(this)
    this.listBooks = this.listBooks.bind(this)
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
          <div>
            <button onClick={this.debugReq}>
              /api/debug
            </button>
            <button onClick={this.listBooks}>
              /api/book/list
            </button>
          </div>
        )}
        <div>
          <h3> Debug Div </h3>
          <div> {this.state.currentDebugInfo} </div>
        </div>
      </div>
    );
  }
}

export default App;
