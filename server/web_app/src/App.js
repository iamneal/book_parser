import React, { Component } from "react";
import logo from "./logo.svg";
import "./App.css";
import { Watcher } from "./watcher.js";
import Promise from "bluebird";
import {DocList} from "./BookList"
import "semantic-ui-css/semantic.min.css"

class App extends Component {
  constructor() {
    super()
    this.bindAll()
    this.state = {
      watcher: new Watcher(),
      currentDebugInfo: "",
      docs: []
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
      try {
        let docs = JSON.parse(xmlhttp.response)
        console.log(docs)
        this.state.docs = docs.books
      } catch(e) {
        this.state.currenDebugInfo = "failed to parse json: " + e
        console.log(e)
      }
    }).catch((xmlhttp) => {
      this.state.currentDebugInfo = "request failed status: " + xmlhttp.status
    }).finally(() => this.setState(this.state))
  }

  pullDoc(id) {
    console.log("pulling book with: id, obj", id, {id: id})
    this.state.watcher.makePostRequest("/api/book/pull", {id: id}).then((xmlhttp) => {
      console.log("success! ", xmlhttp.response)
      this.state.currentDebugInfo = xmlhttp.response
    }, (xmlhttp) => {
      this.state.currentDebugInfo = "request failed status: ", xmlhttp.status
    }).finally(() => this.setState(this.state))
  }

  bindAll() {
    this.login = this.login.bind(this)
    this.debugReq = this.debugReq.bind(this)
    this.listBooks = this.listBooks.bind(this)
    this.pullDoc = this.pullDoc.bind(this)
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
        <div className="ui two column stackable grid container">
          <div className="column">
            <h3> Debug Div </h3>
            <div> {this.state.currentDebugInfo} </div>
          </div>
          <div className="column">
            <h3> Doc list </h3>
            <DocList pullDoc={this.pullDoc} docs={this.state.docs}/>
          </div>
        </div>
      </div>
    );
  }
}

export default App;
