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
    console.log('pull')
    this.props.pullDoc(this.props.docs[index].id)
  }

  selectDoc(index) {
    console.log("select ", index)
    this.state.selectedDoc = index
    this.setState(this.state)
  }

  deselectDoc(index) {
    console.log("deselect ", index)
    this.state.selectedDoc = NaN
    this.setState(this.state)
  }

  render() {
    let docNodes = this.props.docs.map((doc, index) => {
      return (<Doc
        id={doc.Id}
        name={doc.Name}
        index={index}
        select={this.selectDoc}
        pull={this.pullDoc}
        deselect={this.deselectDoc}
        highlighted={this.state.selectedDoc === index}
        key={index}
      />)
    })
    return (
      <div className="ui one column stackable grid container">
          {docNodes}
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
      <div className="column">
      {(this.props.highlighted) ? (
        <div className="ui raised segmant" onClick={() => this.props.deselect(this.props.index)}>
          <div className="two wide column">
          <div className="ui segment">
            <p> {this.props.name} </p>
          </div>
          <div className="column">
          <div className="ui segment">
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
          </div>
        </div>

      ) : (
        <div className="ui raised segment" onClick={() => this.props.select(this.props.index)}>
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
