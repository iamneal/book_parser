
// TODO all globals should be put in a file that both go and js can read
let TOKEN_KEY = "GRPC-Metadata-book_parser_token"

class Watcher {
  constructor() {
    this.token = ""
  }

  getTokenFromCookie() {

  }

  updateLocalToken(tokenStr) {
    this.token = JSON.parse(tokenStr)
    console.log("the token: ", JSON.stringify(this.token, null, 2))
  }

  refreshToken() {

  }

  login() {
    return this.makePostRequest("/login").then((xmlhttp) => {
      console.log("success: " + xmlhttp.responseText)
      this.updateToken(xmlhttp.responseText)
    }, (xmlhttp) => {
      console.log("failure: " + xmlhttp.status)
    })
  }
  makePostRequest(path, paramsObj) {
    return new Promise((resolve, reject) => {
      let xmlhttp = new XMLHttpRequest()
      xmlhttp.onreadystatechange = function() {
        if (xmlhttp.readyState === XMLHttpRequest.DONE) {
          if (xmlhttp.status === 200) {//StatusOK
            resolve(xmlhttp) 
          } else {
            reject(xmlhttp)
          }
        }
      }
      let params = null
      if (paramsObj) {
        // add the token to our request if it exists
        if (this.token) {
          paramsObj[TOKEN_KEY] = this.token.value
        }
        // translate params to form string
        params = Object.keys(paramsObj).
          map((key) => key + "=" + paramsObj[key]).
          join("&").
          replace(/%20/g, "+")
      }
      xmlhttp.open("POST", path)
      xmlhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
      xmlhttp.setRequestHeader("Access-Control-Allow-Origin", "*")
      //xmlhttp.setRequestHeader(TOKEN_KEY, this.token.value)
      //xmlhttp.send(params)
      //xmlhttp.send()
    })
  }
}

module.exports.Watcher = Watcher
