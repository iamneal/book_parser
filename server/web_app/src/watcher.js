import Promise from "bluebird";
import {TOKEN_KEY} from "../../globals";

class Watcher {
  constructor() {
    this.token = ""
    this.listeners = {}
  }
  setToken(tok) {
    this.token = ""
    console.log("token now: ", this.token)
  }

  getTokenFromCookie() {

  }
  // will give all listeners the successful/failing xmlhttp request
  registerListener(key, func) {
    if (!this.listeners[key]) {
      this.listeners[key] = []
    }
    this.listeners[key].push(func)
  }

  updateLocalToken(tokenStr) {
    this.token = JSON.parse(tokenStr)
    console.log("the token: ", JSON.stringify(this.token, null, 2))
  }

  refreshToken() {
    return this.makePostRequest("/update/token", null, true).then((xmlhttp) => {
      this.token = xmlhttp.response
    }, (xmlhttp) => {
      if (xmlhttp.status >= 401) {
        this.token = ""
      }
    })
  }

  parseLoginToken() {
    let tok = this.getQueryVariable(TOKEN_KEY)
    if (tok) {
      this.token = tok
    }
  }

	getQueryVariable(variable) {
    let query = window.location.search.substring(1);
    let vars = query.split('&');
    for (let i = 0; i < vars.length; i++) {
        let pair = vars[i].split('=');
        if (decodeURIComponent(pair[0]) === variable) {
            return decodeURIComponent(pair[1]);
        }
    }
  }
  makePostRequest(path, paramsObj, recursed) {
    console.log("PARAMSS OBJ: ", JSON.stringify(paramsObj, null, 2))
    return new Promise((resolve, reject) => {
      let xmlhttp = new XMLHttpRequest()
      xmlhttp.onreadystatechange = () => {
        if (xmlhttp.readyState === XMLHttpRequest.DONE) {
          // call all our listeners with the response, and the path
          let listeners = this.listeners[path]
          if (listeners) {
            listeners.map((f) => f(xmlhttp, path))
          }
          // give our response back to who called us
          if (xmlhttp.status === 200) {//StatusOK
            resolve(xmlhttp)
          } else if (xmlhttp.status === 401 && recursed !== true) {//Unauthorized
            console.log("refreshing token, then retrying request: status, recursed ", xmlhttp.status, recursed)
            this.refreshToken()
              .then(() => this.makePostRequest(path, paramsObj, true))
              .then((xmlhttp) => resolve(xmlhttp))
              .catch((xmlhttp) => reject(xmlhttp))
          } else {
            reject(xmlhttp)
          }
        }
      }
      let params = null
      if (paramsObj) {
        // add the token to our request if it exists
        if (this.token) {
          console.log("TOKEN_KEY, this.token", TOKEN_KEY, this.token)
          paramsObj.token = this.token;
          console.log("paramsObj now: ", paramsObj)

        }
        params = JSON.stringify(paramsObj)
        // translate params to form string
        //params = Object.keys(paramsObj).
        //  map((key) => key + "=" + paramsObj[key]).
        //  join("&").
        //  replace(/%20/g, "+")
      }
      xmlhttp.open("POST", path)
      //xmlhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
      xmlhttp.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
      xmlhttp.setRequestHeader("Access-Control-Allow-Origin", "*")
      xmlhttp.setRequestHeader(TOKEN_KEY, this.token)
      console.log("params: ", params)
      xmlhttp.send(params)
    })
  }
}

module.exports = new Watcher()
