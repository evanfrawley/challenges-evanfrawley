import React from 'react';

export default class MessageInput extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            value: '',
        }
    }

    _handleKeyPress = (e) => {
        // e.preventDefault();
        if (e.charCode === 13) {
            this.props.handleSubmit(this.state.value);
            this.setState({value: ''});
        }
    };

    _onChange = (e) => {
        e.preventDefault();
        this.setState({value: e.target.value});
    };

    render() {
        return(
            <div className={this.props.className}>
                <input
                  onChange={this._onChange}
                  onKeyPress={this._handleKeyPress}
                  value={this.state.value}
                  placeholder={"Type a message..."}
                />
            </div>
        );
    }
}