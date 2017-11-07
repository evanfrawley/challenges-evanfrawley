import React from 'react';
import {Row, Input, Button} from 'react-materialize';
import {withRouter} from 'react-router';
import {getUsersFromPrefix} from '../services/AuthAPIService';

class Search extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            foundUsers: [],
        }
    }

    triggerSearch = (prefix) => {
        getUsersFromPrefix(prefix).then((response) => {
            this.setState({foundUsers: response})
        })
    };

    _handleChange = (e) => {
        if (e.target.value !== '') {
            this.triggerSearch(e.target.value);
        }
    };

    render() {
        let searchResults = this.state.foundUsers.map((item) => {
            return (
                <li key={item.id}>
                    <img src={item.photourl} />
                    <span>{`${item.firstname} ${item.lastname}`}</span>
                    <span>{item.username}</span>
                    <span>{item.email}</span>
                </li>
            )
        });
        return (
            <div>
                <div>
                    <p>Search for other users!</p>
                    <Row>
                        <Input type={"text"} onChange={this._handleChange} />
                    </Row>
                </div>
                <div>
                    <ul>
                        {searchResults}
                    </ul>
                </div>
                <div>
                    <Button onClick={this.props.handleSignOut}>Sign Out</Button>
                </div>
            </div>
        );
    }
}

export default withRouter(Search);