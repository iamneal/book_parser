'use strict'

class Watcher {
  constructor() {
    this.token = ""
  }

  getTokenFromCookie() {
  }
  updateToken(tokenStr) {
    this.token = JSON.parse(tokenStr)
    console.log("the token: ", JSON.stringify(this.token, null, 2))
  }
  login() {
    return makePostReqest("/login").then((xmlhttp) => {
      console.log("success")
    }, (xmlhttp) => {
      console.log("failure")
    })
  }
  makePostRequest(path, paramsObj) {
    return new Promise((resolve, reject) => {
      xmlhttp := new XmlHttpRequest()
      xmlhttp.onreadystatechange = function() {
        if (xmlhttp.readyState === XMLHttpRequest.DONE) {
          if (xmlhttp.status === 200) {//StatusOK
            resolve(xmlhttp) 
          } else {
            reject(xmlhttp)
          }
        }
      }
      params = null
      if paramsObj {
        // add the token to our request if it exists
        if (this.token) {
          paramsObj[TOKEN_KEY] = this.token.value
        }
        // translate params to form string
        params = Object.keys(paramsObj).
          map((key) => key + '=' + paramsObj[key]).
          join(&).
          replace(/%20/g, '+')
      }
      xmlhttp.open("POST", path)
      xmlhttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
      xmlhttp.setRequestHeader(TOKEN_HEADER, this.token.value)
      xmlhttp.send(params)
    })
  }
    
}
