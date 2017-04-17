import React, { Component } from "react";
import PropTypes from "prop-types";
import Promise from "bluebird";

class DocList extends Component {
  constructor() {
    super()
    this.state = {selectedDoc: NaN}
    this.bindAll()
  }
  
  bindAll() {
    this.pullDoc = this.pullDoc.bind(this)
    this.selectDoc = this.selectDoc.bind(this)
    this.deselectDoc = this.deselectDoc.bind(this)
  }

  pullDoc(index) {
    this.props.pullDoc(this.props.docs[index].id)
  }

  selectDoc(index) {
    this.state.selectedDoc = index
    this.setState(this.state)
  }

  deselectDoc(index) {
    this.state.selectedDoc = NaN
    this.setState(this.state)
  }

  render() {
    let docNodes = this.props.docs.map((doc, index) => {
      return (<Doc
        id={doc.id}
        name={doc.name}
        index={index}
        select={this.selectDoc}
        deselect={this.deselectDoc}
        pull={this.pullDoc}
        highlighted={this.state.selectedDoc === index}
      />)
    })
    return (
      <div>
        <h3> Book list </h3>
        <div className="BookList">
          {docNodes}
        </div>
      </div>
    )
  }
}

DocList.propTypes = {
  pullDoc: PropTypes.func.isRequired,
  docs: PropTypes.array
}

class Doc extends Component {
  constructor() {
    super() 
  }

  render() {
    return (
      <div>
      {(this.props.highlighted) ? (
        <div className="ExpandedBook" onClick={this.props.deselect(this.props.index)}>
          <div className="ExpandedBook-Header">
            <p> {this.props.name} </p>
          </div>
          <div className="ExpandedBook-Body">
            <div>
              <div> id: {this.props.id}</div>
              <div> pull status: (not used yet) </div>
            </div>
            <div>
              <div> 
                <button onClick={() => console.log("called pull")}>
                  pull book
                </button>
              </div>
            </div>
          </div>
        </div>

      ) : (
        <div className="UnexpandedBook" onClick={this.props.select(this.props.index)}>
          {this.props.name} 
        </div>
      )}
    </div>
    )
  }
}

Doc.propTypes = {
  id: PropTypes.string,
  name: PropTypes.string,
  index: PropTypes.number,
  select: PropTypes.func.isRequired,
  deselect: PropTypes.func.isRequired,
  pull: PropTypes.func.isRequired,
  highlighted: PropTypes.bool
}

module.exports.DocList = DocList
module.exports.Doc = Doc
